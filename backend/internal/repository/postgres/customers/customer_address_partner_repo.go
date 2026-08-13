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

type CustomerAddressPartnerRepo struct {
	pool *pgxpool.Pool
}

func NewCustomerAddressPartnerRepo(pool *pgxpool.Pool) *CustomerAddressPartnerRepo {
	return &CustomerAddressPartnerRepo{pool: pool}
}

func (r *CustomerAddressPartnerRepo) GetByID(ctx context.Context, id uuid.UUID) (*customers.CustomerAddressPartner, error) {
	var cap customers.CustomerAddressPartner
	err := r.pool.QueryRow(ctx,
		`SELECT id, customer_address_id, partner_id, from_date, to_date, notes, created_at
		 FROM customers.customer_address_partners
		 WHERE id = $1`, id,
	).Scan(&cap.ID, &cap.CustomerAddressID, &cap.PartnerID, &cap.FromDate, &cap.ToDate, &cap.Notes, &cap.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "customer_address_partner", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get customer_address_partner by id: %w", err)
	}
	return &cap, nil
}

func (r *CustomerAddressPartnerRepo) ListByCustomerAddress(ctx context.Context, caID uuid.UUID) ([]customers.CustomerAddressPartner, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, customer_address_id, partner_id, from_date, to_date, notes, created_at
		 FROM customers.customer_address_partners
		 WHERE customer_address_id = $1
		 ORDER BY from_date`, caID,
	)
	if err != nil {
		return nil, fmt.Errorf("list customer address partners: %w", err)
	}
	defer rows.Close()

	items := make([]customers.CustomerAddressPartner, 0)
	for rows.Next() {
		var cap customers.CustomerAddressPartner
		if err := rows.Scan(&cap.ID, &cap.CustomerAddressID, &cap.PartnerID, &cap.FromDate, &cap.ToDate, &cap.Notes, &cap.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan customer address partner: %w", err)
		}
		items = append(items, cap)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *CustomerAddressPartnerRepo) Create(ctx context.Context, cap *customers.CustomerAddressPartner) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO customers.customer_address_partners (id, customer_address_id, partner_id, from_date, to_date, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		cap.ID, cap.CustomerAddressID, cap.PartnerID, cap.FromDate, cap.ToDate, cap.Notes, cap.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create customer address partner: %w", err)
	}
	return nil
}

func (r *CustomerAddressPartnerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM customers.customer_address_partners WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete customer address partner: %w", err)
	}
	return nil
}
