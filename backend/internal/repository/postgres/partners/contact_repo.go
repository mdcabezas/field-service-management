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

type PartnerContactRepo struct {
	pool *pgxpool.Pool
}

func NewPartnerContactRepo(pool *pgxpool.Pool) *PartnerContactRepo {
	return &PartnerContactRepo{pool: pool}
}

func (r *PartnerContactRepo) GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerContact, error) {
	var c partners.PartnerContact
	err := r.pool.QueryRow(ctx,
		`SELECT id, partner_id, name, position, phone, email, created_at
		 FROM partners.partner_contacts
		 WHERE id = $1`, id,
	).Scan(&c.ID, &c.PartnerID, &c.Name, &c.Position, &c.Phone, &c.Email, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "partner_contact", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get partner_contact by id: %w", err)
	}
	return &c, nil
}

func (r *PartnerContactRepo) ListByPartner(ctx context.Context, partnerID uuid.UUID) ([]partners.PartnerContact, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, partner_id, name, position, phone, email, created_at
		 FROM partners.partner_contacts
		 WHERE partner_id = $1
		 ORDER BY created_at DESC`, partnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list partner contacts: %w", err)
	}
	defer rows.Close()

	items := make([]partners.PartnerContact, 0)
	for rows.Next() {
		var c partners.PartnerContact
		if err := rows.Scan(&c.ID, &c.PartnerID, &c.Name, &c.Position, &c.Phone, &c.Email, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan partner contact: %w", err)
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *PartnerContactRepo) Create(ctx context.Context, contact *partners.PartnerContact) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO partners.partner_contacts (id, partner_id, name, position, phone, email, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		contact.ID, contact.PartnerID, contact.Name, contact.Position, contact.Phone, contact.Email, contact.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create partner contact: %w", err)
	}
	return nil
}

func (r *PartnerContactRepo) Update(ctx context.Context, id uuid.UUID, contact *partners.PartnerContact) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE partners.partner_contacts
		 SET name = $2, position = $3, phone = $4, email = $5
		 WHERE id = $1`,
		id, contact.Name, contact.Position, contact.Phone, contact.Email,
	)
	if err != nil {
		return fmt.Errorf("update partner contact: %w", err)
	}
	return nil
}

func (r *PartnerContactRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM partners.partner_contacts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete partner contact: %w", err)
	}
	return nil
}
