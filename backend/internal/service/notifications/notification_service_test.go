package notifications

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	notifmodel "localis-backend/internal/model/notifications"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/testutil/mocks"
)

func newNotificationService(t *testing.T) (*NotificationService, *mocks.NotificationRepository) {
	t.Helper()
	notifRepo := mocks.NewNotificationRepository(t)
	return NewNotificationService(notifRepo), notifRepo
}

func TestSend_CreateError(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	notifRepo.On("Create", ctx, mock.Anything).Return(errors.New("db down")).Once()

	_, err := svc.Send(ctx, SendParams{Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelEmail, Recipient: "a@b.cl"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create notification")
}

func TestSend_DispatchFailure(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	notifRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	notifRepo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()

	notif, err := svc.Send(ctx, SendParams{Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelEmail, Recipient: "a@b.cl"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "dispatch notification")
	require.NotNil(t, notif)
	require.Equal(t, shared.NotificationStatusFailed, notif.Status)
	require.Equal(t, 1, notif.Attempts)
	require.NotNil(t, notif.Error)
}

func TestSend_DispatchFailure_UpdateError(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	notifRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	notifRepo.On("Update", ctx, mock.Anything, mock.Anything).Return(errors.New("db down")).Once()

	_, err := svc.Send(ctx, SendParams{Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelSMS, Recipient: "5690000"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "dispatch notification")
}

func TestSend_UnsupportedChannel(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	notifRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	notifRepo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()

	_, err := svc.Send(ctx, SendParams{Type: shared.NotificationTypeOther, Channel: "pigeon", Recipient: "x"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported notification channel")
}

func TestUpdate_MergesFields(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	id := uuid.New()
	existing := &notifmodel.Notification{ID: id, Type: shared.NotificationTypeVisitAssigned, Channel: shared.NotificationChannelEmail, Recipient: "old@b.cl", Status: shared.NotificationStatusPending}
	subject := "New subject"
	body := "New body"

	notifRepo.On("GetByID", ctx, id).Return(existing, nil).Once()
	notifRepo.On("Update", ctx, id, existing).Return(nil).Once()

	updated, err := svc.Update(ctx, id, &notifmodel.Notification{
		Recipient: "new@b.cl",
		Subject:   &subject,
		Body:      &body,
		Status:    shared.NotificationStatusRead,
	})
	require.NoError(t, err)
	require.Equal(t, "new@b.cl", updated.Recipient)
	require.Equal(t, &subject, updated.Subject)
	require.Equal(t, &body, updated.Body)
	require.Equal(t, shared.NotificationStatusRead, updated.Status)
	require.Equal(t, shared.NotificationTypeVisitAssigned, updated.Type)
	require.Equal(t, shared.NotificationChannelEmail, updated.Channel)
}

func TestUpdate_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	id := uuid.New()
	notifRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "notification", ID: id.String()}).Once()

	_, err := svc.Update(ctx, id, &notifmodel.Notification{})
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
	require.Equal(t, "notification", nf.Resource)
}

func TestUpdate_RepoError(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	id := uuid.New()
	notifRepo.On("GetByID", ctx, id).Return(&notifmodel.Notification{ID: id}, nil).Once()
	notifRepo.On("Update", ctx, id, mock.Anything).Return(errors.New("db down")).Once()

	_, err := svc.Update(ctx, id, &notifmodel.Notification{Recipient: "x"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "update notification")
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	id := uuid.New()
	notifRepo.On("GetByID", ctx, id).Return(&notifmodel.Notification{ID: id}, nil).Once()
	notifRepo.On("Delete", ctx, id).Return(nil).Once()

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	id := uuid.New()
	notifRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "notification", ID: id.String()}).Once()

	err := svc.Delete(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestDelete_RepoError(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	id := uuid.New()
	notifRepo.On("GetByID", ctx, id).Return(&notifmodel.Notification{ID: id}, nil).Once()
	notifRepo.On("Delete", ctx, id).Return(errors.New("db down")).Once()

	err := svc.Delete(ctx, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete notification")
}

func TestGetByID_Success(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	id := uuid.New()
	expected := &notifmodel.Notification{ID: id}
	notifRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := svc.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGetByID_Error(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	notifRepo.On("GetByID", ctx, uuid.Nil).Return(nil, errors.New("db down")).Once()

	_, err := svc.GetByID(ctx, uuid.Nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "get notification")
}

func TestList_Success(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	expected := &repository.ListResult[notifmodel.Notification]{Items: []notifmodel.Notification{{ID: uuid.New()}}, Total: 1}
	notifRepo.On("List", ctx, 10, 0).Return(expected, nil).Once()

	got, err := svc.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestList_Error(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	notifRepo.On("List", ctx, 10, 0).Return(nil, errors.New("db down")).Once()

	_, err := svc.List(ctx, 10, 0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list notifications")
}

func TestNotifySLAWarning_Failure(t *testing.T) {
	ctx := context.Background()
	svc, notifRepo := newNotificationService(t)

	notifRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	notifRepo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil).Once()

	err := svc.NotifySLAWarning(ctx, uuid.New(), uuid.New(), "SLA Test", time.Now())
	require.Error(t, err)
	require.Contains(t, err.Error(), "notify SLA warning")
}
