package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/BRO3886/gtasks/internal/config"
	"github.com/BRO3886/gtasks/internal/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
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
	folderPath := config.GetInstallLocation()
	tokFile := folderPath + "/token.json"
	existingToken, err := tokenFromFile(tokFile)
	if err == nil {
		// Token file exists, check if it's still valid
		if isTokenValid(oauthConfig, existingToken) {
			return fmt.Errorf("already logged in (token is valid)")
		}
		// Token exists but is invalid/expired, remove it and proceed
		utils.Info("Existing token is expired or invalid, re-authenticating...\n")
		os.Remove(tokFile)
	}

	// Perform PKCE + localhost authentication
	token, port, err := authenticateWithPKCE(oauthConfig)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Save token
	saveToken(tokFile, token)
	utils.Info("✓ Authorization successful! Server was running on port %d\n", port)
	utils.Info("✓ Credentials saved to %s\n", tokFile)

	return nil
}

// isTokenValid checks if a token is still valid by making a test API call
func isTokenValid(oauthConfig *oauth2.Config, token *oauth2.Token) bool {
	// Create a client with the token
	client := oauthConfig.Client(context.Background(), token)

	// Try to create a Tasks service and make a simple API call
	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return false
	}

	// Try to list task lists (minimal API call to verify token)
	_, err = srv.Tasklists.List().MaxResults(1).Do()
	return err == nil
}

// Logout removes stored authentication token
func Logout() error {
	folderPath := config.GetInstallLocation()
	tokFile := folderPath + "/token.json"
	err := os.Remove(tokFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("not logged in")
		}
		return fmt.Errorf("failed to logout: %v", err)
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
	// Generate PKCE verifier and challenge
	verifier := oauth2.GenerateVerifier()

	// Try to bind to one of the pre-approved ports
	port, listener, err := findAvailablePort()
	if err != nil {
		return nil, 0, fmt.Errorf("unable to find available port: %v", err)
	}
	defer listener.Close()

	// Update config with the selected port
	configCopy := *config
	configCopy.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)

	// Generate authorization URL with PKCE
	state := "gtasks-auth-state"
	authURL := configCopy.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier))

	utils.Info("Opening browser for Google authentication...\n")
	utils.Info("If browser doesn't open, visit: %s\n", authURL)
	utils.Info("Starting local server on http://localhost:%d...\n", port)

	// Try to open browser
	if err := utils.OpenBrowser(authURL); err != nil {
		utils.Warn("Failed to open browser automatically: %v\n", err)
		utils.Warn("Please manually visit the URL above\n")
	}

	// Start callback server
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
		// Verify state parameter (CSRF protection)
		if r.URL.Query().Get("state") != state {
			errorChan <- fmt.Errorf("invalid state parameter")
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			errorChan <- fmt.Errorf("no authorization code received")
			return
		}

		// Exchange authorization code for token with PKCE verifier
		token, err := config.Exchange(context.Background(), code,
			oauth2.VerifierOption(verifier))
		if err != nil {
			errorChan <- fmt.Errorf("failed to exchange code for token: %v", err)
			return
		}

		tokenChan <- token

		// Return success page
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

	// Ensure server shuts down
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	// Wait for result with timeout
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
	folderPath := config.GetInstallLocation()
	tokFile := folderPath + "/token.json"

	token, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, fmt.Errorf("not authenticated. Run 'gtasks login' first")
	}

	return oauthConfig.Client(context.Background(), token), nil
}

// tokenFromFile retrieves a token from a local file
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

// saveToken saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		utils.ErrorP("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
