package postgresoperations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/service"
)

type VisitRentalRepo struct {
	pool *pgxpool.Pool
}

func NewVisitRentalRepo(pool *pgxpool.Pool) *VisitRentalRepo {
	return &VisitRentalRepo{pool: pool}
}

func (r *VisitRentalRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitRental, error) {
	var rv operations.VisitRental
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, rental_id, start_time, end_time, hours, total_cost, reason, created_at
		 FROM operations.visit_rentals
		 WHERE id = $1`, id,
	).Scan(&rv.ID, &rv.VisitID, &rv.RentalID, &rv.StartTime, &rv.EndTime, &rv.Hours, &rv.TotalCost, &rv.Reason, &rv.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_rental", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_rental by id: %w", err)
	}
	return &rv, nil
}

func (r *VisitRentalRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitRental, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, rental_id, start_time, end_time, hours, total_cost, reason, created_at
		 FROM operations.visit_rentals
		 WHERE visit_id = $1
		 ORDER BY start_time`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit rentals by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitRental, 0)
	for rows.Next() {
		var rv operations.VisitRental
		if err := rows.Scan(&rv.ID, &rv.VisitID, &rv.RentalID, &rv.StartTime, &rv.EndTime, &rv.Hours, &rv.TotalCost, &rv.Reason, &rv.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit rental: %w", err)
		}
		items = append(items, rv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitRentalRepo) Create(ctx context.Context, rv *operations.VisitRental) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_rentals (id, visit_id, rental_id, start_time, end_time, hours, total_cost, reason, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		rv.ID, rv.VisitID, rv.RentalID, rv.StartTime, rv.EndTime, rv.Hours, rv.TotalCost, rv.Reason, rv.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit rental: %w", err)
	}
	return nil
}

func (r *VisitRentalRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_rentals WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit rental: %w", err)
	}
	return nil
}
