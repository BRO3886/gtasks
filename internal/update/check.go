package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	githubRepo   = "BRO3886/gtasks"
	cacheTTL     = 24 * time.Hour
	checkTimeout = 2 * time.Second
)

// Result holds the outcome of an update check.
type Result struct {
	Latest    string
	HasUpdate bool
}

// CacheEntry represents the cached update check state.
type CacheEntry struct {
	CheckedAt time.Time
	Latest    string
}

// cacheDir returns the path to the cache directory (~/.cache/gtasks/).
func cacheDir(homeDir string) string {
	return filepath.Join(homeDir, ".cache", "gtasks")
}

// CachePath returns the path to the update check cache file.
func CachePath(homeDir string) string {
	return filepath.Join(cacheDir(homeDir), "update-check")
}

// ReadCache reads the cached update check result.
// Returns nil if the cache doesn't exist or is malformed.
func ReadCache(homeDir string) *CacheEntry {
	data, err := os.ReadFile(CachePath(homeDir))
	if err != nil {
		return nil
	}
	return ParseCache(string(data))
}

// ParseCache parses cache file contents into a CacheEntry.
// Returns nil if the content is malformed.
func ParseCache(content string) *CacheEntry {
	entry := &CacheEntry{}
	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		switch key {
		case "checked_at":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return nil
			}
			entry.CheckedAt = t
		case "latest":
			entry.Latest = val
		}
	}
	if entry.Latest == "" {
		return nil
	}
	return entry
}

// WriteCache writes the update check result to the cache file.
// Silently returns nil on error (cache is best-effort).
func WriteCache(homeDir string, entry *CacheEntry) error {
	dir := cacheDir(homeDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil // silently fail
	}
	content := fmt.Sprintf("checked_at=%s\nlatest=%s\n",
		entry.CheckedAt.UTC().Format(time.RFC3339),
		entry.Latest,
	)
	if err := os.WriteFile(CachePath(homeDir), []byte(content), 0o644); err != nil {
		return nil // silently fail
	}
	return nil
}

// IsCacheFresh returns true if the cache entry is less than cacheTTL old.
func IsCacheFresh(entry *CacheEntry, now time.Time) bool {
	return now.Sub(entry.CheckedAt) < cacheTTL
}

// githubRelease is the subset of GitHub API response we need.
type githubRelease struct {
	TagName string `json:"tag_name"`
}

// FetchLatestVersion fetches the latest release tag from GitHub.
func FetchLatestVersion(ctx context.Context, client *http.Client) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse github response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("empty tag_name in github response")
	}

	return release.TagName, nil
}

// CompareVersions returns true if latest is newer than current.
// Both should be semver-like strings (with or without 'v' prefix).
func CompareVersions(current, latest string) bool {
	cur := parseVersion(current)
	lat := parseVersion(latest)
	if cur == nil || lat == nil {
		return false
	}
	for i := range 3 {
		if lat[i] > cur[i] {
			return true
		}
		if lat[i] < cur[i] {
			return false
		}
	}
	return false
}

// parseVersion parses "v1.2.3" or "1.2.3" into []int.
// Returns nil if parsing fails.
func parseVersion(v string) []int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return nil
	}
	nums := make([]int, 3)
	for i, p := range parts {
		if idx := strings.IndexByte(p, '-'); idx >= 0 {
			p = p[:idx]
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil
		}
		nums[i] = n
	}
	return nums
}

// Check performs the full update check flow:
// 1. Read cache — if fresh, use cached result
// 2. If stale/missing, fetch from GitHub with timeout
// 3. Write result to cache
// Returns nil if the check was skipped or failed (never errors to callers).
func Check(homeDir, currentVersion string) *Result {
	now := time.Now()

	if cached := ReadCache(homeDir); cached != nil && IsCacheFresh(cached, now) {
		if CompareVersions(currentVersion, cached.Latest) {
			return &Result{Latest: cached.Latest, HasUpdate: true}
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()

	latest, err := FetchLatestVersion(ctx, http.DefaultClient)
	if err != nil {
		return nil // silently fail
	}

	WriteCache(homeDir, &CacheEntry{
		CheckedAt: now,
		Latest:    latest,
	})

	if CompareVersions(currentVersion, latest) {
		return &Result{Latest: latest, HasUpdate: true}
	}

	return nil
}
