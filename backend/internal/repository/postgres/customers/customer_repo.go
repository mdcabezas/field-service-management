package postgrescustomers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type CustomerRepo struct {
	pool *pgxpool.Pool
}

func NewCustomerRepo(pool *pgxpool.Pool) *CustomerRepo {
	return &CustomerRepo{pool: pool}
}

func (r *CustomerRepo) GetByID(ctx context.Context, id uuid.UUID) (*customers.Customer, error) {
	var c customers.Customer
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, phone, email, tax_id, created_at, updated_at
		 FROM customers.customers
		 WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.TaxID, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "customer", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get customer by id: %w", err)
	}
	return &c, nil
}

func (r *CustomerRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[customers.Customer], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM customers.customers`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count customers: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, phone, email, tax_id, created_at, updated_at
		 FROM customers.customers
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	items := make([]customers.Customer, 0)
	for rows.Next() {
		var c customers.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.TaxID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[customers.Customer]{Items: items, Total: total}, nil
}

func (r *CustomerRepo) Create(ctx context.Context, customer *customers.Customer) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO customers.customers (id, name, phone, email, tax_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		customer.ID, customer.Name, customer.Phone, customer.Email, customer.TaxID, customer.CreatedAt, customer.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create customer: %w", err)
	}
	return nil
}

func (r *CustomerRepo) Update(ctx context.Context, id uuid.UUID, customer *customers.Customer) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE customers.customers
		 SET name = $2, phone = $3, email = $4, tax_id = $5, updated_at = $6
		 WHERE id = $1`,
		id, customer.Name, customer.Phone, customer.Email, customer.TaxID, customer.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update customer: %w", err)
	}
	return nil
}

func (r *CustomerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM customers.customers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}
	return nil
}
