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

type VisitAssignmentRepo struct {
	pool *pgxpool.Pool
}

func NewVisitAssignmentRepo(pool *pgxpool.Pool) *VisitAssignmentRepo {
	return &VisitAssignmentRepo{pool: pool}
}

func (r *VisitAssignmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitAssignment, error) {
	var a operations.VisitAssignment
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, tech_id, role_id, created_at
		 FROM operations.visit_assignments
		 WHERE id = $1`, id,
	).Scan(&a.ID, &a.VisitID, &a.TechID, &a.RoleID, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_assignment", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_assignment by id: %w", err)
	}
	return &a, nil
}

func (r *VisitAssignmentRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitAssignment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT va.id, va.visit_id, va.tech_id, va.role_id,
		        COALESCE(t.name, '') AS tech_name,
		        COALESCE(tr.name, '') AS role_name,
		        va.created_at
		 FROM operations.visit_assignments va
		 LEFT JOIN customers.technicians t ON t.id = va.tech_id
		 LEFT JOIN core.tech_roles tr ON tr.id = va.role_id
		 WHERE va.visit_id = $1
		 ORDER BY va.created_at`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit assignments by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitAssignment, 0)
	for rows.Next() {
		var a operations.VisitAssignment
		if err := rows.Scan(&a.ID, &a.VisitID, &a.TechID, &a.RoleID, &a.TechName, &a.RoleName, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit assignment: %w", err)
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitAssignmentRepo) Create(ctx context.Context, a *operations.VisitAssignment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_assignments (id, visit_id, tech_id, role_id, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		a.ID, a.VisitID, a.TechID, a.RoleID, a.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit assignment: %w", err)
	}
	return nil
}

func (r *VisitAssignmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_assignments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit assignment: %w", err)
	}
	return nil
}
