// Package dto contains Data Transfer Objects for HTTP request/response in Kodia Framework.
// DTOs validate and sanitize input before passing it to the service layer.
package dto

// --- Auth DTOs ---

// RegisterRequest is the request body for POST /api/auth/register.
type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest is the request body for POST /api/auth/login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest is the request body for POST /api/auth/refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest is the request body for POST /api/auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse is the response for successful authentication.
type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	User         UserResponse `json:"user"`
}

// --- User DTOs ---

// UserResponse is the public-safe user representation.
// Never expose the password field.
type UserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	AvatarURL *string `json:"avatar_url"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// UpdateUserRequest is the request body for PATCH /api/users/:id.
type UpdateUserRequest struct {
	Name      *string `json:"name"       validate:"omitempty,min=2,max=100"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
}

// ChangePasswordRequest is the request body for POST /api/users/me/change-password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=8,max=72"`
}

// PaginatedUsersResponse wraps a list of users with pagination metadata.
type PaginatedUsersResponse struct {
	Users []UserResponse `json:"users"`
}
