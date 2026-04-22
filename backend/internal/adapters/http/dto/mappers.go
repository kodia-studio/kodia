package dto

import (
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
)

// MapUserToResponse converts a domain.User to the public-safe UserResponse DTO.
func MapUserToResponse(u *domain.User) UserResponse {
	if u == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		IsVerified: u.IsVerified,
		TwoFactorEnabled: u.TwoFactorEnabled,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// MapAuthToResponse converts a ports.AuthResponse to the HTTP AuthResponse DTO.
func MapAuthToResponse(r *ports.AuthResponse) AuthResponse {
	resp := AuthResponse{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		TokenType:    "Bearer",
		MFARequired:  r.MFARequired,
		MFAToken:     r.MFAToken,
	}

	if r.User != nil {
		resp.User = MapUserToResponse(r.User)
	}

	return resp
}

// MapUsersToResponse converts a slice of domain.User to a slice of UserResponse DTOs.
func MapUsersToResponse(users []*domain.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = MapUserToResponse(u)
	}
	return result
}

// MapNotificationToResponse converts a domain.Notification to NotificationResponse DTO.
func MapNotificationToResponse(n *domain.Notification) NotificationResponse {
	if n == nil {
		return NotificationResponse{}
	}
	return NotificationResponse{
		ID:        n.ID,
		Type:      string(n.Type),
		Title:     n.Title,
		Message:   n.Message,
		Data:      n.Data,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
}

// MapNotificationsToResponse converts a slice of domain.Notification to NotificationResponse DTOs.
func MapNotificationsToResponse(notifications []*domain.Notification) []NotificationResponse {
	result := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		result[i] = MapNotificationToResponse(n)
	}
	return result
}
