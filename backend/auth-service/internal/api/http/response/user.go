package response

import (
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type User struct {
	ID        string    `json:"id" example:"f81d4fae-7dec-11d0-a765-00a0c91e6bf6"`
	Name      string    `json:"name" example:"Ada Lovelace"`
	Email     string    `json:"email" example:"ada@example.com"`
	Role      string    `json:"role" example:"member"`
	CreatedAt time.Time `json:"createdAt" example:"2026-07-19T08:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2026-07-19T08:00:00Z"`
}

func NewUser(user *entity.User) User {
	return User{
		ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role.String(),
		CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	}
}
