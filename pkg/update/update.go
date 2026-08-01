package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	DefaultRepo = "bradly0cjw/ponysay-go"
	UserAgent   = "ponysay-go-updater"
)

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type Release struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	PublishedAt time.Time      `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

// FetchLatestRelease retrieves release metadata from GitHub API.
func FetchLatestRelease(repo string) (*Release, error) {
	if repo == "" {
		repo = DefaultRepo
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error fetching latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d %s", resp.StatusCode, resp.Status)
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("failed to decode release JSON: %w", err)
	}
	return &rel, nil
}

// GetTargetAssetName returns the expected binary artifact filename for GOOS/GOARCH.
func GetTargetAssetName(goos, goarch string) string {
	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}
	return fmt.Sprintf("ponysay-%s-%s%s", goos, goarch, ext)
}

// SelfUpdate checks GitHub for the latest release and updates the current executable.
func SelfUpdate(currentVersion string, repo string) error {
	fmt.Println("Checking for updates from GitHub...")
	rel, err := FetchLatestRelease(repo)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	fmt.Printf("Current version : %s\n", currentVersion)
	fmt.Printf("Latest release  : %s (%s)\n", rel.TagName, rel.Name)

	targetAsset := GetTargetAssetName(runtime.GOOS, runtime.GOARCH)
	var downloadURL string
	var assetSize int64

	for _, asset := range rel.Assets {
		if asset.Name == targetAsset {
			downloadURL = asset.BrowserDownloadURL
			assetSize = asset.Size
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no prebuilt binary found for %s/%s (expected asset: %s)", runtime.GOOS, runtime.GOARCH, targetAsset)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to locate current executable path: %w", err)
	}
	evalPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		evalPath = execPath
	}

	fmt.Printf("Downloading %s (%.2f MB)...\n", targetAsset, float64(assetSize)/(1024*1024))

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download server returned status %d %s", resp.StatusCode, resp.Status)
	}

	dir := filepath.Dir(evalPath)
	tmpFile, err := os.CreateTemp(dir, "ponysay-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file in %s (do you need elevated permissions/sudo?): %w", dir, err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	_, err = io.Copy(tmpFile, resp.Body)
	_ = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to write updated binary: %w", err)
	}

	if err := os.Chmod(tmpName, 0755); err != nil {
		return fmt.Errorf("failed to set permissions on updated binary: %w", err)
	}

	if runtime.GOOS == "windows" {
		oldPath := evalPath + ".old"
		_ = os.Remove(oldPath)
		if err := os.Rename(evalPath, oldPath); err != nil {
			return fmt.Errorf("failed to rename existing binary on Windows: %w", err)
		}
		if err := os.Rename(tmpName, evalPath); err != nil {
			_ = os.Rename(oldPath, evalPath)
			return fmt.Errorf("failed to install new binary: %w", err)
		}
	} else {
		if err := os.Rename(tmpName, evalPath); err != nil {
			return fmt.Errorf("failed to replace executable at %s (do you need elevated permissions/sudo?): %w", evalPath, err)
		}
	}

	fmt.Printf("Successfully updated ponysay to %s!\n", rel.TagName)
	return nil
}
