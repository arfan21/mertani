package sensorsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/arfan21/mertani/internal/entity"
	"github.com/arfan21/mertani/internal/model"
	"github.com/arfan21/mertani/internal/sensor"
	"github.com/arfan21/mertani/pkg/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo sensor.Repository
}

func New(repo sensor.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req model.SensorCreateRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("sensor.service.Create: failed to validate request: %w", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("device.service.Create: failed to generate id: %w", err)
		return
	}

	data := entity.Sensor{
		ID:          id.String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		DeviceID:    req.DeviceID,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("sensor.service.Create: failed to create sensor: %w", err)
		return
	}

	return
}

func (s Service) GetByID(ctx context.Context, id string) (res model.SensorResponse, err error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		err = fmt.Errorf("sensor.service.GetByID: failed to get sensor: %w", err)
		return
	}
	res = model.SensorResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Type:        data.Type,
		DeviceID:    data.DeviceID,
		CreatedAt:   data.CreatedAt.Format(time.DateTime),
		UpdatedAt:   data.UpdatedAt.Format(time.DateTime),
	}

	return
}

func (s Service) Update(ctx context.Context, req model.SensorUpdateRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("sensor.service.Update: failed to validate request: %w", err)
		return
	}

	data := entity.Sensor{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	err = s.repo.Update(ctx, data)
	if err != nil {
		err = fmt.Errorf("sensor.service.Update: failed to update sensor: %w", err)
		return
	}

	return
}

func (s Service) Delete(ctx context.Context, id string) (err error) {
	err = s.repo.Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("sensor.service.Delete: failed to delete sensor: %w", err)
		return
	}

	return
}
