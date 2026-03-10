package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/BRO3886/gtasks/internal/config"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

const (
	keyringService = "gtasks"
	keyringUser    = "oauth2-token"
)

// Pre-approved ports that are registered in Google Cloud Console
var approvedPorts = []int{8080, 8081, 8082, 9090, 9091}

// Login performs OAuth2 authentication using PKCE + localhost flow
func Login() error {
	// Get OAuth2 configuration
	oauthConfig, err := config.GetOAuth2Config()
	if err != nil {
		return fmt.Errorf("failed to get OAuth2 config: %v", err)
	}

	// Validate configuration
	if err := config.ValidateOAuth2Config(oauthConfig); err != nil {
		return fmt.Errorf("invalid OAuth2 config: %v", err)
	}

	// Check if already logged in with a valid token
	existingToken, err := loadToken()
	if err == nil {
		if isTokenValid(oauthConfig, existingToken) {
			return fmt.Errorf("already logged in (token is valid)")
		}
		// Token exists but is invalid/expired — remove it and proceed
		utils.Info("Existing token is expired or invalid, re-authenticating...\n")
		if err := deleteToken(); err != nil {
			utils.Warn("Failed to clear existing token: %v\n", err)
		}
	}

	// Perform PKCE + localhost authentication
	token, port, err := authenticateWithPKCE(oauthConfig)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Save token
	backend, err := saveToken(token)
	if err != nil {
		return fmt.Errorf("failed to save credentials: %v", err)
	}
	utils.Info("✓ Authorization successful! Server was running on port %d\n", port)
	utils.Info("✓ Credentials saved to %s\n", backend)

	return nil
}

// isTokenValid checks if a token is still valid by making a test API call
func isTokenValid(oauthConfig *oauth2.Config, token *oauth2.Token) bool {
	client := oauthConfig.Client(context.Background(), token)

	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return false
	}

	_, err = srv.Tasklists.List().MaxResults(1).Do()
	return err == nil
}

// Logout removes stored authentication token
func Logout() error {
	if err := deleteToken(); err != nil {
		return err
	}
	utils.Info("✓ Successfully logged out\n")
	return nil
}

// GetService creates a Google Tasks service client
func GetService() (*tasks.Service, error) {
	oauthConfig, err := config.GetOAuth2Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth2 config: %v", err)
	}

	client, err := getClient(oauthConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %v", err)
	}

	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Tasks service: %v", err)
	}

	return srv, nil
}

// authenticateWithPKCE performs OAuth2 authentication with PKCE
func authenticateWithPKCE(config *oauth2.Config) (*oauth2.Token, int, error) {
	verifier := oauth2.GenerateVerifier()

	port, listener, err := findAvailablePort()
	if err != nil {
		return nil, 0, fmt.Errorf("unable to find available port: %v", err)
	}
	defer listener.Close()

	configCopy := *config
	configCopy.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)

	state := "gtasks-auth-state"
	authURL := configCopy.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier))

	utils.Info("Opening browser for Google authentication...\n")
	utils.Info("If browser doesn't open, visit: %s\n", authURL)
	utils.Info("Starting local server on http://localhost:%d...\n", port)

	if err := utils.OpenBrowser(authURL); err != nil {
		utils.Warn("Failed to open browser automatically: %v\n", err)
		utils.Warn("Please manually visit the URL above\n")
	}

	return startCallbackServer(listener, &configCopy, verifier, state, port)
}

// findAvailablePort tries to bind to one of the pre-approved ports
func findAvailablePort() (int, net.Listener, error) {
	for _, port := range approvedPorts {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			return port, listener, nil
		}
	}
	return 0, nil, fmt.Errorf("all approved ports (%v) are occupied", approvedPorts)
}

