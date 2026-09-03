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

type MaintenanceRecordRepo struct {
	pool *pgxpool.Pool
}

func NewMaintenanceRecordRepo(pool *pgxpool.Pool) *MaintenanceRecordRepo {
	return &MaintenanceRecordRepo{pool: pool}
}

func (r *MaintenanceRecordRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.MaintenanceRecord, error) {
	var m inventory.MaintenanceRecord
	err := r.pool.QueryRow(ctx,
		`SELECT mr.id, mr.type, mr.reference_id, mr.date, mr.cost, mr.supplier, mr.description, mr.next_date, mr.created_at,
		 CASE mr.type
		   WHEN 'vehicle' THEN COALESCE(NULLIF(v.license_plate, ''), v.name)
		   WHEN 'tool' THEN COALESCE(NULLIF(t.code, ''), t.name)
		   WHEN 'equipment' THEN COALESCE(NULLIF(e.code, ''), e.name)
		 END as reference_display
		 FROM inventory.maintenance_records mr
		 LEFT JOIN inventory.vehicles v ON mr.type = 'vehicle' AND mr.reference_id = v.id
		 LEFT JOIN inventory.tools t ON mr.type = 'tool' AND mr.reference_id = t.id
		 LEFT JOIN inventory.equipment e ON mr.type = 'equipment' AND mr.reference_id = e.id
		 WHERE mr.id = $1`, id,
	).Scan(&m.ID, &m.Type, &m.ReferenceID, &m.Date, &m.Cost, &m.Supplier, &m.Description, &m.NextDate, &m.CreatedAt, &m.ReferenceDisplay)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "maintenance_record", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get maintenance_record by id: %w", err)
	}
	return &m, nil
}

func (r *MaintenanceRecordRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.MaintenanceRecord], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.maintenance_records`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count maintenance records: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT mr.id, mr.type, mr.reference_id, mr.date, mr.cost, mr.supplier, mr.description, mr.next_date, mr.created_at,
		 CASE mr.type
		   WHEN 'vehicle' THEN COALESCE(NULLIF(v.license_plate, ''), v.name)
		   WHEN 'tool' THEN COALESCE(NULLIF(t.code, ''), t.name)
		   WHEN 'equipment' THEN COALESCE(NULLIF(e.code, ''), e.name)
		 END as reference_display
		 FROM inventory.maintenance_records mr
		 LEFT JOIN inventory.vehicles v ON mr.type = 'vehicle' AND mr.reference_id = v.id
		 LEFT JOIN inventory.tools t ON mr.type = 'tool' AND mr.reference_id = t.id
		 LEFT JOIN inventory.equipment e ON mr.type = 'equipment' AND mr.reference_id = e.id
		 ORDER BY mr.date DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list maintenance records: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.MaintenanceRecord, 0)
	for rows.Next() {
		var m inventory.MaintenanceRecord
		if err := rows.Scan(&m.ID, &m.Type, &m.ReferenceID, &m.Date, &m.Cost, &m.Supplier, &m.Description, &m.NextDate, &m.CreatedAt, &m.ReferenceDisplay); err != nil {
			return nil, fmt.Errorf("scan maintenance record: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.MaintenanceRecord]{Items: items, Total: total}, nil
}

func (r *MaintenanceRecordRepo) Create(ctx context.Context, rec *inventory.MaintenanceRecord) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.maintenance_records (id, type, reference_id, date, cost, supplier, description, next_date, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		rec.ID, rec.Type, rec.ReferenceID, rec.Date, rec.Cost, rec.Supplier, rec.Description, rec.NextDate, rec.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create maintenance record: %w", err)
	}
	return nil
}

func (r *MaintenanceRecordRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.maintenance_records WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete maintenance record: %w", err)
	}
	return nil
}

func (r *MaintenanceRecordRepo) Update(ctx context.Context, id uuid.UUID, rec *inventory.MaintenanceRecord) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.maintenance_records
		 SET type = $2, reference_id = $3, date = $4, cost = $5, supplier = $6, description = $7, next_date = $8
		 WHERE id = $1`,
		id, rec.Type, rec.ReferenceID, rec.Date, rec.Cost, rec.Supplier, rec.Description, rec.NextDate,
	)
	if err != nil {
		return fmt.Errorf("update maintenance record: %w", err)
	}
	return nil
}
