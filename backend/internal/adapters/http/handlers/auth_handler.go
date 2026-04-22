// Package handlers contains all HTTP handler implementations for Kodia Framework.
package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/dto"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/binder"
	"github.com/kodia-studio/kodia/pkg/response"
	"github.com/kodia-studio/kodia/pkg/validation"
	"go.uber.org/zap"
)

// AuthHandler handles all authentication-related HTTP requests.
type AuthHandler struct {
	authService ports.AuthService
	validate    *validation.Validator
	log         *zap.Logger
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService ports.AuthService, validate *validation.Validator, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validate,
		log:         log,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account and return JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RegisterRequest true "Registration data"
// @Success      201 {object} response.Response{data=dto.AuthResponse}
// @Failure      400 {object} response.Response
// @Failure      409 {object} response.Response
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := binder.Bind(c, &req); err != nil {
		return
	}

	result, err := h.authService.Register(c.Request.Context(), ports.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.Created(c, "Registration successful", dto.MapAuthToResponse(result))
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password, returns JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginRequest true "Login credentials"
// @Success      200 {object} response.Response{data=dto.AuthResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := binder.Bind(c, &req); err != nil {
		return
	}

	result, err := h.authService.Login(c.Request.Context(), ports.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, "Login successful", dto.MapAuthToResponse(result))
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Exchange a valid refresh token for a new access token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} response.Response{data=dto.AuthResponse}
// @Failure      401 {object} response.Response
// @Router       /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	response.OK(c, "Token refreshed", dto.MapAuthToResponse(result))
}

// Logout godoc
// @Summary      Logout
// @Description  Revoke a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body dto.LogoutRequest true "Refresh token to revoke"
// @Success      200 {object} response.Response
// @Router       /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		h.log.Warn("Logout error", zap.Error(err))
	}

	response.OK(c, "Logged out successfully", nil)
}

// LogoutAll godoc
// @Summary      Logout from all devices
// @Description  Revoke all refresh tokens for the authenticated user
// @Tags         auth
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Router       /api/auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.authService.LogoutAll(c.Request.Context(), userID); err != nil {
		response.InternalServerError(c, "")
		return
	}
	response.OK(c, "Logged out from all devices", nil)
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the profile of the authenticated user
// @Tags         auth
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=dto.UserResponse}
// @Router       /api/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// The user info is already in the context from the Auth middleware
	userID := middleware.GetUserID(c)
	c.Set("requesting_user_id", userID)
	// Forward to user handler is one pattern; returning claims directly is another.
	// Here we return a simple profile from JWT claims for efficiency.
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved",
		"data": gin.H{
			"user_id": userID,
			"role":    middleware.GetUserRole(c),
		},
	})
}

// ForgotPassword handles requesting a password reset link.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		h.log.Error("Forgot password error", zap.Error(err))
	}

	response.OK(c, "If your email is registered, you will receive a reset link.", nil)
}

// ResetPassword handles resetting a password using a token.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !validation.BindAndValidate(c, h.validate, &req) {
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Password reset successfully", nil)
}

// VerifyEmail handles email verification using a token.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.BadRequest(c, "Token is required", nil)
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), token); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "Email verified successfully", nil)
}

// Enable2FA generates TOTP secret and QR code.
func (h *AuthHandler) Enable2FA(c *gin.Context) {
	userID := middleware.GetUserID(c)
	setup, err := h.authService.Enable2FA(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to initiate 2FA setup")
		return
	}

	response.OK(c, "2FA setup initiated", dto.TwoFactorSetupResponse{
		Secret: setup.Secret,
		QRCode: setup.QRCode,
	})
}

// Verify2FA verifies the initial TOTP setup and returns recovery codes.
func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req dto.Verify2FARequest
	if !validation.BindAndValidate(c, h.validate, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	recovery, err := h.authService.Verify2FA(c.Request.Context(), userID, req.Code)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.OK(c, "2FA enabled successfully. Save your recovery codes.", gin.H{"recovery_codes": recovery})
}

// Disable2FA disables TOTP for the user.
func (h *AuthHandler) Disable2FA(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.authService.Disable2FA(c.Request.Context(), userID); err != nil {
		response.InternalServerError(c, "Failed to disable 2FA")
		return
	}

	response.OK(c, "2FA disabled successfully", nil)
}

// LoginVerify2FA handles TOTP verification during login using a temporary token.
func (h *AuthHandler) LoginVerify2FA(c *gin.Context) {
	var req dto.LoginVerify2FARequest
	if !validation.BindAndValidate(c, h.validate, &req) {
		return
	}

	result, err := h.authService.LoginVerify2FA(c.Request.Context(), req.MFAToken, req.Code)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "Login successful", dto.MapAuthToResponse(result))
}

// handleAuthError maps domain errors to appropriate HTTP responses.
func (h *AuthHandler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrAlreadyExists):
		response.Conflict(c, "Email address is already registered")
	case errors.Is(err, domain.ErrInvalidCredentials):
		response.Unauthorized(c, "Invalid email or password")
	case errors.Is(err, domain.ErrInactiveAccount):
		response.Forbidden(c, "Your account has been deactivated")
	case errors.Is(err, domain.ErrTokenExpired):
		response.Unauthorized(c, "Token has expired")
	case errors.Is(err, domain.ErrTokenRevoked):
		response.Unauthorized(c, "Token has been revoked")
	default:
		h.log.Error("Auth error", zap.Error(err))
		response.InternalServerError(c, "")
	}
}

