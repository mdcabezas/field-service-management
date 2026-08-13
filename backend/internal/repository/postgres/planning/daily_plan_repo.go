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

type DailyPlanRepo struct {
	pool *pgxpool.Pool
}

func NewDailyPlanRepo(pool *pgxpool.Pool) *DailyPlanRepo {
	return &DailyPlanRepo{pool: pool}
}

func (r *DailyPlanRepo) GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyPlan, error) {
	var p planning.DailyPlan
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, date, notes, created_at
		 FROM planning.daily_plans
		 WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Date, &p.Notes, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "daily_plan", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get daily_plan by id: %w", err)
	}
	return &p, nil
}

func (r *DailyPlanRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[planning.DailyPlan], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM planning.daily_plans`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count daily plans: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, date, notes, created_at
		 FROM planning.daily_plans
		 ORDER BY date DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list daily plans: %w", err)
	}
	defer rows.Close()

	items := make([]planning.DailyPlan, 0)
	for rows.Next() {
		var p planning.DailyPlan
		if err := rows.Scan(&p.ID, &p.Name, &p.Date, &p.Notes, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan daily plan: %w", err)
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[planning.DailyPlan]{Items: items, Total: total}, nil
}

func (r *DailyPlanRepo) Create(ctx context.Context, plan *planning.DailyPlan) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO planning.daily_plans (id, name, date, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		plan.ID, plan.Name, plan.Date, plan.Notes, plan.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create daily plan: %w", err)
	}
	return nil
}

func (r *DailyPlanRepo) Update(ctx context.Context, id uuid.UUID, plan *planning.DailyPlan) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE planning.daily_plans
		 SET name = $2, date = $3, notes = $4
		 WHERE id = $1`,
		id, plan.Name, plan.Date, plan.Notes,
	)
	if err != nil {
		return fmt.Errorf("update daily plan: %w", err)
	}
	return nil
}

func (r *DailyPlanRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM planning.daily_plans WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete daily plan: %w", err)
	}
	return nil
}
