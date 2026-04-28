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
// @Router       /auth/register [post]
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
// @Router       /auth/login [post]
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
// @Router       /auth/refresh [post]
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

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Send a password reset link to the user's email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.ForgotPasswordRequest true "Email address"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /auth/forgot-password [post]
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

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset password using a valid reset token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.ResetPasswordRequest true "Reset token and new password"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /auth/reset-password [post]
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

// VerifyEmail godoc
// @Summary      Verify email address
// @Description  Verify email using a verification token from email
// @Tags         auth
// @Produce      json
// @Param        token query string true "Verification token"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /auth/verify-email [get]
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

// Logout godoc
// @Summary      Logout
// @Description  Revoke a refresh token
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.LogoutRequest true "Refresh token to revoke"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/logout [post]
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
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.authService.LogoutAll(c.Request.Context(), userID); err != nil {
		response.InternalServerError(c, "")
		return
	}
	response.OK(c, "Logged out from all devices", nil)
}

// Me godoc
// @Summary      Get current user profile
// @Description  Returns the profile of the authenticated user
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/me [get]
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

// Enable2FA godoc
// @Summary      Enable two-factor authentication
// @Description  Generate TOTP secret and QR code for setting up 2FA
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Response{data=dto.TwoFactorSetupResponse}
// @Failure      401 {object} response.Response
// @Router       /auth/2fa/enable [post]
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

// Verify2FA godoc
// @Summary      Verify two-factor authentication setup
// @Description  Verify TOTP code during 2FA setup and return recovery codes
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.Verify2FARequest true "TOTP verification code"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/2fa/verify [post]
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

// Disable2FA godoc
// @Summary      Disable two-factor authentication
// @Description  Disable TOTP for the authenticated user
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/2fa/disable [delete]
func (h *AuthHandler) Disable2FA(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.authService.Disable2FA(c.Request.Context(), userID); err != nil {
		response.InternalServerError(c, "Failed to disable 2FA")
		return
	}

	response.OK(c, "2FA disabled successfully", nil)
}

// LoginVerify2FA godoc
// @Summary      Verify 2FA during login
// @Description  Complete login by verifying a TOTP code using temporary MFA token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginVerify2FARequest true "MFA token and verification code"
// @Success      200 {object} response.Response{data=dto.AuthResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /auth/2fa/login-verify [post]
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
