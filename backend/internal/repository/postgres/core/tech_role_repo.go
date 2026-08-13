package postgrescore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/core"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type CoreTechRoleRepo struct {
	pool *pgxpool.Pool
}

func NewCoreTechRoleRepo(pool *pgxpool.Pool) *CoreTechRoleRepo {
	return &CoreTechRoleRepo{pool: pool}
}

func (r *CoreTechRoleRepo) GetByID(ctx context.Context, id uuid.UUID) (*core.TechRole, error) {
	var tr core.TechRole
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, active, created_at
		 FROM core.tech_roles
		 WHERE id = $1`, id,
	).Scan(&tr.ID, &tr.Code, &tr.Name, &tr.Active, &tr.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "tech_role", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get tech_role by id: %w", err)
	}
	return &tr, nil
}

func (r *CoreTechRoleRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[core.TechRole], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM core.tech_roles`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count tech roles: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, code, name, active, created_at
		 FROM core.tech_roles
		 ORDER BY code
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list tech roles: %w", err)
	}
	defer rows.Close()

	items := make([]core.TechRole, 0)
	for rows.Next() {
		var tr core.TechRole
		if err := rows.Scan(&tr.ID, &tr.Code, &tr.Name, &tr.Active, &tr.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan tech role: %w", err)
		}
		items = append(items, tr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[core.TechRole]{Items: items, Total: total}, nil
}

func (r *CoreTechRoleRepo) Create(ctx context.Context, role *core.TechRole) error {
	role.ID = uuid.New()
	role.CreatedAt = time.Now()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO core.tech_roles (id, code, name, active, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		role.ID, role.Code, role.Name, role.Active, role.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create tech role: %w", err)
	}
	return nil
}

func (r *CoreTechRoleRepo) Update(ctx context.Context, id uuid.UUID, role *core.TechRole) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE core.tech_roles
		 SET code = $2, name = $3, active = $4
		 WHERE id = $1`,
		id, role.Code, role.Name, role.Active,
	)
	if err != nil {
		return fmt.Errorf("update tech role: %w", err)
	}
	return nil
}

func (r *CoreTechRoleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM core.tech_roles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete tech role: %w", err)
	}
	return nil
}
