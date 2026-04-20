// Package dto contains Data Transfer Objects for HTTP request/response in Kodia Framework.
// DTOs validate and sanitize input before passing it to the service layer.
package dto

// --- Auth DTOs ---

// RegisterRequest is the request body for POST /api/auth/register.
// @swagger:model
type RegisterRequest struct {
	// User display name
	// @required
	// @example "John Doe"
	Name string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`

	// User email address
	// @required
	// @example "user@example.com"
	Email string `json:"email" validate:"required,email" example:"user@example.com"`

	// User password
	// @required
	// @example "SecurePassword123!"
	// @minLength 8
	Password string `json:"password" validate:"required,min=8,max=72" example:"SecurePassword1123!"`
}

// LoginRequest is the request body for POST /api/auth/login.
// @swagger:model
type LoginRequest struct {
	// User email address
	// @required
	// @example "user@example.com"
	Email string `json:"email" validate:"required,email" example:"user@example.com"`

	// User password
	// @required
	Password string `json:"password" validate:"required" example:"SecurePassword123!"`
}

// RefreshTokenRequest is the request body for POST /api/auth/refresh.
// @swagger:model
type RefreshTokenRequest struct {
	// Valid refresh token
	// @required
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// LogoutRequest is the request body for POST /api/auth/logout.
// @swagger:model
type LogoutRequest struct {
	// Refresh token to revoke
	// @required
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// AuthResponse is the response for successful authentication.
// @swagger:model
type AuthResponse struct {
	AccessToken  string       `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string       `json:"token_type"    example:"Bearer"`
	User         UserResponse `json:"user"`
}

// --- User DTOs ---

// UserResponse is the public-safe user representation.
// Never expose the password field.
// @swagger:model
type UserResponse struct {
	// User unique identifier
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// User display name
	Name string `json:"name" example:"John Doe"`

	// User email address
	Email string `json:"email" example:"user@example.com"`

	// User role
	Role string `json:"role" example:"user"`

	// User account status
	IsActive bool `json:"is_active" example:"true"`

	// User avatar image URL
	AvatarURL *string `json:"avatar_url" example:"https://example.com/avatar.jpg"`

	// ISO 8601 creation timestamp
	CreatedAt string `json:"created_at" example:"2024-04-19T10:00:00Z"`

	// ISO 8601 update timestamp
	UpdatedAt string `json:"updated_at" example:"2024-04-19T10:00:00Z"`
}

// UpdateUserRequest is the request body for PATCH /api/users/:id.
// @swagger:model
type UpdateUserRequest struct {
	Name      *string `json:"name"       validate:"omitempty,min=2,max=100" example:"Jane Doe"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"             example:"https://example.com/avatar2.jpg"`
}

// ChangePasswordRequest is the request body for POST /api/users/me/change-password.
// @swagger:model
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"OldSecurePassword123!"`
	NewPassword     string `json:"new_password"     validate:"required,min=8,max=72" example:"NewSecurePassword456!"`
}

// PaginatedUsersResponse wraps a list of users with pagination metadata.
// @swagger:model
type PaginatedUsersResponse struct {
	Users []UserResponse `json:"users"`
}
