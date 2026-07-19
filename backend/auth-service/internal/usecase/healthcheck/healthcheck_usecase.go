package healthcheck

import (
	"context"
	"fmt"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/repository"
)

type healthcheckUsecase struct {
	repository repository.HealthRepository
}

func NewUsecase(repository repository.HealthRepository) Usecase {
	return &healthcheckUsecase{repository: repository}
}

func (u *healthcheckUsecase) Readiness(ctx context.Context) error {
	if err := u.repository.Ping(ctx); err != nil {
		return fmt.Errorf("check database readiness: %w", err)
	}
	return nil
}
