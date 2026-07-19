package identity

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	identityusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/identity"
)

type fakeCreator struct{ calls int }

func (f *fakeCreator) Create(_ context.Context, _ string, _ identityusecase.Input) (*entity.User, error) {
	f.calls++
	return &entity.User{ID: "user-123", Name: "Ada", Email: "ada@example.com", Role: entity.RoleMember}, nil
}

type fakeTokenAuthenticator struct{ err error }

func (f fakeTokenAuthenticator) AuthenticateServiceToken(context.Context, string, string, string, string) error {
	return f.err
}

func TestCreateRequiresUserServiceGrantAndRejectsRoleInput(t *testing.T) {
	for _, test := range []struct {
		name       string
		authorizer fakeTokenAuthenticator
		body       string
		status     int
	}{
		{name: "valid", body: `{"name":"Ada","email":"ada@example.com","password":"correct horse battery staple"}`, status: 201},
		{name: "invalid token", authorizer: fakeTokenAuthenticator{err: errors.New("invalid")}, body: `{"name":"Ada","email":"ada@example.com","password":"correct horse battery staple"}`, status: 401},
		{name: "role injection", body: `{"name":"Ada","email":"ada@example.com","password":"correct horse battery staple","role":"admin"}`, status: 422},
	} {
		t.Run(test.name, func(t *testing.T) {
			creator := &fakeCreator{}
			handler := NewHandler(creator, test.authorizer, validator.New())
			app := fiber.New()
			app.Post("/internal/identities", handler.Create)
			request := httptest.NewRequest("POST", "/internal/identities", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer service-token")
			request.Header.Set("Idempotency-Key", "registration-123")
			response, err := app.Test(request)
			if err != nil || response.StatusCode != test.status {
				t.Fatalf("status = %d, error = %v", response.StatusCode, err)
			}
			if test.status == 201 {
				var payload map[string]any
				if err := json.NewDecoder(response.Body).Decode(&payload); err != nil || creator.calls != 1 {
					t.Fatalf("response = %#v, calls = %d, error = %v", payload, creator.calls, err)
				}
			}
		})
	}
}
