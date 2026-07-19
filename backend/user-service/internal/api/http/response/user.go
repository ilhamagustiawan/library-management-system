package response

import (
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

type User struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      entity.Role `json:"role"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

func NewUser(user *entity.User) User {
	return User{
		ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role,
		CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	}
}
