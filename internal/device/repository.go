package device

import (
	"context"

	"github.com/arfan21/mertani/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data entity.Device) error
	GetByID(ctx context.Context, id string) (data entity.Device, err error)
	GetAll(ctx context.Context) (data []entity.Device, err error)
	Update(ctx context.Context, data entity.Device) error
	Delete(ctx context.Context, id string) error
}
