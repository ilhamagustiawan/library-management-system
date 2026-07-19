package healthcheck

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	err error
}

func (r fakeRepository) Ping(context.Context) error { return r.err }

func TestReadinessReturnsRepositoryFailure(t *testing.T) {
	want := errors.New("database unavailable")
	uc := NewUsecase(fakeRepository{err: want})

	err := uc.Readiness(context.Background())
	if !errors.Is(err, want) {
		t.Fatalf("Readiness() error = %v, want wrapped repository error", err)
	}
}
