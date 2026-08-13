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

type PartnerAgreementRepo struct {
	pool *pgxpool.Pool
}

func NewPartnerAgreementRepo(pool *pgxpool.Pool) *PartnerAgreementRepo {
	return &PartnerAgreementRepo{pool: pool}
}

func (r *PartnerAgreementRepo) GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerAgreement, error) {
	var a partners.PartnerAgreement
	err := r.pool.QueryRow(ctx,
		`SELECT id, partner_id, service_type, rate, active, created_at
		 FROM partners.partner_agreements
		 WHERE id = $1`, id,
	).Scan(&a.ID, &a.PartnerID, &a.ServiceType, &a.Rate, &a.Active, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "partner_agreement", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get partner_agreement by id: %w", err)
	}
	return &a, nil
}

func (r *PartnerAgreementRepo) ListByPartner(ctx context.Context, partnerID uuid.UUID) ([]partners.PartnerAgreement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, partner_id, service_type, rate, active, created_at
		 FROM partners.partner_agreements
		 WHERE partner_id = $1
		 ORDER BY created_at DESC`, partnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list partner agreements: %w", err)
	}
	defer rows.Close()

	items := make([]partners.PartnerAgreement, 0)
	for rows.Next() {
		var a partners.PartnerAgreement
		if err := rows.Scan(&a.ID, &a.PartnerID, &a.ServiceType, &a.Rate, &a.Active, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan partner agreement: %w", err)
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *PartnerAgreementRepo) Create(ctx context.Context, agreement *partners.PartnerAgreement) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO partners.partner_agreements (id, partner_id, service_type, rate, active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		agreement.ID, agreement.PartnerID, agreement.ServiceType, agreement.Rate, agreement.Active, agreement.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create partner agreement: %w", err)
	}
	return nil
}

func (r *PartnerAgreementRepo) Update(ctx context.Context, id uuid.UUID, agreement *partners.PartnerAgreement) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE partners.partner_agreements
		 SET service_type = $2, rate = $3, active = $4
		 WHERE id = $1`,
		id, agreement.ServiceType, agreement.Rate, agreement.Active,
	)
	if err != nil {
		return fmt.Errorf("update partner agreement: %w", err)
	}
	return nil
}

func (r *PartnerAgreementRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM partners.partner_agreements WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete partner agreement: %w", err)
	}
	return nil
}
