package postgresoperations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/service"
)

type VisitCheckpointRepo struct {
	pool *pgxpool.Pool
}

func NewVisitCheckpointRepo(pool *pgxpool.Pool) *VisitCheckpointRepo {
	return &VisitCheckpointRepo{pool: pool}
}

func (r *VisitCheckpointRepo) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *VisitCheckpointRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitCheckpoint, error) {
	var cp operations.VisitCheckpoint
	var geomRaw []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, type, ST_AsEWKB(geom), timestamp, device_info, created_at
		 FROM operations.visit_checkpoints
		 WHERE id = $1`, id,
	).Scan(&cp.ID, &cp.VisitID, &cp.Type, &geomRaw, &cp.Timestamp, &cp.DeviceInfo, &cp.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_checkpoint", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_checkpoint by id: %w", err)
	}
	if len(geomRaw) > 0 {
		if geom, _, err := ewkb.Unmarshal(geomRaw); err != nil {
			slog.Warn("failed to unmarshal checkpoint geometry", "id", cp.ID, "error", err)
		} else if p, ok := geom.(orb.Point); ok {
			cp.Geom = p
		}
	}
	return &cp, nil
}

func (r *VisitCheckpointRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitCheckpoint, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, type, ST_AsEWKB(geom), timestamp, device_info, created_at
		 FROM operations.visit_checkpoints
		 WHERE visit_id = $1
		 ORDER BY timestamp
		 LIMIT 100`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit checkpoints by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitCheckpoint, 0)
	for rows.Next() {
		var cp operations.VisitCheckpoint
		var geomRaw []byte
		if err := rows.Scan(&cp.ID, &cp.VisitID, &cp.Type, &geomRaw, &cp.Timestamp, &cp.DeviceInfo, &cp.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit checkpoint: %w", err)
		}
		if len(geomRaw) > 0 {
			if geom, _, err := ewkb.Unmarshal(geomRaw); err != nil {
				slog.Warn("failed to unmarshal checkpoint geometry", "id", cp.ID, "error", err)
			} else if p, ok := geom.(orb.Point); ok {
				cp.Geom = p
			}
		}
		items = append(items, cp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitCheckpointRepo) Create(ctx context.Context, cp *operations.VisitCheckpoint) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_checkpoints (id, visit_id, type, geom, timestamp, device_info, created_at)
		 VALUES ($1, $2, $3, ST_GeomFromEWKB($4), $5, $6, $7)`,
		cp.ID, cp.VisitID, cp.Type, ewkb.Value(cp.Geom, 4326), cp.Timestamp, cp.DeviceInfo, cp.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit checkpoint: %w", err)
	}
	return nil
}

func (r *VisitCheckpointRepo) CreateInTx(ctx context.Context, tx pgx.Tx, cp *operations.VisitCheckpoint) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO operations.visit_checkpoints (id, visit_id, type, geom, timestamp, device_info, created_at)
		 VALUES ($1, $2, $3, ST_GeomFromEWKB($4), $5, $6, $7)`,
		cp.ID, cp.VisitID, cp.Type, ewkb.Value(cp.Geom, 4326), cp.Timestamp, cp.DeviceInfo, cp.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit checkpoint in tx: %w", err)
	}
	return nil
}

func (r *VisitCheckpointRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_checkpoints WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit checkpoint: %w", err)
	}
	return nil
}
