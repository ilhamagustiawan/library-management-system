package transaction

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/middleware"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	transactionusecase "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/usecase/transaction"
)

type fakeUsecase struct {
	confirmed       bool
	returnInputSeen transactionusecase.ReturnInput
}

func (f *fakeUsecase) Borrow(context.Context, transactionusecase.BorrowInput) (*entity.Loan, error) {
	return &entity.Loan{ID: "loan-1"}, nil
}
func (f *fakeUsecase) QuoteReturn(_ context.Context, input transactionusecase.ReturnQuoteInput) (*entity.ReturnQuote, error) {
	return &entity.ReturnQuote{
		LoanID: input.LoanID, BookID: "book-1", Fine: &entity.FineQuote{
			OverdueDays: 2, DailyRateMinor: 5000, TotalAmountMinor: 10000, Currency: "IDR",
		},
	}, nil
}
func (f *fakeUsecase) Return(_ context.Context, input transactionusecase.ReturnInput) (*entity.Loan, bool, error) {
	f.returnInputSeen = input
	return &entity.Loan{ID: input.LoanID, StockSyncStatus: entity.StockSyncPending}, f.confirmed, nil
}
func (f *fakeUsecase) List(context.Context, transactionusecase.ListInput) (transactionusecase.Page, error) {
	return transactionusecase.Page{}, nil
}

func TestReturnUsesAcceptedWhileStockAckPending(t *testing.T) {
	app := fiber.New()
	handler := NewHandler(&fakeUsecase{}, validator.New())
	app.Post("/loans/:loanId/return", middleware.RequireScope("loans:return:self"), handler.ReturnSelf)
	request, _ := http.NewRequest(http.MethodPost, "/loans/52a88672-a4c2-4876-be5a-65863aeb35e4/return", nil)
	request.Header.Set("X-Credential-Sub", "member-1")
	request.Header.Set("X-Credential-Scope", "loans:return:self")
	result, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if result.StatusCode != http.StatusAccepted || result.Header.Get("Retry-After") != "2" {
		t.Fatalf("response = %d retry %q", result.StatusCode, result.Header.Get("Retry-After"))
	}
}

func TestReturnQuoteShowsAuthoritativeFine(t *testing.T) {
	app := fiber.New()
	handler := NewHandler(&fakeUsecase{}, validator.New())
	app.Get("/loans/:loanId/return", middleware.RequireScope("loans:return:self"), handler.QuoteReturnSelf)
	request, _ := http.NewRequest(http.MethodGet, "/loans/52a88672-a4c2-4876-be5a-65863aeb35e4/return", nil)
	request.Header.Set("X-Credential-Sub", "member-1")
	request.Header.Set("X-Credential-Scope", "loans:return:self")

	result, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer result.Body.Close()
	body, readErr := io.ReadAll(result.Body)
	if readErr != nil {
		t.Fatalf("read response: %v", readErr)
	}
	if result.StatusCode != http.StatusOK || !strings.Contains(string(body), `"totalAmountMinor":10000`) {
		t.Fatalf("response = %d %s", result.StatusCode, body)
	}
}

func TestReturnAcceptsFineAmountReviewedByMember(t *testing.T) {
	app := fiber.New()
	usecase := &fakeUsecase{}
	handler := NewHandler(usecase, validator.New())
	app.Post("/loans/:loanId/return", middleware.RequireScope("loans:return:self"), handler.ReturnSelf)
	request, _ := http.NewRequest(http.MethodPost, "/loans/52a88672-a4c2-4876-be5a-65863aeb35e4/return", strings.NewReader(`{"acceptedFineAmountMinor":10000}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Credential-Sub", "member-1")
	request.Header.Set("X-Credential-Scope", "loans:return:self")

	result, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if result.StatusCode != http.StatusAccepted {
		var body map[string]any
		_ = json.NewDecoder(result.Body).Decode(&body)
		t.Fatalf("response = %d %#v", result.StatusCode, body)
	}
	if usecase.returnInputSeen.AcceptedFineAmountMinor == nil || *usecase.returnInputSeen.AcceptedFineAmountMinor != 10000 {
		t.Fatalf("accepted fine = %v, want 10000", usecase.returnInputSeen.AcceptedFineAmountMinor)
	}
}

func TestBorrowRejectsUnknownInputFields(t *testing.T) {
	app := fiber.New()
	handler := NewHandler(&fakeUsecase{}, validator.New())
	app.Post("/loans", middleware.RequireScope("loans:borrow:self"), handler.Borrow)
	request, _ := http.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{"bookId":"7b36fe43-f31d-4861-884f-42ed7386b1e9","memberId":"attacker"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Credential-Sub", "member-1")
	request.Header.Set("X-Credential-Scope", "loans:borrow:self")
	result, _ := app.Test(request)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d", result.StatusCode)
	}
}

func TestBorrowRequiresExactScope(t *testing.T) {
	app := fiber.New()
	handler := NewHandler(&fakeUsecase{}, validator.New())
	app.Post("/loans", middleware.RequireScope("loans:borrow:self"), handler.Borrow)
	request, _ := http.NewRequest(http.MethodPost, "/loans", strings.NewReader(`{"bookId":"7b36fe43-f31d-4861-884f-42ed7386b1e9"}`))
	request.Header.Set("X-Credential-Sub", "member-1")
	request.Header.Set("X-Credential-Scope", "transactions:read:self")
	result, _ := app.Test(request)
	if result.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d", result.StatusCode)
	}
}
