package config

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/tasks/v1"
)

// Client ID - injected at build time or via environment variable
// Build with: go build -ldflags "-X github.com/BRO3886/gtasks/internal/config.ClientID=your-client-id"
var ClientID = ""

// Client Secret - injected at build time or via environment variable
// Note: For Google OAuth2, even "public" desktop clients require a client secret
var ClientSecret = ""

// GetOAuth2Config creates OAuth2 configuration for Google Tasks API
func GetOAuth2Config() (*oauth2.Config, error) {
	// Try environment variable first for client ID
	clientID := os.Getenv("GTASKS_CLIENT_ID")
	if clientID == "" {
		// Fall back to build-time injected client ID
		clientID = ClientID
	}

	if clientID == "" {
		return nil, fmt.Errorf("no client ID found. Set GTASKS_CLIENT_ID environment variable or rebuild with client ID")
	}

	// Try environment variable first for client secret
	clientSecret := os.Getenv("GTASKS_CLIENT_SECRET")
	if clientSecret == "" {
		// Fall back to build-time injected client secret
		clientSecret = ClientSecret
	}

	if clientSecret == "" {
		return nil, fmt.Errorf("no client secret found. Set GTASKS_CLIENT_SECRET environment variable or rebuild with client secret")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{tasks.TasksScope},
		Endpoint:     google.Endpoint,
		// RedirectURL will be set dynamically by auth flow
	}

	return config, nil
}

// ValidateOAuth2Config ensures the OAuth2 configuration is valid
func ValidateOAuth2Config(config *oauth2.Config) error {
	if config.ClientID == "" {
		return fmt.Errorf("OAuth2 client ID is required")
	}
	if len(config.Scopes) == 0 {
		return fmt.Errorf("OAuth2 scopes are required")
	}
	if config.Endpoint.AuthURL == "" || config.Endpoint.TokenURL == "" {
		return fmt.Errorf("OAuth2 endpoint is invalid")
	}
	return nil
}
