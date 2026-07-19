package entity

import "time"

type RegistrationStatus string

const (
	RegistrationPending   RegistrationStatus = "pending"
	RegistrationCompleted RegistrationStatus = "completed"
	RegistrationConflict  RegistrationStatus = "conflict"
)

type Registration struct {
	ID        string             `db:"id"`
	Name      string             `db:"name"`
	Email     string             `db:"email"`
	Status    RegistrationStatus `db:"status"`
	CreatedAt time.Time          `db:"created_at"`
	UpdatedAt time.Time          `db:"updated_at"`
}
