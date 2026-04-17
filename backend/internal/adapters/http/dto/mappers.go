package dto

import (
	"time"

	"github.com/kodia/framework/backend/internal/core/domain"
	"github.com/kodia/framework/backend/internal/core/ports"
)

// MapUserToResponse converts a domain.User to the public-safe UserResponse DTO.
func MapUserToResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// MapAuthToResponse converts a ports.AuthResponse to the HTTP AuthResponse DTO.
func MapAuthToResponse(r *ports.AuthResponse) AuthResponse {
	return AuthResponse{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		TokenType:    "Bearer",
		User:         MapUserToResponse(r.User),
	}
}

// MapUsersToResponse converts a slice of domain.User to a slice of UserResponse DTOs.
func MapUsersToResponse(users []*domain.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = MapUserToResponse(u)
	}
	return result
}
