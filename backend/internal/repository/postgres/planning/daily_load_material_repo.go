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

type DailyLoadMaterialRepo struct {
	pool *pgxpool.Pool
}

func NewDailyLoadMaterialRepo(pool *pgxpool.Pool) *DailyLoadMaterialRepo {
	return &DailyLoadMaterialRepo{pool: pool}
}

func (r *DailyLoadMaterialRepo) GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadMaterial, error) {
	var m planning.DailyLoadMaterial
	err := r.pool.QueryRow(ctx,
		`SELECT id, daily_plan_id, material_id, loaded_quantity, returned_quantity, created_at
		 FROM planning.daily_load_materials
		 WHERE id = $1`, id,
	).Scan(&m.ID, &m.DailyPlanID, &m.MaterialID, &m.LoadedQuantity, &m.ReturnedQuantity, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "daily_load_material", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get daily_load_material by id: %w", err)
	}
	return &m, nil
}

func (r *DailyLoadMaterialRepo) ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadMaterial, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, daily_plan_id, material_id, loaded_quantity, returned_quantity, created_at
		 FROM planning.daily_load_materials
		 WHERE daily_plan_id = $1
		 ORDER BY created_at`, planID,
	)
	if err != nil {
		return nil, fmt.Errorf("list daily load materials by plan: %w", err)
	}
	defer rows.Close()

	items := make([]planning.DailyLoadMaterial, 0)
	for rows.Next() {
		var m planning.DailyLoadMaterial
		if err := rows.Scan(&m.ID, &m.DailyPlanID, &m.MaterialID, &m.LoadedQuantity, &m.ReturnedQuantity, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan daily load material: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *DailyLoadMaterialRepo) Create(ctx context.Context, m *planning.DailyLoadMaterial) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO planning.daily_load_materials (id, daily_plan_id, material_id, loaded_quantity, returned_quantity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		m.ID, m.DailyPlanID, m.MaterialID, m.LoadedQuantity, m.ReturnedQuantity, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create daily load material: %w", err)
	}
	return nil
}

func (r *DailyLoadMaterialRepo) Update(ctx context.Context, id uuid.UUID, m *planning.DailyLoadMaterial) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE planning.daily_load_materials
		 SET daily_plan_id = $2, material_id = $3, loaded_quantity = $4, returned_quantity = $5
		 WHERE id = $1`,
		id, m.DailyPlanID, m.MaterialID, m.LoadedQuantity, m.ReturnedQuantity,
	)
	if err != nil {
		return fmt.Errorf("update daily load material: %w", err)
	}
	return nil
}

func (r *DailyLoadMaterialRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM planning.daily_load_materials WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete daily load material: %w", err)
	}
	return nil
}