// startCallbackServer handles the OAuth2 callback
func startCallbackServer(listener net.Listener, config *oauth2.Config, verifier, state string, port int) (*oauth2.Token, int, error) {
	tokenChan := make(chan *oauth2.Token, 1)
	errorChan := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			errorChan <- fmt.Errorf("invalid state parameter")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errorChan <- fmt.Errorf("no authorization code received")
			return
		}

		token, err := config.Exchange(context.Background(), code,
			oauth2.VerifierOption(verifier))
		if err != nil {
			errorChan <- fmt.Errorf("failed to exchange code for token: %v", err)
			return
		}

		tokenChan <- token

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>GTasks Authentication</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .success { color: #28a745; }
        .container { max-width: 500px; margin: 0 auto; }
    </style>
</head>
<body>
    <div class="container">
        <h2 class="success">Authorization Successful!</h2>
        <p>You can close this browser window and return to your terminal.</p>
        <p><small>GTasks CLI is now authenticated and ready to use.</small></p>
    </div>
    <script>
        // Auto-close after 3 seconds
        setTimeout(function(){ window.close(); }, 3000);
    </script>
</body>
</html>`)
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errorChan <- fmt.Errorf("server error: %v", err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	select {
	case token := <-tokenChan:
		return token, port, nil
	case err := <-errorChan:
		return nil, port, err
	case <-time.After(5 * time.Minute):
		return nil, port, fmt.Errorf("authentication timeout (5 minutes)")
	}
}

// getClient retrieves HTTP client with valid token
func getClient(oauthConfig *oauth2.Config) (*http.Client, error) {
	token, err := loadToken()
	if err != nil {
		return nil, fmt.Errorf("not authenticated. Run 'gtasks login' first")
	}

	return oauthConfig.Client(context.Background(), token), nil
}

// saveToken serializes a token and stores it in the system keyring.
// Falls back to a plain file if the keyring is unavailable.
// Returns a human-readable description of where the token was stored.
func saveToken(token *oauth2.Token) (string, error) {
	data, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token: %v", err)
	}

	if err := keyring.Set(keyringService, keyringUser, string(data)); err != nil {
		// Keyring unavailable (e.g. headless server) — fall back to file
		utils.Warn("System keyring unavailable (%v), falling back to file storage\n", err)
		return "file", saveTokenToFile(token)
	}
	return "system keyring", nil
}

// loadToken retrieves the token from the keyring, falling back to a legacy file.
func loadToken() (*oauth2.Token, error) {
	data, err := keyring.Get(keyringService, keyringUser)
	if err == nil {
		var token oauth2.Token
		if jsonErr := json.Unmarshal([]byte(data), &token); jsonErr == nil {
			return &token, nil
		} else {
			// Corrupt keyring entry — warn and clean it up before falling through
			utils.Warn("Keyring entry is corrupt, clearing it: %v\n", jsonErr)
			keyring.Delete(keyringService, keyringUser)
		}
	}

	// Keyring unavailable or empty — try legacy file
	folderPath := config.GetInstallLocation()
	tokFile := folderPath + "/token.json"
	token, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}

	// Migrate to keyring — verify the write round-trips before deleting the file
	b, jsonErr := json.Marshal(token)
	if jsonErr != nil {
		return token, nil // can't migrate, but token is still usable
	}
	if migrateErr := keyring.Set(keyringService, keyringUser, string(b)); migrateErr == nil {
		// Verify the entry is readable before removing the file
		if verify, readErr := keyring.Get(keyringService, keyringUser); readErr == nil && verify == string(b) {
			os.Remove(tokFile)
			utils.Info("✓ Migrated credentials from file to system keyring\n")
		}
	}

	return token, nil
}

// deleteToken removes the token from both keyring and legacy file.
func deleteToken() error {
	keyringErr := keyring.Delete(keyringService, keyringUser)

	folderPath := config.GetInstallLocation()
	tokFile := folderPath + "/token.json"
	fileErr := os.Remove(tokFile)

	// Success if at least one token was actually deleted
	if keyringErr == nil || fileErr == nil {
		return nil
	}

	// Both failed — check if it's simply "not found" vs a real error
	if errors.Is(keyringErr, keyring.ErrNotFound) && os.IsNotExist(fileErr) {
		return fmt.Errorf("not logged in")
	}

	// At least one had a real I/O error
	if keyringErr != nil && !errors.Is(keyringErr, keyring.ErrNotFound) {
		return fmt.Errorf("failed to delete from keyring: %v", keyringErr)
	}
	return fmt.Errorf("failed to delete token file: %v", fileErr)
}

// tokenFromFile retrieves a token from a local file (legacy / fallback)
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	return &token, err
}

// saveTokenToFile is the fallback when keyring is unavailable
func saveTokenToFile(token *oauth2.Token) error {
	folderPath := config.GetInstallLocation()
	path := folderPath + "/token.json"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
