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

type DailyLoadEPPRepo struct {
	pool *pgxpool.Pool
}

func NewDailyLoadEPPRepo(pool *pgxpool.Pool) *DailyLoadEPPRepo {
	return &DailyLoadEPPRepo{pool: pool}
}

func (r *DailyLoadEPPRepo) GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadEPP, error) {
	var e planning.DailyLoadEPP
	err := r.pool.QueryRow(ctx,
		`SELECT id, daily_plan_id, epp_id, loaded_quantity, returned_quantity, created_at
		 FROM planning.daily_load_epps
		 WHERE id = $1`, id,
	).Scan(&e.ID, &e.DailyPlanID, &e.EPPID, &e.LoadedQuantity, &e.ReturnedQuantity, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "daily_load_epp", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get daily_load_epp by id: %w", err)
	}
	return &e, nil
}

func (r *DailyLoadEPPRepo) ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadEPP, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, daily_plan_id, epp_id, loaded_quantity, returned_quantity, created_at
		 FROM planning.daily_load_epps
		 WHERE daily_plan_id = $1
		 ORDER BY created_at`, planID,
	)
	if err != nil {
		return nil, fmt.Errorf("list daily load epps by plan: %w", err)
	}
	defer rows.Close()

	items := make([]planning.DailyLoadEPP, 0)
	for rows.Next() {
		var e planning.DailyLoadEPP
		if err := rows.Scan(&e.ID, &e.DailyPlanID, &e.EPPID, &e.LoadedQuantity, &e.ReturnedQuantity, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan daily load epp: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *DailyLoadEPPRepo) Create(ctx context.Context, e *planning.DailyLoadEPP) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO planning.daily_load_epps (id, daily_plan_id, epp_id, loaded_quantity, returned_quantity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ID, e.DailyPlanID, e.EPPID, e.LoadedQuantity, e.ReturnedQuantity, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create daily load epp: %w", err)
	}
	return nil
}

func (r *DailyLoadEPPRepo) Update(ctx context.Context, id uuid.UUID, e *planning.DailyLoadEPP) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE planning.daily_load_epps
		 SET daily_plan_id = $2, epp_id = $3, loaded_quantity = $4, returned_quantity = $5
		 WHERE id = $1`,
		id, e.DailyPlanID, e.EPPID, e.LoadedQuantity, e.ReturnedQuantity,
	)
	if err != nil {
		return fmt.Errorf("update daily load epp: %w", err)
	}
	return nil
}

func (r *DailyLoadEPPRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM planning.daily_load_epps WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete daily load epp: %w", err)
	}
	return nil
}
