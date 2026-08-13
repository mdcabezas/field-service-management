package postgrespartners

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type PartnerRepo struct {
	pool *pgxpool.Pool
}

func NewPartnerRepo(pool *pgxpool.Pool) *PartnerRepo {
	return &PartnerRepo{pool: pool}
}

func (r *PartnerRepo) GetByID(ctx context.Context, id uuid.UUID) (*partners.Partner, error) {
	var p partners.Partner
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, tax_id, status, created_at, updated_at
		 FROM partners.partners
		 WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.TaxID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "partner", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get partner by id: %w", err)
	}
	return &p, nil
}

func (r *PartnerRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[partners.Partner], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM partners.partners`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count partners: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, tax_id, status, created_at, updated_at
		 FROM partners.partners
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list partners: %w", err)
	}
	defer rows.Close()

	items := make([]partners.Partner, 0)
	for rows.Next() {
		var p partners.Partner
		if err := rows.Scan(&p.ID, &p.Name, &p.TaxID, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan partner: %w", err)
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[partners.Partner]{Items: items, Total: total}, nil
}

func (r *PartnerRepo) Create(ctx context.Context, partner *partners.Partner) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO partners.partners (id, name, tax_id, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		partner.ID, partner.Name, partner.TaxID, partner.Status, partner.CreatedAt, partner.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create partner: %w", err)
	}
	return nil
}

func (r *PartnerRepo) Update(ctx context.Context, id uuid.UUID, partner *partners.Partner) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE partners.partners
		 SET name = $2, tax_id = $3, status = $4, updated_at = $5
		 WHERE id = $1`,
		id, partner.Name, partner.TaxID, partner.Status, partner.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update partner: %w", err)
	}
	return nil
}

func (r *PartnerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM partners.partners WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete partner: %w", err)
	}
	return nil
}
