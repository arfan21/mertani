package user

import (
	"context"

	"github.com/arfan21/mertani/internal/entity"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) Repository

	Create(ctx context.Context, data entity.User) (err error)
	GetByEmail(ctx context.Context, email string) (data entity.User, err error)
}
