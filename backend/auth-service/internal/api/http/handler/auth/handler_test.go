package auth

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	authusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/auth"
)

type fakeUsecase struct{}

func (fakeUsecase) Register(context.Context, authusecase.RegisterInput) (*entity.User, error) {
	return &entity.User{ID: "user-1", Name: "Ada", Email: "ada@example.com"}, nil
}
func (fakeUsecase) Login(context.Context, authusecase.LoginInput) (*authusecase.LoginResult, error) {
	return &authusecase.LoginResult{
		User:         &entity.User{ID: "user-1", Name: "Ada", Email: "ada@example.com"},
		SessionToken: "opaque-session", ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}
func (fakeUsecase) AuthenticateSession(_ context.Context, token string) (*entity.User, error) {
	if token == "opaque-session" {
		return &entity.User{ID: "user-1", Name: "Ada", Email: "ada@example.com"}, nil
	}
	return nil, context.Canceled
}
func (fakeUsecase) Logout(context.Context, string) error { return nil }

func TestRegisterRejectsInvalidInput(t *testing.T) {
	app := testApp()
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"name":"Ada","email":"not-email","password":"short"}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnprocessableEntity {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("status = %d, want 422; body=%s", response.StatusCode, body)
	}
}

func TestLoginSetsSessionCookieAndRedirectsOnlyToLocalAuthorizeEndpoint(t *testing.T) {
	app := testApp()
	form := "email=ada%40example.com&password=correct+horse+battery&return_to=" +
		"http%3A%2F%2Fauth.example%2Foauth%2Fauthorize%3Fclient_id%3Dnextjs"
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", response.StatusCode)
	}
	if response.Header.Get("Location") != "http://auth.example/oauth/authorize?client_id=nextjs" {
		t.Fatalf("Location = %q, want local authorize URL", response.Header.Get("Location"))
	}
	cookies := response.Cookies()
	if len(cookies) != 1 || cookies[0].Name != "lms_session" || !cookies[0].HttpOnly {
		t.Fatalf("cookies = %#v, want HttpOnly session cookie", cookies)
	}
}

func TestLoginRejectsExternalReturnURL(t *testing.T) {
	app := testApp()
	form := "email=ada%40example.com&password=correct+horse+battery&return_to=https%3A%2F%2Fevil.example%2Fsteal"
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", response.StatusCode)
	}
}

func testApp() *fiber.App {
	handler := NewHandler(fakeUsecase{}, validator.New(), Config{
		Issuer: "http://auth.example", SessionCookieName: "lms_session", SessionCookieSecure: true,
	})
	app := fiber.New()
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)
	return app
}
