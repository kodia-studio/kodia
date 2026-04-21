// Package authsocial provides standardized abstractions for OAuth2 and Social Login providers.
package authsocial

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *GoogleProvider) GetUser(ctx context.Context, token *oauth2.Token) (*User, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &User{
		ID:        fmt.Sprintf("%v", data["id"]),
		Email:     fmt.Sprintf("%v", data["email"]),
		Name:      fmt.Sprintf("%v", data["name"]),
		AvatarURL: fmt.Sprintf("%v", data["picture"]),
		RawData:   data,
	}, nil
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

func (p *GitHubProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *GitHubProvider) GetUser(ctx context.Context, token *oauth2.Token) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userData map[string]interface{}
	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, err
	}

	user := &User{
		ID:        fmt.Sprintf("%v", userData["id"]),
		Email:     fmt.Sprintf("%v", userData["email"]),
		Name:      fmt.Sprintf("%v", userData["name"]),
		AvatarURL: fmt.Sprintf("%v", userData["avatar_url"]),
		RawData:   userData,
	}

	// If email is nil from user endpoint, fetch from emails endpoint
	if user.Email == "<nil>" || user.Email == "" {
		req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				defer resp.Body.Close()
				var emails []map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&emails); err == nil {
					for _, email := range emails {
						if primary, ok := email["primary"].(bool); ok && primary {
							user.Email = fmt.Sprintf("%v", email["email"])
							break
						}
					}
				}
			}
		}
	}

	return user, nil
}
