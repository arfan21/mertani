package devicesvc

import (
	"context"
	"fmt"
	"time"

	"github.com/arfan21/mertani/internal/device"
	"github.com/arfan21/mertani/internal/entity"
	"github.com/arfan21/mertani/internal/model"
	"github.com/arfan21/mertani/pkg/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo device.Repository
}

func New(repo device.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req model.DeviceCreateRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("device.service.Create: failed to validate request: %w", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("device.service.Create: failed to generate id: %w", err)
		return
	}

	data := entity.Device{
		ID:          id.String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Location:    req.Location,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("device.service.Create: failed to create device: %w", err)
		return
	}

	return
}

func (s Service) GetByID(ctx context.Context, id string) (res model.DeviceResponse, err error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		err = fmt.Errorf("device.service.GetByID: failed to get device: %w", err)
		return
	}
	res = model.DeviceResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Type:        data.Type,
		Location:    data.Location,
		CreatedAt:   data.CreatedAt.Format(time.DateTime),
		UpdatedAt:   data.UpdatedAt.Format(time.DateTime),
	}

	res.Sensors = make([]model.SensorResponse, len(data.Sensors))

	for i, s := range data.Sensors {
		res.Sensors[i] = model.SensorResponse{
			ID:          s.ID,
			DeviceID:    s.DeviceID,
			Name:        s.Name,
			Description: s.Description,
			Type:        s.Type,
			CreatedAt:   s.CreatedAt.Format(time.DateTime),
			UpdatedAt:   s.UpdatedAt.Format(time.DateTime),
		}
	}

	return
}

func (s Service) GetAll(ctx context.Context) (res []model.DeviceResponse, err error) {
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		err = fmt.Errorf("device.service.GetAll: failed to get devices: %w", err)
		return
	}

	res = make([]model.DeviceResponse, len(data))
	for i, d := range data {
		res[i] = model.DeviceResponse{
			ID:          d.ID,
			Name:        d.Name,
			Description: d.Description,
			Type:        d.Type,
			Location:    d.Location,
			CreatedAt:   d.CreatedAt.Format(time.DateTime),
			UpdatedAt:   d.UpdatedAt.Format(time.DateTime),
		}
	}

	return
}

func (s Service) Update(ctx context.Context, req model.DeviceUpdateRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("device.service.Update: failed to validate request: %w", err)
		return
	}

	data := entity.Device{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Location:    req.Location,
	}

	err = s.repo.Update(ctx, data)
	if err != nil {
		err = fmt.Errorf("device.service.Update: failed to update device: %w", err)
		return
	}

	return
}

func (s Service) Delete(ctx context.Context, id string) (err error) {
	err = s.repo.Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("device.service.Delete: failed to delete device: %w", err)
		return
	}

	return
}
