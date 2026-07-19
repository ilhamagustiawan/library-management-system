package helper

import (
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SessionCookieConfig struct {
	Name   string
	Domain string
	Secure bool
}

func SetSessionCookie(c *fiber.Ctx, config SessionCookieConfig, value string, expiresAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name: config.Name, Value: value, Path: "/", Domain: config.Domain,
		Expires: expiresAt, HTTPOnly: true, Secure: config.Secure, SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func IsSafeAuthorizeURL(raw, issuerRaw string) bool {
	issuer, err := url.Parse(issuerRaw)
	if err != nil {
		return false
	}
	target, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return target.Scheme == issuer.Scheme && target.Host == issuer.Host &&
		target.Path == "/oauth/authorize" && target.User == nil && target.Fragment == ""
}

func BearerToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}
