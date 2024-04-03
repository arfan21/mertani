package sensor

import (
	"context"

	"github.com/arfan21/mertani/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data entity.Sensor) error
	GetByID(ctx context.Context, id string) (data entity.Sensor, err error)
	Update(ctx context.Context, data entity.Sensor) error
	Delete(ctx context.Context, id string) error
}
