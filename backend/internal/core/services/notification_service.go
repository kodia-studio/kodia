package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	wsocket "github.com/kodia-studio/kodia/internal/adapters/websocket"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"go.uber.org/zap"
)

// NotificationCreatedEvent is the event dispatched when a notification is created.
type NotificationCreatedEvent struct {
	notification *domain.Notification
}

func NewNotificationCreatedEvent(n *domain.Notification) *NotificationCreatedEvent {
	return &NotificationCreatedEvent{notification: n}
}

func (e *NotificationCreatedEvent) Name() string {
	return "NotificationCreated"
}

func (e *NotificationCreatedEvent) Payload() interface{} {
	return e.notification
}

type notificationService struct {
	repo        ports.NotificationRepository
	broadcaster *wsocket.Broadcaster
	dispatcher  ports.EventDispatcher
	log         *zap.Logger
}

func NewNotificationService(
	repo ports.NotificationRepository,
	broadcaster *wsocket.Broadcaster,
	dispatcher ports.EventDispatcher,
	log *zap.Logger,
) ports.NotificationService {
	return &notificationService{
		repo:        repo,
		broadcaster: broadcaster,
		dispatcher:  dispatcher,
		log:         log,
	}
}

func (s *notificationService) Send(ctx context.Context, input ports.SendNotificationInput) (*domain.Notification, error) {
	n := &domain.Notification{
		ID:        uuid.NewString(),
		UserID:    input.UserID,
		Type:      input.Type,
		Title:     input.Title,
		Message:   input.Message,
		Data:      input.Data,
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 1. Persist to database
	if err := s.repo.Create(ctx, n); err != nil {
		s.log.Error("Failed to create notification", zap.String("userID", input.UserID), zap.Error(err))
		return nil, err
	}

	// 2. Real-time push via WebSocket
	s.broadcaster.NotifyUser(input.UserID, "notification", wsocket.NotificationPayload{
		Title:   n.Title,
		Message: n.Message,
		Data:    n.Data,
	})

	// 3. Async email (optional, via event dispatch)
	if input.SendEmail {
		s.dispatcher.Dispatch(ctx, NewNotificationCreatedEvent(n))
	}

	return n, nil
}

func (s *notificationService) GetAll(ctx context.Context, userID string, params *pagination.Params) ([]*domain.Notification, int64, error) {
	return s.repo.FindByUserID(ctx, userID, params)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id string, userID string) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) Delete(ctx context.Context, id string, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *notificationService) CountUnread(ctx context.Context, userID string) (int64, error) {
	return s.repo.CountUnread(ctx, userID)
}
