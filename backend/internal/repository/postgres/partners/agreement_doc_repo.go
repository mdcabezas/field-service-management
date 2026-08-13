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

type PartnerAgreementDocRepo struct {
	pool *pgxpool.Pool
}

func NewPartnerAgreementDocRepo(pool *pgxpool.Pool) *PartnerAgreementDocRepo {
	return &PartnerAgreementDocRepo{pool: pool}
}

func (r *PartnerAgreementDocRepo) GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerAgreementDoc, error) {
	var d partners.PartnerAgreementDoc
	err := r.pool.QueryRow(ctx,
		`SELECT id, agreement_id, work_type, doc_name, required, created_at
		 FROM partners.partner_agreement_docs
		 WHERE id = $1`, id,
	).Scan(&d.ID, &d.AgreementID, &d.WorkType, &d.DocName, &d.Required, &d.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "partner_agreement_doc", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get partner_agreement_doc by id: %w", err)
	}
	return &d, nil
}

func (r *PartnerAgreementDocRepo) ListByAgreement(ctx context.Context, agreementID uuid.UUID) ([]partners.PartnerAgreementDoc, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, agreement_id, work_type, doc_name, required, created_at
		 FROM partners.partner_agreement_docs
		 WHERE agreement_id = $1
		 ORDER BY created_at DESC`, agreementID,
	)
	if err != nil {
		return nil, fmt.Errorf("list partner agreement docs: %w", err)
	}
	defer rows.Close()

	items := make([]partners.PartnerAgreementDoc, 0)
	for rows.Next() {
		var d partners.PartnerAgreementDoc
		if err := rows.Scan(&d.ID, &d.AgreementID, &d.WorkType, &d.DocName, &d.Required, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan partner agreement doc: %w", err)
		}
		items = append(items, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *PartnerAgreementDocRepo) Create(ctx context.Context, doc *partners.PartnerAgreementDoc) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO partners.partner_agreement_docs (id, agreement_id, work_type, doc_name, required, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		doc.ID, doc.AgreementID, doc.WorkType, doc.DocName, doc.Required, doc.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create partner agreement doc: %w", err)
	}
	return nil
}

func (r *PartnerAgreementDocRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM partners.partner_agreement_docs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete partner agreement doc: %w", err)
	}
	return nil
}
