package registration

import (
	"context"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

type Input struct {
	Name     string
	Email    string
	Password string
}

type IdentityInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Identity struct {
	ID    string
	Name  string
	Email string
	Role  entity.Role
}

type IdentityCreator interface {
	Create(ctx context.Context, idempotencyKey string, input IdentityInput) (*Identity, error)
}

type Usecase interface {
	Register(ctx context.Context, input Input) (*entity.User, error)
}
