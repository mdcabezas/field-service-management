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

type MaintenanceScheduleRepo struct {
	pool *pgxpool.Pool
}

func NewMaintenanceScheduleRepo(pool *pgxpool.Pool) *MaintenanceScheduleRepo {
	return &MaintenanceScheduleRepo{pool: pool}
}

func (r *MaintenanceScheduleRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.MaintenanceSchedule, error) {
	var m inventory.MaintenanceSchedule
	err := r.pool.QueryRow(ctx,
		`SELECT id, type, reference_id, frequency_km, frequency_days, last_service_date, last_service_km, active, created_at
		 FROM inventory.maintenance_schedules
		 WHERE id = $1`, id,
	).Scan(&m.ID, &m.Type, &m.ReferenceID, &m.FrequencyKM, &m.FrequencyDays, &m.LastServiceDate, &m.LastServiceKM, &m.Active, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "maintenance_schedule", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get maintenance_schedule by id: %w", err)
	}
	return &m, nil
}

func (r *MaintenanceScheduleRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.MaintenanceSchedule], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.maintenance_schedules`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count maintenance schedules: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, type, reference_id, frequency_km, frequency_days, last_service_date, last_service_km, active, created_at
		 FROM inventory.maintenance_schedules
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list maintenance schedules: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.MaintenanceSchedule, 0)
	for rows.Next() {
		var m inventory.MaintenanceSchedule
		if err := rows.Scan(&m.ID, &m.Type, &m.ReferenceID, &m.FrequencyKM, &m.FrequencyDays, &m.LastServiceDate, &m.LastServiceKM, &m.Active, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan maintenance schedule: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.MaintenanceSchedule]{Items: items, Total: total}, nil
}

func (r *MaintenanceScheduleRepo) Create(ctx context.Context, sched *inventory.MaintenanceSchedule) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.maintenance_schedules (id, type, reference_id, frequency_km, frequency_days, last_service_date, last_service_km, active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		sched.ID, sched.Type, sched.ReferenceID, sched.FrequencyKM, sched.FrequencyDays, sched.LastServiceDate, sched.LastServiceKM, sched.Active, sched.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create maintenance schedule: %w", err)
	}
	return nil
}

func (r *MaintenanceScheduleRepo) Update(ctx context.Context, id uuid.UUID, sched *inventory.MaintenanceSchedule) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.maintenance_schedules
		 SET type = $2, reference_id = $3, frequency_km = $4, frequency_days = $5, last_service_date = $6, last_service_km = $7, active = $8
		 WHERE id = $1`,
		id, sched.Type, sched.ReferenceID, sched.FrequencyKM, sched.FrequencyDays, sched.LastServiceDate, sched.LastServiceKM, sched.Active,
	)
	if err != nil {
		return fmt.Errorf("update maintenance schedule: %w", err)
	}
	return nil
}

func (r *MaintenanceScheduleRepo) ListActive(ctx context.Context) ([]inventory.MaintenanceSchedule, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, type, reference_id, frequency_km, frequency_days, last_service_date, last_service_km, active, created_at
		 FROM inventory.maintenance_schedules
		 WHERE active = true AND frequency_days IS NOT NULL
		 ORDER BY created_at DESC
		 LIMIT 1000`,
	)
	if err != nil {
		return nil, fmt.Errorf("list active maintenance schedules: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.MaintenanceSchedule, 0)
	for rows.Next() {
		var m inventory.MaintenanceSchedule
		if err := rows.Scan(&m.ID, &m.Type, &m.ReferenceID, &m.FrequencyKM, &m.FrequencyDays, &m.LastServiceDate, &m.LastServiceKM, &m.Active, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan maintenance schedule: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *MaintenanceScheduleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.maintenance_schedules WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete maintenance schedule: %w", err)
	}
	return nil
}
