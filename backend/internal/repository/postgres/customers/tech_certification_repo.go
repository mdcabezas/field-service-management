package postgrescustomers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/service"
)

type TechCertificationRepo struct {
	pool *pgxpool.Pool
}

func NewTechCertificationRepo(pool *pgxpool.Pool) *TechCertificationRepo {
	return &TechCertificationRepo{pool: pool}
}

func (r *TechCertificationRepo) GetByID(ctx context.Context, id uuid.UUID) (*customers.TechCertification, error) {
	var tc customers.TechCertification
	err := r.pool.QueryRow(ctx,
		`SELECT id, tech_id, cert_id, number, issuer, issue_date, expiry_date, created_at
		 FROM customers.tech_certifications
		 WHERE id = $1`, id,
	).Scan(&tc.ID, &tc.TechID, &tc.CertID, &tc.Number, &tc.Issuer, &tc.IssueDate, &tc.ExpiryDate, &tc.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "tech_certification", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get tech_certification by id: %w", err)
	}
	return &tc, nil
}

func (r *TechCertificationRepo) ListByTech(ctx context.Context, techID uuid.UUID) ([]customers.TechCertification, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tech_id, cert_id, number, issuer, issue_date, expiry_date, created_at
		 FROM customers.tech_certifications
		 WHERE tech_id = $1
		 ORDER BY created_at`, techID,
	)
	if err != nil {
		return nil, fmt.Errorf("list tech certifications: %w", err)
	}
	defer rows.Close()

	items := make([]customers.TechCertification, 0)
	for rows.Next() {
		var tc customers.TechCertification
		if err := rows.Scan(&tc.ID, &tc.TechID, &tc.CertID, &tc.Number, &tc.Issuer, &tc.IssueDate, &tc.ExpiryDate, &tc.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan tech certification: %w", err)
		}
		items = append(items, tc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *TechCertificationRepo) Create(ctx context.Context, cert *customers.TechCertification) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO customers.tech_certifications (id, tech_id, cert_id, number, issuer, issue_date, expiry_date, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		cert.ID, cert.TechID, cert.CertID, cert.Number, cert.Issuer, cert.IssueDate, cert.ExpiryDate, cert.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create tech certification: %w", err)
	}
	return nil
}

func (r *TechCertificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM customers.tech_certifications WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete tech certification: %w", err)
	}
	return nil
}
