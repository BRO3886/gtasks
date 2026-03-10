package config

import (
	"fmt"

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

// GetOAuth2Config creates OAuth2 configuration for Google Tasks API.
//
// Credential resolution order (highest priority first):
//  1. GTASKS_CLIENT_ID / GTASKS_CLIENT_SECRET env vars or config file (via koanf)
//  2. Client ID / secret embedded at build time via -ldflags
func GetOAuth2Config() (*oauth2.Config, error) {
	cfgClientID, cfgClientSecret := GetCredentials()

	clientID := cfgClientID
	if clientID == "" {
		clientID = ClientID
	}
	if clientID == "" {
		return nil, fmt.Errorf("no client ID found. Set GTASKS_CLIENT_ID env var, add credentials.client_id to config file, or rebuild with client ID")
	}

	clientSecret := cfgClientSecret
	if clientSecret == "" {
		clientSecret = ClientSecret
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("no client secret found. Set GTASKS_CLIENT_SECRET env var, add credentials.client_secret to config file, or rebuild with client secret")
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
