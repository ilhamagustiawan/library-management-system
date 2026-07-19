package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
	registrationusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/registration"
)

type fakeRegistrationUsecase struct{ calls int }

func (f *fakeRegistrationUsecase) Register(context.Context, registrationusecase.Input) (*entity.User, error) {
	f.calls++
	return &entity.User{ID: "user-123", Name: "Ada", Email: "ada@example.com", Role: entity.RoleMember}, nil
}

func TestRegisterRejectsRoleAndAcceptsStrictMemberInput(t *testing.T) {
	for _, test := range []struct {
		name   string
		body   string
		status int
		calls  int
	}{
		{name: "member", body: `{"name":"Ada","email":"ada@example.com","password":"correct horse battery staple"}`, status: 201, calls: 1},
		{name: "role", body: `{"name":"Ada","email":"ada@example.com","password":"correct horse battery staple","role":"admin"}`, status: 422},
		{name: "weak password", body: `{"name":"Ada","email":"ada@example.com","password":"short"}`, status: 422},
	} {
		t.Run(test.name, func(t *testing.T) {
			usecase := &fakeRegistrationUsecase{}
			handler := NewHandler(usecase, validator.New())
			app := fiber.New()
			app.Post("/api/v1/users", handler.Register)
			request := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			response, err := app.Test(request)
			if err != nil || response.StatusCode != test.status || usecase.calls != test.calls {
				t.Fatalf("status = %d, calls = %d, error = %v", response.StatusCode, usecase.calls, err)
			}
		})
	}
}
