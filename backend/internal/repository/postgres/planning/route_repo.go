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

type RouteRepo struct {
	pool *pgxpool.Pool
}

func NewRouteRepo(pool *pgxpool.Pool) *RouteRepo {
	return &RouteRepo{pool: pool}
}

func (r *RouteRepo) GetByID(ctx context.Context, id uuid.UUID) (*planning.Route, error) {
	var rt planning.Route
	err := r.pool.QueryRow(ctx,
		`SELECT id, type, date, daily_plan_id, notes, status, created_at
		 FROM planning.routes
		 WHERE id = $1`, id,
	).Scan(&rt.ID, &rt.Type, &rt.Date, &rt.DailyPlanID, &rt.Notes, &rt.Status, &rt.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "route", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get route by id: %w", err)
	}
	return &rt, nil
}

func (r *RouteRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[planning.Route], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM planning.routes`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count routes: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, type, date, daily_plan_id, notes, status, created_at
		 FROM planning.routes
		 ORDER BY date DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list routes: %w", err)
	}
	defer rows.Close()

	items := make([]planning.Route, 0)
	for rows.Next() {
		var rt planning.Route
		if err := rows.Scan(&rt.ID, &rt.Type, &rt.Date, &rt.DailyPlanID, &rt.Notes, &rt.Status, &rt.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan route: %w", err)
		}
		items = append(items, rt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[planning.Route]{Items: items, Total: total}, nil
}

func (r *RouteRepo) ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.Route, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, type, date, daily_plan_id, notes, status, created_at
		 FROM planning.routes
		 WHERE daily_plan_id = $1
		 ORDER BY date`, planID,
	)
	if err != nil {
		return nil, fmt.Errorf("list routes by plan: %w", err)
	}
	defer rows.Close()

	items := make([]planning.Route, 0)
	for rows.Next() {
		var rt planning.Route
		if err := rows.Scan(&rt.ID, &rt.Type, &rt.Date, &rt.DailyPlanID, &rt.Notes, &rt.Status, &rt.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan route: %w", err)
		}
		items = append(items, rt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *RouteRepo) Create(ctx context.Context, route *planning.Route) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO planning.routes (id, type, date, daily_plan_id, notes, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		route.ID, route.Type, route.Date, route.DailyPlanID, route.Notes, route.Status, route.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create route: %w", err)
	}
	return nil
}

func (r *RouteRepo) Update(ctx context.Context, id uuid.UUID, route *planning.Route) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE planning.routes
		 SET type = $2, date = $3, daily_plan_id = $4, notes = $5, status = $6
		 WHERE id = $1`,
		id, route.Type, route.Date, route.DailyPlanID, route.Notes, route.Status,
	)
	if err != nil {
		return fmt.Errorf("update route: %w", err)
	}
	return nil
}

func (r *RouteRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM planning.routes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete route: %w", err)
	}
	return nil
}
