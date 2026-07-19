package messaging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

const (
	loanReturnedKey = "transactions.loan.returned.v1"
	stockUpdatedKey = "books.stock.updated.v1"
)

type RabbitConfig struct {
	URL, Exchange, DeadExchange, ReturnQueue, AckQueue, DeadQueue string
	ConfirmTimeout                                                time.Duration
}

func (c RabbitConfig) defaults() RabbitConfig {
	if c.Exchange == "" {
		c.Exchange = "library.events"
	}
	if c.DeadExchange == "" {
		c.DeadExchange = "library.events.dlx"
	}
	if c.ReturnQueue == "" {
		c.ReturnQueue = "book-service.loan-returned.v1"
	}
	if c.AckQueue == "" {
		c.AckQueue = "transaction-service.book-stock-updated.v1"
	}
	if c.DeadQueue == "" {
		c.DeadQueue = "book-service.dead-letter"
	}
	if c.ConfirmTimeout <= 0 {
		c.ConfirmTimeout = 5 * time.Second
	}
	return c
}

type Rabbit struct {
	config    RabbitConfig
	processor *ReturnProcessor
	mutex     sync.Mutex
	conn      *amqp.Connection
	channel   *amqp.Channel
	returns   <-chan amqp.Return
}

func NewRabbit(config RabbitConfig, processor *ReturnProcessor) (*Rabbit, error) {
	config = config.defaults()
	if config.URL == "" || processor == nil {
		return nil, fmt.Errorf("RabbitMQ URL and return processor are required")
	}
	return &Rabbit{config: config, processor: processor}, nil
}

func (r *Rabbit) Publish(ctx context.Context, message repository.OutboxMessage) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if err := r.ensurePublisher(); err != nil {
		return err
	}
	publishCtx, cancel := context.WithTimeout(ctx, r.config.ConfirmTimeout)
	defer cancel()
	confirmation, err := r.channel.PublishWithDeferredConfirmWithContext(publishCtx, r.config.Exchange, message.RoutingKey, true, false, amqp.Publishing{ContentType: "application/json", DeliveryMode: amqp.Persistent, MessageId: message.ID, Type: message.EventType, Timestamp: time.Now().UTC(), Body: message.Payload})
	if err != nil {
		r.reset()
		return err
	}
	confirmed, err := confirmation.WaitContext(publishCtx)
	if err != nil {
		r.reset()
		return err
	}
	if !confirmed {
		return fmt.Errorf("RabbitMQ negatively acknowledged %s", message.ID)
	}
	select {
	case returned := <-r.returns:
		if returned.MessageId == message.ID {
			return fmt.Errorf("RabbitMQ could not route %s", message.ID)
		}
	default:
	}
	return nil
}

func (r *Rabbit) RunConsumer(ctx context.Context) {
	for ctx.Err() == nil {
		if err := r.consume(ctx); err != nil && ctx.Err() == nil {
			slog.Error("Book return consumer stopped", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}

func (r *Rabbit) consume(ctx context.Context) error {
	conn, err := amqp.Dial(r.config.URL)
	if err != nil {
		return err
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	if err := declareTopology(channel, r.config); err != nil {
		return err
	}
	if err := channel.Qos(20, 0, false); err != nil {
		return err
	}
	deliveries, err := channel.ConsumeWithContext(ctx, r.config.ReturnQueue, "book-service-loan-returned", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		if len(delivery.Body) > 1<<20 {
			_ = delivery.Reject(false)
			continue
		}
		if err := r.processor.Process(ctx, delivery.Body); err != nil {
			slog.Error("reject LoanReturned event", "message_id", delivery.MessageId, "error", err)
			if errors.Is(err, ErrInvalidMessage) || delivery.Redelivered {
				_ = delivery.Reject(false)
			} else {
				_ = delivery.Nack(false, true)
			}
			continue
		}
		if err := delivery.Ack(false); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (r *Rabbit) ensurePublisher() error {
	if r.conn != nil && !r.conn.IsClosed() && r.channel != nil && !r.channel.IsClosed() {
		return nil
	}
	r.reset()
	conn, err := amqp.Dial(r.config.URL)
	if err != nil {
		return err
	}
	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	if err := declareTopology(channel, r.config); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}
	if err := channel.Confirm(false); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}
	r.conn, r.channel, r.returns = conn, channel, channel.NotifyReturn(make(chan amqp.Return, 1))
	return nil
}

func (r *Rabbit) Close() { r.mutex.Lock(); defer r.mutex.Unlock(); r.reset() }
func (r *Rabbit) reset() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
	r.channel, r.conn, r.returns = nil, nil, nil
}

func declareTopology(channel *amqp.Channel, config RabbitConfig) error {
	if err := channel.ExchangeDeclare(config.Exchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	if err := channel.ExchangeDeclare(config.DeadExchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	deadArgs := amqp.Table{"x-dead-letter-exchange": config.DeadExchange}
	if _, err := channel.QueueDeclare(config.ReturnQueue, true, false, false, false, deadArgs); err != nil {
		return err
	}
	if err := channel.QueueBind(config.ReturnQueue, loanReturnedKey, config.Exchange, false, nil); err != nil {
		return err
	}
	if _, err := channel.QueueDeclare(config.AckQueue, true, false, false, false, deadArgs); err != nil {
		return err
	}
	if err := channel.QueueBind(config.AckQueue, stockUpdatedKey, config.Exchange, false, nil); err != nil {
		return err
	}
	if _, err := channel.QueueDeclare(config.DeadQueue, true, false, false, false, nil); err != nil {
		return err
	}
	return channel.QueueBind(config.DeadQueue, "#", config.DeadExchange, false, nil)
}
