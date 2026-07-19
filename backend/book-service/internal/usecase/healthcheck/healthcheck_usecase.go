package healthcheck

import (
	"context"
	"fmt"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

type Usecase struct{ repository repository.HealthRepository }

func NewUsecase(repository repository.HealthRepository) *Usecase {
	return &Usecase{repository: repository}
}
func (u *Usecase) Liveness(context.Context) error { return nil }
func (u *Usecase) Readiness(ctx context.Context) error {
	if err := u.repository.Ping(ctx); err != nil {
		return fmt.Errorf("ping book database: %w", err)
	}
	return nil
}
