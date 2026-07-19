package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/gofiber/fiber/v2"
)

func DecodeStrictJSON(ctx *fiber.Ctx, output any) error {
	if !ctx.Is("json") {
		return errors.New("content type must be application/json")
	}
	decoder := json.NewDecoder(bytes.NewReader(ctx.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(output); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}
