package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	notifmodel "localis-backend/internal/model/notifications"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type NotificationService struct {
	notifRepo repository.NotificationRepository
}

func NewNotificationService(notifRepo repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notifRepo: notifRepo,
	}
}

// SendParams defines the parameters for sending a notification.
type SendParams struct {
	Type       shared.NotificationType
	Channel    shared.NotificationChannel
	Recipient  string
	Subject    *string
	Body       *string
	EntityType *string
	EntityID   *uuid.UUID
}

// Send sends a notification using the provided parameters.
//
// Multiple configuration options may be provided:
//
//   - Subscribe to the Notifications service.
//
//   - Update correspondences.
//
//   - Make changes to notifications.
//
//     When the service is ready, notifications may be sent to customers.
//
// At the end of the method, the notification is completed.
func (s *NotificationService) Send(ctx context.Context, params SendParams) (*notifmodel.Notification, error) {
	notif := &notifmodel.Notification{
		ID:         uuid.New(),
		Type:       params.Type,
		Channel:    params.Channel,
		Recipient:  params.Recipient,
		Subject:    params.Subject,
		Body:       params.Body,
		Status:     shared.NotificationStatusPending,
		EntityType: params.EntityType,
		EntityID:   params.EntityID,
		Attempts:   0,
		CreatedAt:  time.Now(),
	}

	if err := s.notifRepo.Create(ctx, notif); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}

	if err := s.dispatch(ctx, notif); err != nil {
		notif.Status = shared.NotificationStatusFailed
		notif.Attempts++
		errMsg := err.Error()
		notif.Error = &errMsg
		if updateErr := s.notifRepo.Update(ctx, notif.ID, notif); updateErr != nil {
			slog.Warn("failed to update notification status to failed", "notif_id", notif.ID, "error", updateErr)
		}
		return notif, fmt.Errorf("dispatch notification: %w", err)
	}

	notif.Status = shared.NotificationStatusSent
	notif.Attempts++
	if updateErr := s.notifRepo.Update(ctx, notif.ID, notif); updateErr != nil {
		slog.Warn("failed to update notification status to sent", "notif_id", notif.ID, "error", updateErr)
	}

	return notif, nil
}

func (s *NotificationService) Update(ctx context.Context, id uuid.UUID, notif *notifmodel.Notification) (*notifmodel.Notification, error) {
	existing, err := s.notifRepo.GetByID(ctx, id)
	if err != nil {
		return nil, service.HandleRepoGetByIDError(err, "notification", id.String())
	}

	if notif.Type != "" {
		existing.Type = notif.Type
	}
	if notif.Channel != "" {
		existing.Channel = notif.Channel
	}
	if notif.Recipient != "" {
		existing.Recipient = notif.Recipient
	}
	if notif.Subject != nil {
		existing.Subject = notif.Subject
	}
	if notif.Body != nil {
		existing.Body = notif.Body
	}
	if notif.Status != "" {
		existing.Status = notif.Status
	}

	if err := s.notifRepo.Update(ctx, id, existing); err != nil {
		return nil, fmt.Errorf("update notification: %w", err)
	}

	return existing, nil
}

func (s *NotificationService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.notifRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "notification", id.String())
	}
	if err := s.notifRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

func (s *NotificationService) GetByID(ctx context.Context, id uuid.UUID) (*notifmodel.Notification, error) {
	item, err := s.notifRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return item, nil
}

func (s *NotificationService) List(ctx context.Context, limit, offset int) (*repository.ListResult[notifmodel.Notification], error) {
	items, err := s.notifRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	return items, nil
}

func (s *NotificationService) NotifySLAWarning(ctx context.Context, visitID uuid.UUID, slaID uuid.UUID, slaName string, deadline time.Time) error {
	subject := fmt.Sprintf("Alerta SLA: %s", slaName)
	body := fmt.Sprintf("El SLA \"%s\" vence el %s", slaName, deadline.Format("2006-01-02 15:04"))

	notifType := shared.NotificationTypeSLAWarning
	entityType := string(shared.AuditEntityTypeVisitSLATracking)

	_, err := s.Send(ctx, SendParams{
		Type:       notifType,
		Channel:    shared.NotificationChannelEmail,
		Recipient:  slaName,
		Subject:    &subject,
		Body:       &body,
		EntityType: &entityType,
		EntityID:   &visitID,
	})
	if err != nil {
		return fmt.Errorf("notify SLA warning: %w", err)
	}
	return nil
}

func (s *NotificationService) dispatch(ctx context.Context, notif *notifmodel.Notification) error {
	switch notif.Channel {
	case shared.NotificationChannelEmail:
		return s.dispatchEmail(ctx, notif)
	case shared.NotificationChannelSMS:
		return s.dispatchSMS(ctx, notif)
	case shared.NotificationChannelWhatsapp:
		return s.dispatchWhatsApp(ctx, notif)
	case shared.NotificationChannelPush:
		return s.dispatchPush(ctx, notif)
	case shared.NotificationChannelInApp:
		return s.dispatchInApp(ctx, notif)
	default:
		return fmt.Errorf("unsupported notification channel: %s", notif.Channel)
	}
}

func (s *NotificationService) dispatchEmail(ctx context.Context, notif *notifmodel.Notification) error {
	return fmt.Errorf("email delivery not implemented for notification %s", notif.ID)
}

func (s *NotificationService) dispatchSMS(ctx context.Context, notif *notifmodel.Notification) error {
	return fmt.Errorf("SMS delivery not implemented for notification %s", notif.ID)
}

func (s *NotificationService) dispatchWhatsApp(ctx context.Context, notif *notifmodel.Notification) error {
	return fmt.Errorf("WhatsApp delivery not implemented for notification %s", notif.ID)
}

func (s *NotificationService) dispatchPush(ctx context.Context, notif *notifmodel.Notification) error {
	return fmt.Errorf("push notification delivery not implemented for notification %s", notif.ID)
}

func (s *NotificationService) dispatchInApp(ctx context.Context, notif *notifmodel.Notification) error {
	return fmt.Errorf("in-app notification delivery not implemented for notification %s", notif.ID)
}
