package auth

import (
	"crypto/rand"
	"encoding/base64"
)

type TokenGenerator struct {
	bytes int
}

func NewTokenGenerator(bytes int) *TokenGenerator {
	return &TokenGenerator{bytes: bytes}
}

func (g *TokenGenerator) Generate() (string, error) {
	buffer := make([]byte, g.bytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
