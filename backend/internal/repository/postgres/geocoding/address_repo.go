package postgresgeocoding

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
	"localis-backend/internal/model/geocoding"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type GeocodingAddressRepo struct {
	pool *pgxpool.Pool
}

func NewGeocodingAddressRepo(pool *pgxpool.Pool) *GeocodingAddressRepo {
	return &GeocodingAddressRepo{pool: pool}
}

func (r *GeocodingAddressRepo) GetByID(ctx context.Context, id uuid.UUID) (*geocoding.Address, error) {
	var a geocoding.Address
	var geomRaw []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, street, number, apartment, neighborhood, city, region,
		        location_references, postal_code, ST_AsEWKB(geom), created_at
		 FROM geocoding.addresses
		 WHERE id = $1`, id,
	).Scan(&a.ID, &a.Street, &a.Number, &a.Apartment, &a.Neighborhood,
		&a.City, &a.Region, &a.LocationReferences, &a.PostalCode,
		&geomRaw, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "address", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get address by id: %w", err)
	}
	if len(geomRaw) > 0 {
		geom, _, err := ewkb.Unmarshal(geomRaw)
		if err != nil {
			slog.Warn("failed to unmarshal address geometry", "id", a.ID, "error", err)
		} else if p, ok := geom.(orb.Point); ok {
			a.Geom = p
		}
	}
	return &a, nil
}

func (r *GeocodingAddressRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[geocoding.Address], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM geocoding.addresses`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count addresses: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, street, number, apartment, neighborhood, city, region,
		        location_references, postal_code, ST_AsEWKB(geom), created_at
		 FROM geocoding.addresses
		 ORDER BY street, number
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list addresses: %w", err)
	}
	defer rows.Close()

	items := make([]geocoding.Address, 0)
	for rows.Next() {
		var a geocoding.Address
		var geomRaw []byte
		if err := rows.Scan(&a.ID, &a.Street, &a.Number, &a.Apartment, &a.Neighborhood,
			&a.City, &a.Region, &a.LocationReferences, &a.PostalCode,
			&geomRaw, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan address: %w", err)
		}
		if len(geomRaw) > 0 {
			geom, _, err := ewkb.Unmarshal(geomRaw)
			if err != nil {
				slog.Warn("failed to unmarshal address geometry", "id", a.ID, "error", err)
			} else if p, ok := geom.(orb.Point); ok {
				a.Geom = p
			}
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[geocoding.Address]{Items: items, Total: total}, nil
}

func (r *GeocodingAddressRepo) Create(ctx context.Context, addr *geocoding.Address) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO geocoding.addresses
		     (id, street, number, apartment, neighborhood, city, region,
		      location_references, postal_code, geom, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, ST_GeomFromEWKB($10), $11)`,
		addr.ID, addr.Street, addr.Number, addr.Apartment, addr.Neighborhood,
		addr.City, addr.Region, addr.LocationReferences, addr.PostalCode,
		ewkb.Value(addr.Geom, 4326), addr.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create address: %w", err)
	}
	return nil
}

func (r *GeocodingAddressRepo) Update(ctx context.Context, id uuid.UUID, addr *geocoding.Address) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE geocoding.addresses
		 SET street = $2, number = $3, apartment = $4, neighborhood = $5,
		     city = $6, region = $7, location_references = $8, postal_code = $9,
		     geom = ST_GeomFromEWKB($10)
		 WHERE id = $1`,
		id, addr.Street, addr.Number, addr.Apartment, addr.Neighborhood,
		addr.City, addr.Region, addr.LocationReferences, addr.PostalCode,
		ewkb.Value(addr.Geom, 4326),
	)
	if err != nil {
		return fmt.Errorf("update address: %w", err)
	}
	return nil
}

func (r *GeocodingAddressRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM geocoding.addresses WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete address: %w", err)
	}
	return nil
}

func (r *GeocodingAddressRepo) FindNearby(ctx context.Context, lat, lng, radiusMeters float64) ([]geocoding.Address, error) {
	const maxFindNearbyLimit = 100

	rows, err := r.pool.Query(ctx,
		`SELECT id, street, number, apartment, neighborhood, city, region,
		        location_references, postal_code, ST_AsEWKB(geom), created_at
		 FROM geocoding.addresses
		 WHERE geom IS NOT NULL AND ST_DWithin(
		     geom::geography,
		     ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
		     $3
		 )
		 ORDER BY ST_Distance(geom::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography)
		 LIMIT $4`,
		lng, lat, radiusMeters, maxFindNearbyLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("find nearby addresses: %w", err)
	}
	defer rows.Close()

	items := make([]geocoding.Address, 0)
	for rows.Next() {
		var a geocoding.Address
		var geomRaw []byte
		if err := rows.Scan(&a.ID, &a.Street, &a.Number, &a.Apartment, &a.Neighborhood,
			&a.City, &a.Region, &a.LocationReferences, &a.PostalCode,
			&geomRaw, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan nearby address: %w", err)
		}
		if len(geomRaw) > 0 {
			geom, _, err := ewkb.Unmarshal(geomRaw)
			if err != nil {
				slog.Warn("failed to unmarshal nearby address geometry", "id", a.ID, "error", err)
			} else if p, ok := geom.(orb.Point); ok {
				a.Geom = p
			}
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}
