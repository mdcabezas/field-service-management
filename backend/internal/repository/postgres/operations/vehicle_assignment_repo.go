package postgresoperations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type VehicleAssignmentRepo struct {
	pool *pgxpool.Pool
}

func NewVehicleAssignmentRepo(pool *pgxpool.Pool) *VehicleAssignmentRepo {
	return &VehicleAssignmentRepo{pool: pool}
}

func (r *VehicleAssignmentRepo) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *VehicleAssignmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VehicleAssignment, error) {
	var va operations.VehicleAssignment
	err := r.pool.QueryRow(ctx,
		`SELECT id, vehicle_id, daily_plan_id, visit_id, departure_time, return_time, departure_mileage, return_mileage, notes, created_at
		 FROM operations.vehicle_assignments
		 WHERE id = $1`, id,
	).Scan(&va.ID, &va.VehicleID, &va.DailyPlanID, &va.VisitID, &va.DepartureTime, &va.ReturnTime, &va.DepartureMileage, &va.ReturnMileage, &va.Notes, &va.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "vehicle_assignment", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get vehicle_assignment by id: %w", err)
	}
	return &va, nil
}

func (r *VehicleAssignmentRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.VehicleAssignment], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM operations.vehicle_assignments`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count vehicle assignments: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, vehicle_id, daily_plan_id, visit_id, departure_time, return_time, departure_mileage, return_mileage, notes, created_at
		 FROM operations.vehicle_assignments
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list vehicle assignments: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VehicleAssignment, 0)
	for rows.Next() {
		var va operations.VehicleAssignment
		if err := rows.Scan(&va.ID, &va.VehicleID, &va.DailyPlanID, &va.VisitID, &va.DepartureTime, &va.ReturnTime, &va.DepartureMileage, &va.ReturnMileage, &va.Notes, &va.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan vehicle assignment: %w", err)
		}
		items = append(items, va)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[operations.VehicleAssignment]{Items: items, Total: total}, nil
}

func (r *VehicleAssignmentRepo) Create(ctx context.Context, va *operations.VehicleAssignment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.vehicle_assignments (id, vehicle_id, daily_plan_id, visit_id, departure_time, return_time, departure_mileage, return_mileage, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		va.ID, va.VehicleID, va.DailyPlanID, va.VisitID, va.DepartureTime, va.ReturnTime, va.DepartureMileage, va.ReturnMileage, va.Notes, va.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create vehicle assignment: %w", err)
	}
	return nil
}

func (r *VehicleAssignmentRepo) CreateInTx(ctx context.Context, tx pgx.Tx, va *operations.VehicleAssignment) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO operations.vehicle_assignments (id, vehicle_id, daily_plan_id, visit_id, departure_time, return_time, departure_mileage, return_mileage, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		va.ID, va.VehicleID, va.DailyPlanID, va.VisitID, va.DepartureTime, va.ReturnTime, va.DepartureMileage, va.ReturnMileage, va.Notes, va.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create vehicle assignment: %w", err)
	}
	return nil
}

func (r *VehicleAssignmentRepo) Update(ctx context.Context, id uuid.UUID, va *operations.VehicleAssignment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.vehicle_assignments
		 SET vehicle_id = $2, daily_plan_id = $3, visit_id = $4, departure_time = $5, return_time = $6, departure_mileage = $7, return_mileage = $8, notes = $9
		 WHERE id = $1`,
		id, va.VehicleID, va.DailyPlanID, va.VisitID, va.DepartureTime, va.ReturnTime, va.DepartureMileage, va.ReturnMileage, va.Notes,
	)
	if err != nil {
		return fmt.Errorf("update vehicle assignment: %w", err)
	}
	return nil
}

func (r *VehicleAssignmentRepo) UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, va *operations.VehicleAssignment) error {
	_, err := tx.Exec(ctx,
		`UPDATE operations.vehicle_assignments
		 SET vehicle_id = $2, daily_plan_id = $3, visit_id = $4, departure_time = $5, return_time = $6, departure_mileage = $7, return_mileage = $8, notes = $9
		 WHERE id = $1`,
		id, va.VehicleID, va.DailyPlanID, va.VisitID, va.DepartureTime, va.ReturnTime, va.DepartureMileage, va.ReturnMileage, va.Notes,
	)
	if err != nil {
		return fmt.Errorf("update vehicle assignment: %w", err)
	}
	return nil
}

func (r *VehicleAssignmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.vehicle_assignments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete vehicle assignment: %w", err)
	}
	return nil
}
