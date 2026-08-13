package postgresinventory

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type ToolRepo struct {
	pool *pgxpool.Pool
}

func NewToolRepo(pool *pgxpool.Pool) *ToolRepo {
	return &ToolRepo{pool: pool}
}

func (r *ToolRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.Tool, error) {
	var t inventory.Tool
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, status, created_at
		 FROM inventory.tools
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.Code, &t.Name, &t.Status, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "tool", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get tool by id: %w", err)
	}
	return &t, nil
}

func (r *ToolRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.Tool], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.tools`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count tools: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, code, name, status, created_at
		 FROM inventory.tools
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list tools: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.Tool, 0)
	for rows.Next() {
		var t inventory.Tool
		if err := rows.Scan(&t.ID, &t.Code, &t.Name, &t.Status, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan tool: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.Tool]{Items: items, Total: total}, nil
}

func (r *ToolRepo) Create(ctx context.Context, tool *inventory.Tool) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.tools (id, code, name, status, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		tool.ID, tool.Code, tool.Name, tool.Status, tool.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create tool: %w", err)
	}
	return nil
}

func (r *ToolRepo) Update(ctx context.Context, id uuid.UUID, tool *inventory.Tool) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.tools
		 SET code = $2, name = $3, status = $4
		 WHERE id = $1`,
		id, tool.Code, tool.Name, tool.Status,
	)
	if err != nil {
		return fmt.Errorf("update tool: %w", err)
	}
	return nil
}

func (r *ToolRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.tools WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete tool: %w", err)
	}
	return nil
}
