package postgrescore

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/core"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type CoreUserRepo struct {
	pool *pgxpool.Pool
}

func NewCoreUserRepo(pool *pgxpool.Pool) *CoreUserRepo {
	return &CoreUserRepo{pool: pool}
}

func (r *CoreUserRepo) GetByID(ctx context.Context, id string) (*core.User, error) {
	var u core.User
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	err = r.pool.QueryRow(ctx,
		`SELECT id, email, role, name, created_at
		 FROM core.users
		 WHERE id = $1`, parsedID,
	).Scan(&u.ID, &u.Email, &u.Role, &u.Name, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "user", ID: id}
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *CoreUserRepo) GetByEmail(ctx context.Context, email string) (*core.UserWithPassword, error) {
	var u core.UserWithPassword
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, role, name, password_hash, created_at
		 FROM core.users
		 WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Role, &u.Name, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "user", ID: email}
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *CoreUserRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[core.User], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM core.users`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, email, role, name, created_at
		 FROM core.users
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	items := make([]core.User, 0)
	for rows.Next() {
		var u core.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.Name, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[core.User]{Items: items, Total: total}, nil
}

func (r *CoreUserRepo) Create(ctx context.Context, user *core.User) error {
	user.ID = uuid.New()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO core.users (id, email, role, name, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Email, user.Role, user.Name, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *CoreUserRepo) Update(ctx context.Context, id string, user *core.User) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}
	_, err = r.pool.Exec(ctx,
		`UPDATE core.users
		 SET email = $2, name = $3, role = $4
		 WHERE id = $1`,
		parsedID, user.Email, user.Name, user.Role,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *CoreUserRepo) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}
	_, err = r.pool.Exec(ctx, `DELETE FROM core.users WHERE id = $1`, parsedID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
