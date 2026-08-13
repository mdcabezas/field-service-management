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

type PartnerAgreementFormRepo struct {
	pool *pgxpool.Pool
}

func NewPartnerAgreementFormRepo(pool *pgxpool.Pool) *PartnerAgreementFormRepo {
	return &PartnerAgreementFormRepo{pool: pool}
}

func (r *PartnerAgreementFormRepo) GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerAgreementForm, error) {
	var f partners.PartnerAgreementForm
	err := r.pool.QueryRow(ctx,
		`SELECT id, agreement_id, work_type, form_template_id, quantity, created_at
		 FROM partners.partner_agreement_forms
		 WHERE id = $1`, id,
	).Scan(&f.ID, &f.AgreementID, &f.WorkType, &f.FormTemplateID, &f.Quantity, &f.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "partner_agreement_form", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get partner_agreement_form by id: %w", err)
	}
	return &f, nil
}

func (r *PartnerAgreementFormRepo) ListByAgreement(ctx context.Context, agreementID uuid.UUID) ([]partners.PartnerAgreementForm, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, agreement_id, work_type, form_template_id, quantity, created_at
		 FROM partners.partner_agreement_forms
		 WHERE agreement_id = $1
		 ORDER BY created_at DESC`, agreementID,
	)
	if err != nil {
		return nil, fmt.Errorf("list partner agreement forms: %w", err)
	}
	defer rows.Close()

	items := make([]partners.PartnerAgreementForm, 0)
	for rows.Next() {
		var f partners.PartnerAgreementForm
		if err := rows.Scan(&f.ID, &f.AgreementID, &f.WorkType, &f.FormTemplateID, &f.Quantity, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan partner agreement form: %w", err)
		}
		items = append(items, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *PartnerAgreementFormRepo) Create(ctx context.Context, form *partners.PartnerAgreementForm) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO partners.partner_agreement_forms (id, agreement_id, work_type, form_template_id, quantity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		form.ID, form.AgreementID, form.WorkType, form.FormTemplateID, form.Quantity, form.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create partner agreement form: %w", err)
	}
	return nil
}

func (r *PartnerAgreementFormRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM partners.partner_agreement_forms WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete partner agreement form: %w", err)
	}
	return nil
}
