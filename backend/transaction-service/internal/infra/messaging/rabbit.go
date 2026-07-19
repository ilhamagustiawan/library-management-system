package messaging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/repository"
)

const (
	DefaultExchange        = "library.events"
	DefaultDeadExchange    = "library.events.dlx"
	DefaultBookReturnQueue = "book-service.loan-returned.v1"
	DefaultStockAckQueue   = "transaction-service.book-stock-updated.v1"
	DefaultDeadLetterQueue = "transaction-service.dead-letter"
	LoanReturnedRoutingKey = "transactions.loan.returned.v1"
	StockUpdatedRoutingKey = "books.stock.updated.v1"
)

type RabbitConfig struct {
	URL             string
	Exchange        string
	DeadExchange    string
	BookReturnQueue string
	StockAckQueue   string
	DeadLetterQueue string
	ConfirmTimeout  time.Duration
}

func (c RabbitConfig) defaults() RabbitConfig {
	if c.Exchange == "" {
		c.Exchange = DefaultExchange
	}
	if c.DeadExchange == "" {
		c.DeadExchange = DefaultDeadExchange
	}
	if c.BookReturnQueue == "" {
		c.BookReturnQueue = DefaultBookReturnQueue
	}
	if c.StockAckQueue == "" {
		c.StockAckQueue = DefaultStockAckQueue
	}
	if c.DeadLetterQueue == "" {
		c.DeadLetterQueue = DefaultDeadLetterQueue
	}
	if c.ConfirmTimeout <= 0 {
		c.ConfirmTimeout = 5 * time.Second
	}
	return c
}

type RabbitPublisher struct {
	config  RabbitConfig
	mutex   sync.Mutex
	conn    *amqp.Connection
	channel *amqp.Channel
	returns <-chan amqp.Return
}

func NewRabbitPublisher(config RabbitConfig) (*RabbitPublisher, error) {
	config = config.defaults()
	if config.URL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}
	return &RabbitPublisher{config: config}, nil
}

func (p *RabbitPublisher) Publish(ctx context.Context, message repository.OutboxMessage) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if err := p.ensureChannel(); err != nil {
		return err
	}
	publishCtx, cancel := context.WithTimeout(ctx, p.config.ConfirmTimeout)
	defer cancel()
	// Persistent delivery plus publisher confirms follows RabbitMQ's data-safety guidance.
	// Source: https://www.rabbitmq.com/docs/publishers#data-safety
	confirmation, err := p.channel.PublishWithDeferredConfirmWithContext(
		publishCtx, p.config.Exchange, message.RoutingKey, true, false,
		amqp.Publishing{
			ContentType: "application/json", DeliveryMode: amqp.Persistent,
			MessageId: message.ID, Type: message.EventType, Timestamp: time.Now().UTC(), Body: message.Payload,
		},
	)
	if err != nil {
		p.reset()
		return fmt.Errorf("publish RabbitMQ event: %w", err)
	}
	confirmed, err := confirmation.WaitContext(publishCtx)
	if err != nil {
		p.reset()
		return fmt.Errorf("wait for RabbitMQ publisher confirm: %w", err)
	}
	if !confirmed {
		return fmt.Errorf("RabbitMQ negatively acknowledged event %s", message.ID)
	}
	select {
	case returned := <-p.returns:
		if returned.MessageId == message.ID {
			return fmt.Errorf("RabbitMQ could not route event %s", message.ID)
		}
	default:
	}
	return nil
}

func (p *RabbitPublisher) ensureChannel() error {
	if p.conn != nil && !p.conn.IsClosed() && p.channel != nil && !p.channel.IsClosed() {
		return nil
	}
	p.reset()
	conn, err := amqp.Dial(p.config.URL)
	if err != nil {
		return fmt.Errorf("connect RabbitMQ publisher: %w", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("open RabbitMQ publisher channel: %w", err)
	}
	if err := declareTopology(channel, p.config); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}
	if err := channel.Confirm(false); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return fmt.Errorf("enable RabbitMQ publisher confirms: %w", err)
	}
	p.conn, p.channel = conn, channel
	p.returns = channel.NotifyReturn(make(chan amqp.Return, 1))
	return nil
}

