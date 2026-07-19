package book

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
)

type staticToken string

func (t staticToken) Token(context.Context) (string, error) { return string(t), nil }

func TestReserveCallsAtomicBookEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut || request.URL.Path != "/internal/v1/books/book-1/reservations/loan-1" {
			t.Fatalf("request = %s %s", request.Method, request.URL.Path)
		}
		if request.Header.Get("Authorization") != "Bearer service-token" {
			t.Fatalf("authorization missing")
		}
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusCreated)
		_, _ = response.Write([]byte(`{"code":"LMS-200000","data":{"transactionId":"loan-1","bookId":"book-1","status":"active"}}`))
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL}, staticToken("service-token"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if err := client.Reserve(context.Background(), "book-1", "loan-1"); err != nil {
		t.Fatalf("Reserve() error = %v", err)
	}
}

func TestReserveMapsUnavailableStock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusConflict)
		_, _ = response.Write([]byte(`{"code":"LMS-409003","message":"book has no available copies"}`))
	}))
	defer server.Close()
	client, _ := NewClient(Config{BaseURL: server.URL}, staticToken("service-token"))
	if err := client.Reserve(context.Background(), "book-1", "loan-1"); !errors.Is(err, errs.ErrStockUnavailable) {
		t.Fatalf("Reserve() error = %v, want stock unavailable", err)
	}
}
