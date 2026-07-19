package healthcheck

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct{ database *sqlx.DB }

func NewRepository(database *sqlx.DB) *Repository    { return &Repository{database: database} }
func (r *Repository) Ping(ctx context.Context) error { return r.database.PingContext(ctx) }
