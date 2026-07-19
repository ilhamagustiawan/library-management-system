package rabbitmq

import (
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

func TestMessageIsPersistentTypedJSONWithStableID(t *testing.T) {
	message := entity.OutboxMessage{OutboxEvent: entity.OutboxEvent{
		ID: "event-123", Type: entity.UserRegisteredV1, AggregateID: "user-123",
		Payload:    []byte(`{"eventId":"event-123","type":"UserRegistered.v1","data":{"userId":"user-123","role":"member"}}`),
		OccurredAt: time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC),
	}}
	publishing := newPublishing(message)
	if publishing.DeliveryMode != amqp.Persistent || publishing.ContentType != "application/json" {
		t.Fatalf("publishing = %#v", publishing)
	}
	if publishing.MessageId != "event-123" || publishing.Type != entity.UserRegisteredV1 || string(publishing.Body) != string(message.Payload) {
		t.Fatalf("publishing metadata = %#v", publishing)
	}
}
