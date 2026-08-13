package postgresplanning

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/service"
)

type DailyLoadToolRepo struct {
	pool *pgxpool.Pool
}

func NewDailyLoadToolRepo(pool *pgxpool.Pool) *DailyLoadToolRepo {
	return &DailyLoadToolRepo{pool: pool}
}

func (r *DailyLoadToolRepo) GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadTool, error) {
	var t planning.DailyLoadTool
	err := r.pool.QueryRow(ctx,
		`SELECT id, daily_plan_id, tool_id, loaded_quantity, returned_quantity, created_at
		 FROM planning.daily_load_tools
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.DailyPlanID, &t.ToolID, &t.LoadedQuantity, &t.ReturnedQuantity, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "daily_load_tool", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get daily_load_tool by id: %w", err)
	}
	return &t, nil
}

func (r *DailyLoadToolRepo) ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadTool, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, daily_plan_id, tool_id, loaded_quantity, returned_quantity, created_at
		 FROM planning.daily_load_tools
		 WHERE daily_plan_id = $1
		 ORDER BY created_at`, planID,
	)
	if err != nil {
		return nil, fmt.Errorf("list daily load tools by plan: %w", err)
	}
	defer rows.Close()

	items := make([]planning.DailyLoadTool, 0)
	for rows.Next() {
		var t planning.DailyLoadTool
		if err := rows.Scan(&t.ID, &t.DailyPlanID, &t.ToolID, &t.LoadedQuantity, &t.ReturnedQuantity, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan daily load tool: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *DailyLoadToolRepo) Create(ctx context.Context, t *planning.DailyLoadTool) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO planning.daily_load_tools (id, daily_plan_id, tool_id, loaded_quantity, returned_quantity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		t.ID, t.DailyPlanID, t.ToolID, t.LoadedQuantity, t.ReturnedQuantity, t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create daily load tool: %w", err)
	}
	return nil
}

func (r *DailyLoadToolRepo) Update(ctx context.Context, id uuid.UUID, t *planning.DailyLoadTool) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE planning.daily_load_tools
		 SET daily_plan_id = $2, tool_id = $3, loaded_quantity = $4, returned_quantity = $5
		 WHERE id = $1`,
		id, t.DailyPlanID, t.ToolID, t.LoadedQuantity, t.ReturnedQuantity,
	)
	if err != nil {
		return fmt.Errorf("update daily load tool: %w", err)
	}
	return nil
}

func (r *DailyLoadToolRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM planning.daily_load_tools WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete daily load tool: %w", err)
	}
	return nil
}
