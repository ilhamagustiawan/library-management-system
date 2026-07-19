package healthcheck

import (
	"context"
	"fmt"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/repository"
)

type Usecase interface{ Readiness(context.Context) error }

type usecase struct{ repository repository.HealthRepository }

func New(repository repository.HealthRepository) Usecase { return &usecase{repository: repository} }

func (u *usecase) Readiness(ctx context.Context) error {
	if err := u.repository.Ping(ctx); err != nil {
		return fmt.Errorf("check database readiness: %w", err)
	}
	return nil
}
