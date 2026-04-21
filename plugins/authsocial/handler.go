package authsocial

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/response"
	"go.uber.org/zap"
)

// SocialHandler handles OAuth redirect and callback endpoints.
type SocialHandler struct {
	service             *SocialAuthService
	frontendURL         string
	googleProvider      Provider
	githubProvider      Provider
	log                 *zap.Logger
}

// NewSocialHandler creates a new social handler.
func NewSocialHandler(
	service *SocialAuthService,
	frontendURL string,
	googleProvider Provider,
	githubProvider Provider,
	log *zap.Logger,
) *SocialHandler {
	return &SocialHandler{
		service:        service,
		frontendURL:    frontendURL,
		googleProvider: googleProvider,
		githubProvider: githubProvider,
		log:            log,
	}
}

// Redirect initiates the OAuth flow by redirecting to the provider.
// GET /api/auth/social/:provider/redirect
func (h *SocialHandler) Redirect(c *gin.Context) {
	providerName := c.Param("provider")

	var provider Provider
	switch providerName {
	case "google":
		provider = h.googleProvider
	case "github":
		provider = h.githubProvider
	default:
		c.JSON(http.StatusBadRequest, response.BadRequest("unsupported provider"))
		return
	}

	// Generate CSRF state
	state, err := h.service.GenerateState(c.Request.Context())
	if err != nil {
		h.log.Error("failed to generate state", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.InternalServerError("failed to initiate login"))
		return
	}

	// Get authorization URL
	authURL := provider.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback handles the OAuth callback from the provider.
// GET /api/auth/social/:provider/callback
func (h *SocialHandler) Callback(c *gin.Context) {
	providerName := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	// Check for OAuth errors
	if errMsg := c.Query("error"); errMsg != "" {
		h.log.Warn("oauth error from provider",
			zap.String("provider", providerName),
			zap.String("error", errMsg),
		)
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?message="+errMsg)
		return
	}

	if code == "" {
		c.JSON(http.StatusBadRequest, response.BadRequest("missing authorization code"))
		return
	}

	if state == "" {
		c.JSON(http.StatusBadRequest, response.BadRequest("missing state parameter"))
		return
	}

	var provider Provider
	switch providerName {
	case "google":
		provider = h.googleProvider
	case "github":
		provider = h.githubProvider
	default:
		c.JSON(http.StatusBadRequest, response.BadRequest("unsupported provider"))
		return
	}

	// Handle OAuth callback
	result, err := h.service.HandleCallback(c.Request.Context(), providerName, code, state, provider)
	if err != nil {
		h.log.Error("oauth callback failed",
			zap.String("provider", providerName),
			zap.Error(err),
		)
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?message="+err.Error())
		return
	}

	// Redirect to frontend with tokens in query params
	redirectURL := h.frontendURL + "/auth/social/success?token=" + result.AccessToken + "&refresh=" + result.RefreshToken
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
