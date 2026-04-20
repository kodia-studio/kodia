// Package authsocial provides standardized abstractions for OAuth2 and Social Login providers.
package authsocial

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/github"
)

// User represents the standard user info returned by social providers.
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	RawData   map[string]interface{} `json:"-"`
}

// Provider defines the interface for OAuth2 social authentication.
type Provider interface {
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUser(ctx context.Context, token *oauth2.Token) (*User, error)
}

// Config holds the OAuth2 credentials.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// GoogleProvider handles Google OAuth2.
type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(cfg Config) *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     google.Endpoint,
		},
	}
}

func (p *GoogleProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

// GitHubProvider handles GitHub OAuth2.
type GitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider(cfg Config) *GitHubProvider {
	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     github.Endpoint,
		},
	}
}

func (p *GitHubProvider) GetAuthURL(state string) string {
    return p.config.AuthCodeURL(state)
}
