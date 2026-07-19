package request

import (
	"bytes"
	"encoding/json"
)

type CreateBook struct {
	ISBN            string  `json:"isbn" validate:"required,max=32" example:"9780132350884"`
	Title           string  `json:"title" validate:"required,max=255" example:"Clean Code"`
	Author          string  `json:"author" validate:"required,max=255" example:"Robert C. Martin"`
	Description     *string `json:"description" validate:"omitempty,max=5000" example:"A handbook of agile software craftsmanship."`
	CoverURL        *string `json:"coverUrl" validate:"omitempty,max=512" example:"https://covers.openlibrary.org/b/id/10521270-M.jpg"`
	PublicationYear *int    `json:"publicationYear" validate:"omitempty,min=1000" example:"2008"`
	TotalCopies     int     `json:"totalCopies" validate:"required,min=1" example:"3"`
}

type OptionalString struct {
	Present bool
	Value   *string
}

func (value *OptionalString) UnmarshalJSON(payload []byte) error {
	value.Present = true
	if bytes.Equal(payload, []byte("null")) {
		value.Value = nil
		return nil
	}
	return json.Unmarshal(payload, &value.Value)
}

type OptionalInt struct {
	Present bool
	Value   *int
}

func (value *OptionalInt) UnmarshalJSON(payload []byte) error {
	value.Present = true
	if bytes.Equal(payload, []byte("null")) {
		value.Value = nil
		return nil
	}
	return json.Unmarshal(payload, &value.Value)
}

type UpdateBook struct {
	ISBN            OptionalString `json:"isbn" swaggertype:"string" example:"9780132350884"`
	Title           OptionalString `json:"title" swaggertype:"string" example:"Clean Code"`
	Author          OptionalString `json:"author" swaggertype:"string" example:"Robert C. Martin"`
	Description     OptionalString `json:"description" swaggertype:"string" example:"Updated description"`
	CoverURL        OptionalString `json:"coverUrl" swaggertype:"string" example:"https://covers.openlibrary.org/b/id/10521270-M.jpg"`
	PublicationYear OptionalInt    `json:"publicationYear" swaggertype:"integer" example:"2009"`
	TotalCopies     OptionalInt    `json:"totalCopies" swaggertype:"integer" example:"4"`
}
