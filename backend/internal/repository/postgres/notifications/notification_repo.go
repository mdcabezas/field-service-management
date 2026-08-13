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

type NotificationRepo struct {
	pool *pgxpool.Pool
}

func NewNotificationRepo(pool *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{pool: pool}
}

func (r *NotificationRepo) GetByID(ctx context.Context, id uuid.UUID) (*notifications.Notification, error) {
	var n notifications.Notification
	err := r.pool.QueryRow(ctx,
		`SELECT id, type, channel, recipient, subject, body, status, entity_type, entity_id,
		        scheduled_for, sent_at, read_at, attempts, error, created_at
		 FROM notifications.notifications
		 WHERE id = $1`, id,
	).Scan(&n.ID, &n.Type, &n.Channel, &n.Recipient, &n.Subject, &n.Body, &n.Status, &n.EntityType, &n.EntityID,
		&n.ScheduledFor, &n.SentAt, &n.ReadAt, &n.Attempts, &n.Error, &n.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "notification", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get notification by id: %w", err)
	}
	return &n, nil
}

func (r *NotificationRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[notifications.Notification], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications.notifications`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count notifications: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, type, channel, recipient, subject, body, status, entity_type, entity_id,
		        scheduled_for, sent_at, read_at, attempts, error, created_at
		 FROM notifications.notifications
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	items := make([]notifications.Notification, 0)
	for rows.Next() {
		var n notifications.Notification
		if err := rows.Scan(&n.ID, &n.Type, &n.Channel, &n.Recipient, &n.Subject, &n.Body, &n.Status, &n.EntityType, &n.EntityID,
			&n.ScheduledFor, &n.SentAt, &n.ReadAt, &n.Attempts, &n.Error, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		items = append(items, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[notifications.Notification]{Items: items, Total: total}, nil
}

func (r *NotificationRepo) Create(ctx context.Context, notif *notifications.Notification) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notifications.notifications (id, type, channel, recipient, subject, body, status, entity_type, entity_id,
		        scheduled_for, sent_at, read_at, attempts, error, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		notif.ID, notif.Type, notif.Channel, notif.Recipient, notif.Subject, notif.Body, notif.Status, notif.EntityType, notif.EntityID,
		notif.ScheduledFor, notif.SentAt, notif.ReadAt, notif.Attempts, notif.Error, notif.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepo) Update(ctx context.Context, id uuid.UUID, notif *notifications.Notification) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notifications.notifications
		 SET type = $2, channel = $3, recipient = $4, subject = $5, body = $6, status = $7,
		     entity_type = $8, entity_id = $9, scheduled_for = $10, sent_at = $11, read_at = $12,
		     attempts = $13, error = $14
		 WHERE id = $1`,
		id, notif.Type, notif.Channel, notif.Recipient, notif.Subject, notif.Body, notif.Status, notif.EntityType, notif.EntityID,
		notif.ScheduledFor, notif.SentAt, notif.ReadAt, notif.Attempts, notif.Error,
	)
	if err != nil {
		return fmt.Errorf("update notification: %w", err)
	}
	return nil
}

func (r *NotificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM notifications.notifications WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}
