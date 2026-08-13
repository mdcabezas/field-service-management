package postgrespartners

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/service"
)

type SLARepo struct {
	pool *pgxpool.Pool
}

func NewSLARepo(pool *pgxpool.Pool) *SLARepo {
	return &SLARepo{pool: pool}
}

func (r *SLARepo) GetByID(ctx context.Context, id uuid.UUID) (*partners.SLA, error) {
	var s partners.SLA
	err := r.pool.QueryRow(ctx,
		`SELECT id, partner_id, name, description, work_type, response_hours, resolution_hours,
		        compliance_target, active, valid_from, valid_until, created_at
		 FROM partners.slas
		 WHERE id = $1`, id,
	).Scan(
		&s.ID, &s.PartnerID, &s.Name, &s.Description, &s.WorkType,
		&s.ResponseHours, &s.ResolutionHours, &s.ComplianceTarget,
		&s.Active, &s.ValidFrom, &s.ValidUntil, &s.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "sla", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get sla by id: %w", err)
	}
	return &s, nil
}

func (r *SLARepo) ListByPartner(ctx context.Context, partnerID uuid.UUID) ([]partners.SLA, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, partner_id, name, description, work_type, response_hours, resolution_hours,
		        compliance_target, active, valid_from, valid_until, created_at
		 FROM partners.slas
		 WHERE partner_id = $1
		 ORDER BY created_at DESC`, partnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list slas: %w", err)
	}
	defer rows.Close()

	items := make([]partners.SLA, 0)
	for rows.Next() {
		var s partners.SLA
		if err := rows.Scan(
			&s.ID, &s.PartnerID, &s.Name, &s.Description, &s.WorkType,
			&s.ResponseHours, &s.ResolutionHours, &s.ComplianceTarget,
			&s.Active, &s.ValidFrom, &s.ValidUntil, &s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan sla: %w", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *SLARepo) Create(ctx context.Context, sla *partners.SLA) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO partners.slas (id, partner_id, name, description, work_type, response_hours, resolution_hours,
		                           compliance_target, active, valid_from, valid_until, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		sla.ID, sla.PartnerID, sla.Name, sla.Description, sla.WorkType,
		sla.ResponseHours, sla.ResolutionHours, sla.ComplianceTarget,
		sla.Active, sla.ValidFrom, sla.ValidUntil, sla.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create sla: %w", err)
	}
	return nil
}

func (r *SLARepo) Update(ctx context.Context, id uuid.UUID, sla *partners.SLA) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE partners.slas
		 SET name = $2, description = $3, work_type = $4, response_hours = $5, resolution_hours = $6,
		     compliance_target = $7, active = $8, valid_from = $9, valid_until = $10
		 WHERE id = $1`,
		id, sla.Name, sla.Description, sla.WorkType,
		sla.ResponseHours, sla.ResolutionHours, sla.ComplianceTarget,
		sla.Active, sla.ValidFrom, sla.ValidUntil,
	)
	if err != nil {
		return fmt.Errorf("update sla: %w", err)
	}
	return nil
}

func (r *SLARepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM partners.slas WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete sla: %w", err)
	}
	return nil
}
