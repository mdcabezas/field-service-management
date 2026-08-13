package postgresplanning

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/repository"
)

type RouteTypeRepo struct {
	pool *pgxpool.Pool
}

func NewRouteTypeRepo(pool *pgxpool.Pool) *RouteTypeRepo {
	return &RouteTypeRepo{pool: pool}
}

func (r *RouteTypeRepo) List(ctx context.Context) (*repository.ListResult[planning.RouteType], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM planning.route_types WHERE active = true`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count route types: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, code, name, active
		 FROM planning.route_types
		 WHERE active = true
		 ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list route types: %w", err)
	}
	defer rows.Close()

	items := make([]planning.RouteType, 0)
	for rows.Next() {
		var rt planning.RouteType
		if err := rows.Scan(&rt.ID, &rt.Code, &rt.Name, &rt.Active); err != nil {
			return nil, fmt.Errorf("scan route type: %w", err)
		}
		items = append(items, rt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[planning.RouteType]{Items: items, Total: total}, nil
}
