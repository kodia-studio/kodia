package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/providers"
	"github.com/kodia-studio/kodia/pkg/kodia"
	"github.com/kodia-studio/kodia/pkg/validation"
	"github.com/kodia-studio/kodia/tests"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Login(t *testing.T) {
	// 1. Setup Test Infrastructure
	tests.SkipIfShort(t)
	
	// Start DB container
	td := tests.NewTestDatabase(t)
	defer td.Cleanup()
	
	// Create Kodia App
	cfg := tests.NewTestConfig()
	log := tests.NewTestLogger()
	app := kodia.NewApp(cfg, log)
	
	// Inject test dependencies
	app.DB = td.DB
	
	// Register Providers
	app.RegisterProviders(
		providers.NewHttpProvider(),
		providers.NewAuthProvider(),
	)
	
	// Start Test Server
	ts := app.NewTestServer(t)
	factory := tests.NewFactory(t, td.DB)
	
	t.Run("it should login successfully with valid credentials", func(t *testing.T) {
		td.Reset()
		
		// Create a user
		email := "test@kodia.id"
		factory.CreateUser(func(u *domain.User) {
			u.Email = email
		})
		
		// Prepare payload
		payload := map[string]string{
			"email":    email,
			"password": "password123",
		}
		
		// Perform Request
		resp := tests.JSONRequest(t, http.DefaultClient, "POST", ts.URL+"/api/auth/login", payload)
		
		// Assertions
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		// Read body for both parsing and contract validation
		body, _ := io.ReadAll(resp.Body)
		
		var result map[string]interface{}
		_ = json.Unmarshal(body, &result)
		
		assert.True(t, result["success"].(bool))
		
		// Contract Validation (Data part)
		dataBytes, _ := json.Marshal(result["data"])
		validation.ValidateContract(t, dataBytes, map[string]interface{}{
			"access_token":  "string",
			"refresh_token": "string",
			"user":          "object",
		})
	})
}
