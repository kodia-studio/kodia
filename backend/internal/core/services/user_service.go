package services

import (
	"context"
	"fmt"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/hash"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"go.uber.org/zap"
)

// UserService implements ports.UserService.
type UserService struct {
	userRepo ports.UserRepository
	log      *zap.Logger
}

// NewUserService creates a new UserService with its dependencies injected.
func NewUserService(userRepo ports.UserRepository, log *zap.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		log:      log,
	}
}

// GetByID fetches a single user by ID.
func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAll returns a paginated list of all users.
func (s *UserService) GetAll(ctx context.Context, params *pagination.Params) ([]*domain.User, int64, error) {
	return s.userRepo.FindAll(ctx, params)
}

// Update updates a user's mutable profile fields.
func (s *UserService) Update(ctx context.Context, id string, input ports.UpdateUserInput) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Error("Failed to update user", zap.String("user_id", id), zap.Error(err))
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

// Delete soft-deletes a user by ID.
func (s *UserService) Delete(ctx context.Context, id string) error {
	if _, err := s.userRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, id)
}

// ChangePassword updates a user's password after verifying the current one.
func (s *UserService) ChangePassword(ctx context.Context, id string, input ports.ChangePasswordInput) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !hash.Check(input.CurrentPassword, user.Password) {
		return domain.ErrInvalidCredentials
	}

	newHashedPassword, err := hash.Make(input.NewPassword)
	if err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	user.Password = newHashedPassword
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

// UpdateAvatar updates the avatar URL for a user.
func (s *UserService) UpdateAvatar(ctx context.Context, id string, avatarURL string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	user.AvatarURL = &avatarURL
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}
