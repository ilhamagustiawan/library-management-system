package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	domainerrs "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *entity.User) error {
	const query = `
		INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return domainerrs.ErrConflict
	}
	return err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	const query = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users WHERE email = ?`
	return r.get(ctx, query, email)
}

func (r *Repository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	const query = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users WHERE id = ?`
	return r.get(ctx, query, id)
}

func (r *Repository) get(ctx context.Context, query string, arg any) (*entity.User, error) {
	var user entity.User
	if err := r.db.GetContext(ctx, &user, query, arg); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerrs.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
