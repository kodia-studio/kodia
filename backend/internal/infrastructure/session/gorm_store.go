package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/pkg/session"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

// NewGormStore creates a new session store that persists to database via GORM.
func NewGormStore(db *gorm.DB) *GormStore {
	// AutoMigrate the session table
	_ = db.AutoMigrate(&session.Session{})
	
	return &GormStore{db: db}
}

func (s *GormStore) Create(sess *session.Session) error {
	if sess.ID == "" {
		sess.ID = uuid.NewString()
	}
	if sess.CreatedAt.IsZero() {
		sess.CreatedAt = time.Now()
	}
	sess.LastSeen = time.Now()
	return s.db.Create(sess).Error
}

func (s *GormStore) Get(id string) (*session.Session, error) {
	var sess session.Session
	err := s.db.Where("id = ? AND is_expired = ?", id, false).First(&sess).Error
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *GormStore) GetByUserID(userID string) ([]session.Session, error) {
	var sessions []session.Session
	err := s.db.Where("user_id = ? AND is_expired = ?", userID, false).Find(&sessions).Error
	return sessions, err
}

func (s *GormStore) UpdateLastSeen(id string) error {
	return s.db.Model(&session.Session{}).Where("id = ?", id).Update("last_seen", time.Now()).Error
}

func (s *GormStore) Revoke(id string) error {
	return s.db.Model(&session.Session{}).Where("id = ?", id).Update("is_expired", true).Error
}

func (s *GormStore) RevokeAllForUser(userID string) error {
	return s.db.Model(&session.Session{}).Where("user_id = ?", userID).Update("is_expired", true).Error
}
