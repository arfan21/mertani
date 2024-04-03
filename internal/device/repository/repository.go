package devicerepo

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

func (r Repository) Create(ctx context.Context, data entity.Device) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO devices (id, name, description, type, location)
		VALUES ($1, $2, $3, $4, $5)
	`, data.ID, data.Name, data.Description, data.Type, data.Location)
	if err != nil {
		err = fmt.Errorf("device.repository.Create: failed to create device: %w", err)
		return err
	}

	return nil
}

func (r Repository) GetByID(ctx context.Context, id string) (data entity.Device, err error) {
	// 2006-01-02T15:04:05Z07:00
	query := `
		SELECT 
			dv.id, dv.name, dv.description, dv.type, dv.location, dv.created_at, dv.updated_at,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', ss.id,
						'name', ss.name,
						'description', ss.description,
						'type', ss.type,
						'created_at', to_char(ss.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
						'updated_at', to_char(ss.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
					) 
				) FILTER (WHERE ss.id IS NOT NULL),
				NULL
			) AS sensors
		FROM devices dv
		LEFT JOIN sensors ss ON dv.id = ss.device_id
		WHERE dv.id = $1
		GROUP BY dv.id
	`

	err = r.db.QueryRow(ctx, query, id).Scan(
		&data.ID,
		&data.Name,
		&data.Description,
		&data.Type,
		&data.Location,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.Sensors,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrDeviceNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrDeviceNotFound
			}
		}

		err = fmt.Errorf("device.repository.GetByID: failed to get device: %w", err)
	}

	return
}

func (r Repository) GetAll(ctx context.Context) (data []entity.Device, err error) {
	query := `
		SELECT id, name, description, type, location, created_at, updated_at
		FROM devices
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("device.repository.GetAll: failed to get devices: %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d entity.Device
		err = rows.Scan(&d.ID, &d.Name, &d.Description, &d.Type, &d.Location, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			err = fmt.Errorf("device.repository.GetAll: failed to scan devices: %w", err)
			return
		}

		data = append(data, d)
	}

	return
}

func (r Repository) Update(ctx context.Context, data entity.Device) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE devices
		SET name = $1, description = $2, type = $3, location = $4, updated_at = $5
		WHERE id = $6
	`, data.Name, data.Description, data.Type, data.Location, data.UpdatedAt, data.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrDeviceNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrDeviceNotFound
			}
		}

		err = fmt.Errorf("device.repository.Update: failed to update device: %w", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		err = constant.ErrDeviceNotFound
		err = fmt.Errorf("device.repository.Update: failed to update device: %w", err)
		return err
	}

	return nil
}

func (r Repository) Delete(ctx context.Context, id string) error {
	cmd, err := r.db.Exec(ctx, `
		DELETE FROM devices
		WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrDeviceNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrDeviceNotFound
			}
		}

		err = fmt.Errorf("device.repository.Delete: failed to delete device: %w", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		err = constant.ErrDeviceNotFound
		err = fmt.Errorf("device.repository.Update: failed to update device: %w", err)
		return err
	}

	return nil
}
