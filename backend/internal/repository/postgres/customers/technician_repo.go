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

type TechnicianRepo struct {
	pool *pgxpool.Pool
}

func NewTechnicianRepo(pool *pgxpool.Pool) *TechnicianRepo {
	return &TechnicianRepo{pool: pool}
}

func (r *TechnicianRepo) GetByID(ctx context.Context, id uuid.UUID) (*customers.Technician, error) {
	var t customers.Technician
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, is_active, specialties, created_at
		 FROM customers.technicians
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.UserID, &t.Name, &t.IsActive, &t.Specialties, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "technician", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get technician by id: %w", err)
	}
	return &t, nil
}

func (r *TechnicianRepo) GetByUserID(ctx context.Context, userID string) (*customers.Technician, error) {
	var t customers.Technician
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, is_active, specialties, created_at
		 FROM customers.technicians
		 WHERE user_id = $1`, userID,
	).Scan(&t.ID, &t.UserID, &t.Name, &t.IsActive, &t.Specialties, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "technician", ID: userID}
	}
	if err != nil {
		return nil, fmt.Errorf("get technician by user id: %w", err)
	}
	return &t, nil
}

func (r *TechnicianRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[customers.Technician], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM customers.technicians`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count technicians: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, is_active, specialties, created_at
		 FROM customers.technicians
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list technicians: %w", err)
	}
	defer rows.Close()

	items := make([]customers.Technician, 0)
	for rows.Next() {
		var t customers.Technician
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.IsActive, &t.Specialties, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan technician: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[customers.Technician]{Items: items, Total: total}, nil
}

func (r *TechnicianRepo) Create(ctx context.Context, tech *customers.Technician) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO customers.technicians (id, user_id, name, is_active, specialties, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		tech.ID, tech.UserID, tech.Name, tech.IsActive, tech.Specialties, tech.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create technician: %w", err)
	}
	return nil
}

func (r *TechnicianRepo) Update(ctx context.Context, id uuid.UUID, tech *customers.Technician) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE customers.technicians
		 SET user_id = $2, name = $3, is_active = $4, specialties = $5
		 WHERE id = $1`,
		id, tech.UserID, tech.Name, tech.IsActive, tech.Specialties,
	)
	if err != nil {
		return fmt.Errorf("update technician: %w", err)
	}
	return nil
}

func (r *TechnicianRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM customers.technicians WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete technician: %w", err)
	}
	return nil
}
