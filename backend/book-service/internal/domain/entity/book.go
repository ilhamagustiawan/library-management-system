package entity

import "time"

type Book struct {
	ID              string     `json:"id" db:"id"`
	ISBN            string     `json:"isbn" db:"isbn"`
	Title           string     `json:"title" db:"title"`
	Author          string     `json:"author" db:"author"`
	Description     *string    `json:"description" db:"description"`
	CoverURL        *string    `json:"coverUrl" db:"cover_url"`
	PublicationYear *int       `json:"publicationYear" db:"publication_year"`
	TotalCopies     int        `json:"totalCopies" db:"total_copies"`
	AvailableCopies int        `json:"availableCopies" db:"available_copies"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
	ArchivedAt      *time.Time `json:"-" db:"archived_at"`
}

type Change[T any] struct {
	Set   bool
	Value *T
}

type BookPatch struct {
	ISBN            Change[string]
	Title           Change[string]
	Author          Change[string]
	Description     Change[string]
	CoverURL        Change[string]
	PublicationYear Change[int]
	TotalCopies     Change[int]
}

type ReservationStatus string

const (
	ReservationActive   ReservationStatus = "active"
	ReservationReleased ReservationStatus = "released"
)

type Reservation struct {
	TransactionID string            `json:"transactionId" db:"transaction_id"`
	BookID        string            `json:"bookId" db:"book_id"`
	Status        ReservationStatus `json:"status" db:"status"`
	CreatedAt     time.Time         `json:"createdAt" db:"created_at"`
	ReleasedAt    *time.Time        `json:"releasedAt,omitempty" db:"released_at"`
}

type Stock struct {
	BookID          string `json:"bookId" db:"book_id"`
	TotalCopies     int    `json:"totalCopies" db:"total_copies"`
	AvailableCopies int    `json:"availableCopies" db:"available_copies"`
}
