package sensor

import (
	"context"

	"github.com/arfan21/mertani/internal/model"
)

type Service interface {
	Create(ctx context.Context, req model.SensorCreateRequest) (err error)
	GetByID(ctx context.Context, id string) (res model.SensorResponse, err error)
	Update(ctx context.Context, req model.SensorUpdateRequest) (err error)
	Delete(ctx context.Context, id string) (err error)
}
