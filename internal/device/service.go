package device

import (
	"context"

	"github.com/arfan21/mertani/internal/model"
)

type Service interface {
	Create(ctx context.Context, req model.DeviceCreateRequest) (err error)
	GetByID(ctx context.Context, id string) (res model.DeviceResponse, err error)
	GetAll(ctx context.Context) (res []model.DeviceResponse, err error)
	Update(ctx context.Context, req model.DeviceUpdateRequest) (err error)
	Delete(ctx context.Context, id string) (err error)
}
