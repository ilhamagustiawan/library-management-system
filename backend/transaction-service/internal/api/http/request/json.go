package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
)

type Borrow struct {
	BookID string `json:"bookId" validate:"required,uuid4"`
}

type Return struct {
	AcceptedFineAmountMinor *int64 `json:"acceptedFineAmountMinor" validate:"omitempty,gte=0"`
}

func DecodeStrictJSON(c *fiber.Ctx, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("request contains multiple JSON values")
	}
	return nil
}
