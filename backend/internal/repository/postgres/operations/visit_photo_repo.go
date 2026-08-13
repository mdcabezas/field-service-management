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

type VisitPhotoRepo struct {
	pool *pgxpool.Pool
}

func NewVisitPhotoRepo(pool *pgxpool.Pool) *VisitPhotoRepo {
	return &VisitPhotoRepo{pool: pool}
}

func (r *VisitPhotoRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitPhoto, error) {
	var ph operations.VisitPhoto
	var geomRaw []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, checkpoint_id, url, ST_AsEWKB(geom), timestamp, stage, finding_type, created_at
		 FROM operations.visit_photos
		 WHERE id = $1`, id,
	).Scan(&ph.ID, &ph.VisitID, &ph.CheckpointID, &ph.URL, &geomRaw,
		&ph.Timestamp, &ph.Stage, &ph.FindingType, &ph.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_photo", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_photo by id: %w", err)
	}
	if len(geomRaw) > 0 {
		if geom, _, err := ewkb.Unmarshal(geomRaw); err != nil {
			slog.Warn("failed to unmarshal photo geometry", "id", ph.ID, "error", err)
		} else if p, ok := geom.(orb.Point); ok {
			ph.Geom = p
		}
	}
	return &ph, nil
}

func (r *VisitPhotoRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitPhoto, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, checkpoint_id, url, ST_AsEWKB(geom), timestamp, stage, finding_type, created_at
		 FROM operations.visit_photos
		 WHERE visit_id = $1
		 ORDER BY timestamp
		 LIMIT 100`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit photos by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitPhoto, 0)
	for rows.Next() {
		var ph operations.VisitPhoto
		var geomRaw []byte
		if err := rows.Scan(&ph.ID, &ph.VisitID, &ph.CheckpointID, &ph.URL, &geomRaw,
			&ph.Timestamp, &ph.Stage, &ph.FindingType, &ph.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit photo: %w", err)
		}
		if len(geomRaw) > 0 {
			if geom, _, err := ewkb.Unmarshal(geomRaw); err != nil {
				slog.Warn("failed to unmarshal photo geometry", "id", ph.ID, "error", err)
			} else if p, ok := geom.(orb.Point); ok {
				ph.Geom = p
			}
		}
		items = append(items, ph)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitPhotoRepo) Create(ctx context.Context, ph *operations.VisitPhoto) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_photos (id, visit_id, checkpoint_id, url, geom, timestamp, stage, finding_type, created_at)
		 VALUES ($1, $2, $3, $4, ST_GeomFromEWKB($5), $6, $7, $8, $9)`,
		ph.ID, ph.VisitID, ph.CheckpointID, ph.URL, ewkb.Value(ph.Geom, 4326),
		ph.Timestamp, ph.Stage, ph.FindingType, ph.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit photo: %w", err)
	}
	return nil
}

func (r *VisitPhotoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_photos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit photo: %w", err)
	}
	return nil
}
