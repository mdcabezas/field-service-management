package postgresinventory

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/repository"
)

type VehicleTypeRepo struct {
	pool *pgxpool.Pool
}

func NewVehicleTypeRepo(pool *pgxpool.Pool) *VehicleTypeRepo {
	return &VehicleTypeRepo{pool: pool}
}

func (r *VehicleTypeRepo) List(ctx context.Context) (*repository.ListResult[inventory.VehicleType], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.vehicle_types WHERE active = true`).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, code, name, active
		 FROM inventory.vehicle_types
		 WHERE active = true
		 ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]inventory.VehicleType, 0)
	for rows.Next() {
		var vt inventory.VehicleType
		if err := rows.Scan(&vt.ID, &vt.Code, &vt.Name, &vt.Active); err != nil {
			return nil, err
		}
		items = append(items, vt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &repository.ListResult[inventory.VehicleType]{Items: items, Total: total}, nil
}

func (r *VehicleTypeRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.VehicleType, error) {
	var vt inventory.VehicleType
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, active
		 FROM inventory.vehicle_types
		 WHERE id = $1`, id,
	).Scan(&vt.ID, &vt.Code, &vt.Name, &vt.Active)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &vt, nil
}
