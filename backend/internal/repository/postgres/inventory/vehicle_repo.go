package postgresinventory

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type VehicleRepo struct {
	pool *pgxpool.Pool
}

func NewVehicleRepo(pool *pgxpool.Pool) *VehicleRepo {
	return &VehicleRepo{pool: pool}
}

func (r *VehicleRepo) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *VehicleRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.Vehicle, error) {
	var v inventory.Vehicle
	err := r.pool.QueryRow(ctx,
		`SELECT id, type, license_plate, name, brand, model, year, status, capacity, created_at
		 FROM inventory.vehicles
		 WHERE id = $1`, id,
	).Scan(&v.ID, &v.Type, &v.LicensePlate, &v.Name, &v.Brand, &v.Model, &v.Year, &v.Status, &v.Capacity, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "vehicle", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get vehicle by id: %w", err)
	}
	return &v, nil
}

func (r *VehicleRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.Vehicle], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.vehicles`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count vehicles: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, type, license_plate, name, brand, model, year, status, capacity, created_at
		 FROM inventory.vehicles
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list vehicles: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.Vehicle, 0)
	for rows.Next() {
		var v inventory.Vehicle
		if err := rows.Scan(&v.ID, &v.Type, &v.LicensePlate, &v.Name, &v.Brand, &v.Model, &v.Year, &v.Status, &v.Capacity, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan vehicle: %w", err)
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.Vehicle]{Items: items, Total: total}, nil
}

func (r *VehicleRepo) Create(ctx context.Context, vehicle *inventory.Vehicle) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.vehicles (id, type, license_plate, name, brand, model, year, status, capacity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		vehicle.ID, vehicle.Type, vehicle.LicensePlate, vehicle.Name, vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Status, vehicle.Capacity, vehicle.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}
	return nil
}

func (r *VehicleRepo) Update(ctx context.Context, id uuid.UUID, vehicle *inventory.Vehicle) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.vehicles
		 SET type = $2, license_plate = $3, name = $4, brand = $5, model = $6, year = $7, status = $8, capacity = $9
		 WHERE id = $1`,
		id, vehicle.Type, vehicle.LicensePlate, vehicle.Name, vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Status, vehicle.Capacity,
	)
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}
	return nil
}

func (r *VehicleRepo) UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, vehicle *inventory.Vehicle) error {
	_, err := tx.Exec(ctx,
		`UPDATE inventory.vehicles
		 SET type = $2, license_plate = $3, name = $4, brand = $5, model = $6, year = $7, status = $8, capacity = $9
		 WHERE id = $1`,
		id, vehicle.Type, vehicle.LicensePlate, vehicle.Name, vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Status, vehicle.Capacity,
	)
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}
	return nil
}

func (r *VehicleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.vehicles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete vehicle: %w", err)
	}
	return nil
}
