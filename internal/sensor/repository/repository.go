package sensorrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/mertani/internal/entity"
	"github.com/arfan21/mertani/pkg/constant"
	dbpostgres "github.com/arfan21/mertani/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) Create(ctx context.Context, data entity.Sensor) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sensors (id, device_id, name, description, type)
		VALUES ($1, $2, $3, $4, $5)
	`, data.ID, data.DeviceID, data.Name, data.Description, data.Type)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID || pgxError.Code == constant.ErrSQLFKViolation {
				err = constant.ErrDeviceNotFound
			}
		}

		err = fmt.Errorf("sensor.repository.Create: failed to create sensor: %w", err)
		return err
	}

	return nil
}

func (r Repository) GetByID(ctx context.Context, id string) (data entity.Sensor, err error) {
	query := `
		SELECT id, device_id, name, description, type, created_at, updated_at
		FROM sensors
		WHERE id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(&data.ID, &data.DeviceID, &data.Name, &data.Description, &data.Type, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrSensorNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrSensorNotFound
			}
		}
		err = fmt.Errorf("sensor.repository.GetByID: failed to get sensor: %w", err)
		return
	}

	return
}

func (r Repository) Update(ctx context.Context, data entity.Sensor) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE sensors
		SET name = $1, description = $2, type = $3
		WHERE id = $4
	`, data.Name, data.Description, data.Type, data.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrSensorNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrSensorNotFound
			}
		}
		err = fmt.Errorf("sensor.repository.Update: failed to update sensor: %w", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		err = constant.ErrSensorNotFound
		err = fmt.Errorf("device.repository.Update: failed to update device: %w", err)
		return err
	}

	return nil
}

func (r Repository) Delete(ctx context.Context, id string) error {
	cmd, err := r.db.Exec(ctx, `
		DELETE FROM sensors
		WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrSensorNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrSensorNotFound
			}
		}
		err = fmt.Errorf("sensor.repository.Delete: failed to delete sensor: %w", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		err = constant.ErrSensorNotFound
		err = fmt.Errorf("device.repository.Update: failed to update device: %w", err)
		return err
	}

	return nil
}
