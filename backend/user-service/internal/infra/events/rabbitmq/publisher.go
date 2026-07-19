package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

type Config struct {
	URL            string
	Exchange       string
	RoutingKey     string
	Queue          string
	ConfirmTimeout time.Duration
}

type Publisher struct {
	config  Config
	mu      sync.Mutex
	conn    *amqp.Connection
	ch      *amqp.Channel
	returns <-chan amqp.Return
}

func NewPublisher(config Config) *Publisher {
	if config.Exchange == "" {
		config.Exchange = "library.events"
	}
	if config.RoutingKey == "" {
		config.RoutingKey = "user.registered.v1"
	}
	if config.Queue == "" {
		config.Queue = "library.user-registered.v1"
	}
	if config.ConfirmTimeout <= 0 {
		config.ConfirmTimeout = 5 * time.Second
	}
	return &Publisher{config: config}
}

func (p *Publisher) Publish(ctx context.Context, message entity.OutboxMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ensureChannel(); err != nil {
		return err
	}

	confirmation, err := p.ch.PublishWithDeferredConfirmWithContext(
		ctx, p.config.Exchange, p.config.RoutingKey, true, false, newPublishing(message),
	)
	if err != nil {
		p.reset()
		return fmt.Errorf("publish event %s: %w", message.ID, err)
	}
	if confirmation == nil {
		p.reset()
		return fmt.Errorf("publish event %s: publisher confirmation unavailable", message.ID)
	}
	confirmCtx, cancel := context.WithTimeout(ctx, p.config.ConfirmTimeout)
	defer cancel()
	acknowledged, err := confirmation.WaitContext(confirmCtx)
	if err != nil || !acknowledged {
		p.reset()
		return fmt.Errorf("confirm event %s: acknowledged=%t: %w", message.ID, acknowledged, err)
	}
	select {
	case returned := <-p.returns:
		p.reset()
		return fmt.Errorf("route event %s: %d %s", message.ID, returned.ReplyCode, returned.ReplyText)
	default:
		return nil
	}
}

func (p *Publisher) ensureChannel() error {
	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil && !p.ch.IsClosed() {
		return nil
	}
	p.reset()
	conn, err := amqp.Dial(p.config.URL)
	if err != nil {
		return fmt.Errorf("connect RabbitMQ: %w", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("open RabbitMQ channel: %w", err)
	}
	fail := func(operation string, err error) error {
		_ = channel.Close()
		_ = conn.Close()
		return fmt.Errorf("%s: %w", operation, err)
	}
	if err := channel.ExchangeDeclare(p.config.Exchange, "topic", true, false, false, false, nil); err != nil {
		return fail("declare event exchange", err)
	}
	if _, err := channel.QueueDeclare(p.config.Queue, true, false, false, false, amqp.Table{"x-queue-type": "quorum"}); err != nil {
		return fail("declare event queue", err)
	}
	if err := channel.QueueBind(p.config.Queue, p.config.RoutingKey, p.config.Exchange, false, nil); err != nil {
		return fail("bind event queue", err)
	}
	if err := channel.Confirm(false); err != nil {
		return fail("enable publisher confirms", err)
	}
	p.conn = conn
	p.ch = channel
	p.returns = channel.NotifyReturn(make(chan amqp.Return, 1))
	return nil
}

func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	var result error
	if p.ch != nil && !p.ch.IsClosed() {
		result = p.ch.Close()
	}
	if p.conn != nil && !p.conn.IsClosed() {
		result = errors.Join(result, p.conn.Close())
	}
	p.ch = nil
	p.conn = nil
	p.returns = nil
	return result
}

func (p *Publisher) reset() {
	if p.ch != nil && !p.ch.IsClosed() {
		_ = p.ch.Close()
	}
	if p.conn != nil && !p.conn.IsClosed() {
		_ = p.conn.Close()
	}
	p.ch = nil
	p.conn = nil
	p.returns = nil
}

func newPublishing(message entity.OutboxMessage) amqp.Publishing {
	return amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		MessageId:    message.ID,
		Type:         message.Type,
		Timestamp:    message.OccurredAt,
		Body:         append([]byte(nil), message.Payload...),
	}
}
