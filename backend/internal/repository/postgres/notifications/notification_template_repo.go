package postgresnotifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/notifications"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type NotificationTemplateRepo struct {
	pool *pgxpool.Pool
}

func NewNotificationTemplateRepo(pool *pgxpool.Pool) *NotificationTemplateRepo {
	return &NotificationTemplateRepo{pool: pool}
}

func (r *NotificationTemplateRepo) GetByID(ctx context.Context, id uuid.UUID) (*notifications.NotificationTemplate, error) {
	var t notifications.NotificationTemplate
	err := r.pool.QueryRow(ctx,
		`SELECT id, type, channel, subject, body, active, created_at
		 FROM notifications.notification_templates
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.Type, &t.Channel, &t.Subject, &t.Body, &t.Active, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "notification_template", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get notification_template by id: %w", err)
	}
	return &t, nil
}

func (r *NotificationTemplateRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[notifications.NotificationTemplate], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications.notification_templates`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count notification templates: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, type, channel, subject, body, active, created_at
		 FROM notifications.notification_templates
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list notification templates: %w", err)
	}
	defer rows.Close()

	items := make([]notifications.NotificationTemplate, 0)
	for rows.Next() {
		var t notifications.NotificationTemplate
		if err := rows.Scan(&t.ID, &t.Type, &t.Channel, &t.Subject, &t.Body, &t.Active, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan notification template: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[notifications.NotificationTemplate]{Items: items, Total: total}, nil
}

func (r *NotificationTemplateRepo) Create(ctx context.Context, tmpl *notifications.NotificationTemplate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notifications.notification_templates (id, type, channel, subject, body, active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		tmpl.ID, tmpl.Type, tmpl.Channel, tmpl.Subject, tmpl.Body, tmpl.Active, tmpl.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create notification template: %w", err)
	}
	return nil
}

func (r *NotificationTemplateRepo) Update(ctx context.Context, id uuid.UUID, tmpl *notifications.NotificationTemplate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notifications.notification_templates
		 SET type = $2, channel = $3, subject = $4, body = $5, active = $6
		 WHERE id = $1`,
		id, tmpl.Type, tmpl.Channel, tmpl.Subject, tmpl.Body, tmpl.Active,
	)
	if err != nil {
		return fmt.Errorf("update notification template: %w", err)
	}
	return nil
}

func (r *NotificationTemplateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM notifications.notification_templates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete notification template: %w", err)
	}
	return nil
}
