package postgresplanning

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type DailyPlanAssignmentRepo struct {
	pool *pgxpool.Pool
}

func NewDailyPlanAssignmentRepo(pool *pgxpool.Pool) *DailyPlanAssignmentRepo {
	return &DailyPlanAssignmentRepo{pool: pool}
}

func (r *DailyPlanAssignmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyPlanAssignment, error) {
	var a planning.DailyPlanAssignment
	err := r.pool.QueryRow(ctx,
		`SELECT id, daily_plan_id, tech_id, role_id, created_at
		 FROM planning.daily_plan_assignments
		 WHERE id = $1`, id,
	).Scan(&a.ID, &a.DailyPlanID, &a.TechID, &a.RoleID, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "daily_plan_assignment", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get daily_plan_assignment by id: %w", err)
	}
	return &a, nil
}

func (r *DailyPlanAssignmentRepo) ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyPlanAssignment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, daily_plan_id, tech_id, role_id, created_at
		 FROM planning.daily_plan_assignments
		 WHERE daily_plan_id = $1
		 ORDER BY created_at`, planID,
	)
	if err != nil {
		return nil, fmt.Errorf("list daily plan assignments by plan: %w", err)
	}
	defer rows.Close()

	items := make([]planning.DailyPlanAssignment, 0)
	for rows.Next() {
		var a planning.DailyPlanAssignment
		if err := rows.Scan(&a.ID, &a.DailyPlanID, &a.TechID, &a.RoleID, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan daily plan assignment: %w", err)
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *DailyPlanAssignmentRepo) Create(ctx context.Context, a *planning.DailyPlanAssignment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO planning.daily_plan_assignments (id, daily_plan_id, tech_id, role_id, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		a.ID, a.DailyPlanID, a.TechID, a.RoleID, a.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create daily plan assignment: %w", err)
	}
	return nil
}

func (r *DailyPlanAssignmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM planning.daily_plan_assignments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete daily plan assignment: %w", err)
	}
	return nil
}

func (r *DailyPlanAssignmentRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[planning.DailyPlanAssignment], error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, daily_plan_id, tech_id, role_id, created_at
		 FROM planning.daily_plan_assignments
		 ORDER BY created_at
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list daily plan assignments: %w", err)
	}
	defer rows.Close()

	items := make([]planning.DailyPlanAssignment, 0)
	for rows.Next() {
		var a planning.DailyPlanAssignment
		if err := rows.Scan(&a.ID, &a.DailyPlanID, &a.TechID, &a.RoleID, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan daily plan assignment: %w", err)
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	var total int
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM planning.daily_plan_assignments`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count daily plan assignments: %w", err)
	}

	return &repository.ListResult[planning.DailyPlanAssignment]{Items: items, Total: total}, nil
}

func (r *DailyPlanAssignmentRepo) Update(ctx context.Context, id uuid.UUID, a *planning.DailyPlanAssignment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE planning.daily_plan_assignments SET daily_plan_id = $1, tech_id = $2, role_id = $3 WHERE id = $4`,
		a.DailyPlanID, a.TechID, a.RoleID, id,
	)
	if err != nil {
		return fmt.Errorf("update daily plan assignment: %w", err)
	}
	return nil
}
