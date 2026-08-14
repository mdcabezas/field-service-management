package postgresnotifications

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/notifications"
	"localis-backend/internal/model/shared"
)

func TestNotificationRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewNotificationRepo(pool)

	notif := &notifications.Notification{
		ID:        uuid.New(),
		Type:      shared.NotificationTypeOther,
		Channel:   shared.NotificationChannelEmail,
		Recipient: "test_" + uuid.New().String() + "@example.com",
		Subject:   strPtr("Test Subject"),
		Body:      strPtr("Test Body"),
		Status:    shared.NotificationStatusPending,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, notif)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, notif.ID)

	got, err := repo.GetByID(ctx, notif.ID)
	require.NoError(t, err)
	require.Equal(t, notif.Type, got.Type)
	require.Equal(t, notif.Channel, got.Channel)
	require.Equal(t, notif.Recipient, got.Recipient)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, n := range result.Items {
		if n.ID == notif.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	notif.Status = shared.NotificationStatusSent
	notif.SentAt = timePtr(time.Now())
	err = repo.Update(ctx, notif.ID, notif)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, notif.ID)
	require.NoError(t, err)
	require.Equal(t, shared.NotificationStatusSent, updated.Status)

	err = repo.Delete(ctx, notif.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, notif.ID)
	require.Error(t, err)
}

func TestNotificationRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewNotificationRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestNotificationTemplateRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewNotificationTemplateRepo(pool)

	tmpl := &notifications.NotificationTemplate{
		ID:        uuid.New(),
		Type:      shared.NotificationTypeVisitReminder,
		Channel:   shared.NotificationChannelEmail,
		Subject:   strPtr("Template " + uuid.New().String()),
		Body:      strPtr("Test template body"),
		Active:    true,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, tmpl)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, tmpl.ID)

	got, err := repo.GetByID(ctx, tmpl.ID)
	require.NoError(t, err)
	require.Equal(t, tmpl.Type, got.Type)
	require.Equal(t, tmpl.Channel, got.Channel)
	require.Equal(t, tmpl.Subject, got.Subject)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, t := range result.Items {
		if t.ID == tmpl.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	tmpl.Active = false
	err = repo.Update(ctx, tmpl.ID, tmpl)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, tmpl.ID)
	require.NoError(t, err)
	require.False(t, updated.Active)

	err = repo.Delete(ctx, tmpl.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, tmpl.ID)
	require.Error(t, err)
}

func TestNotificationTemplateRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewNotificationTemplateRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
