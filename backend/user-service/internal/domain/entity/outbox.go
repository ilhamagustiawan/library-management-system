package entity

import "time"

const UserRegisteredV1 = "UserRegistered.v1"

type UserRegistered struct {
	EventID    string             `json:"eventId"`
	Type       string             `json:"type"`
	OccurredAt time.Time          `json:"occurredAt"`
	Data       UserRegisteredData `json:"data"`
}

type UserRegisteredData struct {
	UserID string `json:"userId"`
	Role   Role   `json:"role"`
}

type OutboxEvent struct {
	ID          string    `db:"id"`
	Type        string    `db:"event_type"`
	AggregateID string    `db:"aggregate_id"`
	Payload     []byte    `db:"payload"`
	OccurredAt  time.Time `db:"occurred_at"`
}

type OutboxMessage struct {
	OutboxEvent
	Attempts int `db:"attempts"`
}
