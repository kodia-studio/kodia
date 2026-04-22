package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"gorm.io/gorm"
)

// gormNotification is the GORM model for the notifications table.
// It mirrors the domain.Notification entity but with GORM-specific tags.
type gormNotification struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id;not null;index"`
	Type      string    `gorm:"column:type;not null"`
	Title     string    `gorm:"column:title;not null"`
	Message   string    `gorm:"column:message;not null"`
	Data      []byte    `gorm:"column:data;type:jsonb"`
	IsRead    bool      `gorm:"column:is_read;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (gormNotification) TableName() string { return "notifications" }

// toDomain converts a gormNotification to a domain.Notification entity.
func (g *gormNotification) toDomain() *domain.Notification {
	var data map[string]interface{}
	if len(g.Data) > 0 {
		json.Unmarshal(g.Data, &data)
	}
	return &domain.Notification{
		ID:        g.ID,
		UserID:    g.UserID,
		Type:      domain.NotificationType(g.Type),
		Title:     g.Title,
		Message:   g.Message,
		Data:      data,
		IsRead:    g.IsRead,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

// fromDomainNotification converts a domain.Notification to a gormNotification.
func fromDomainNotification(n *domain.Notification) *gormNotification {
	data, _ := json.Marshal(n.Data)
	return &gormNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      string(n.Type),
		Title:     n.Title,
		Message:   n.Message,
		Data:      data,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) ports.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	m := fromDomainNotification(n)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *notificationRepository) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	var item gormNotification
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return item.toDomain(), nil
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID string, params *pagination.Params) ([]*domain.Notification, int64, error) {
	var items []gormNotification
	var total int64

	query := r.db.WithContext(ctx).Model(&gormNotification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(params.Offset()).Limit(params.Limit()).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*domain.Notification, len(items))
	for i, item := range items {
		result[i] = item.toDomain()
	}

	return result, total, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id string, userID string) error {
	return r.db.WithContext(ctx).Model(&gormNotification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&gormNotification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) Delete(ctx context.Context, id string, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&gormNotification{}).Error
}

func (r *notificationRepository) CountUnread(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&gormNotification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}