func (p *RabbitPublisher) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.reset()
	return nil
}

func (p *RabbitPublisher) reset() {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
	p.channel, p.conn, p.returns = nil, nil, nil
}

type RabbitAckConsumer struct {
	config    RabbitConfig
	processor *AckProcessor
}

func NewRabbitAckConsumer(config RabbitConfig, processor *AckProcessor) (*RabbitAckConsumer, error) {
	config = config.defaults()
	if config.URL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}
	if processor == nil {
		return nil, fmt.Errorf("stock ack processor is required")
	}
	return &RabbitAckConsumer{config: config, processor: processor}, nil
}

func (c *RabbitAckConsumer) Run(ctx context.Context) {
	for ctx.Err() == nil {
		if err := c.consume(ctx); err != nil && ctx.Err() == nil {
			slog.Error("RabbitMQ stock ack consumer stopped", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}

func (c *RabbitAckConsumer) consume(ctx context.Context) error {
	conn, err := amqp.Dial(c.config.URL)
	if err != nil {
		return fmt.Errorf("connect RabbitMQ consumer: %w", err)
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open RabbitMQ consumer channel: %w", err)
	}
	defer channel.Close()
	if err := declareTopology(channel, c.config); err != nil {
		return err
	}
	if err := channel.Qos(20, 0, false); err != nil {
		return fmt.Errorf("set RabbitMQ consumer prefetch: %w", err)
	}
	// Manual acknowledgements transfer responsibility only after the DB commit.
	// Source: https://www.rabbitmq.com/docs/confirms#consumer-acks
	deliveries, err := channel.ConsumeWithContext(ctx, c.config.StockAckQueue, "transaction-service-stock-ack", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume RabbitMQ stock acknowledgements: %w", err)
	}
	for delivery := range deliveries {
		if len(delivery.Body) > 1<<20 {
			_ = delivery.Reject(false)
			continue
		}
		if err := c.processor.Process(ctx, delivery.Body); err != nil {
			slog.Error("stock acknowledgement failed", "message_id", delivery.MessageId, "error", err)
			if errors.Is(err, ErrInvalidMessage) || delivery.Redelivered {
				_ = delivery.Reject(false)
			} else {
				_ = delivery.Nack(false, true)
			}
			continue
		}
		if err := delivery.Ack(false); err != nil {
			return fmt.Errorf("acknowledge RabbitMQ stock event: %w", err)
		}
	}
	if ctx.Err() != nil {
		return nil
	}
	return fmt.Errorf("RabbitMQ delivery channel closed")
}

func declareTopology(channel *amqp.Channel, config RabbitConfig) error {
	if err := channel.ExchangeDeclare(config.Exchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare RabbitMQ event exchange: %w", err)
	}
	if err := channel.ExchangeDeclare(config.DeadExchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare RabbitMQ dead-letter exchange: %w", err)
	}
	deadArgs := amqp.Table{"x-dead-letter-exchange": config.DeadExchange}
	if _, err := channel.QueueDeclare(config.BookReturnQueue, true, false, false, false, deadArgs); err != nil {
		return fmt.Errorf("declare Book return queue: %w", err)
	}
	if err := channel.QueueBind(config.BookReturnQueue, LoanReturnedRoutingKey, config.Exchange, false, nil); err != nil {
		return fmt.Errorf("bind Book return queue: %w", err)
	}
	if _, err := channel.QueueDeclare(config.StockAckQueue, true, false, false, false, deadArgs); err != nil {
		return fmt.Errorf("declare stock ack queue: %w", err)
	}
	if err := channel.QueueBind(config.StockAckQueue, StockUpdatedRoutingKey, config.Exchange, false, nil); err != nil {
		return fmt.Errorf("bind stock ack queue: %w", err)
	}
	if _, err := channel.QueueDeclare(config.DeadLetterQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare transaction dead-letter queue: %w", err)
	}
	if err := channel.QueueBind(config.DeadLetterQueue, "#", config.DeadExchange, false, nil); err != nil {
		return fmt.Errorf("bind transaction dead-letter queue: %w", err)
	}
	return nil
}
