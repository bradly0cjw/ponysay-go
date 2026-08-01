package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTargetAssetName(t *testing.T) {
	tests := []struct {
		goos     string
		goarch   string
		expected string
	}{
		{"darwin", "amd64", "ponysay-darwin-amd64"},
		{"darwin", "arm64", "ponysay-darwin-arm64"},
		{"linux", "amd64", "ponysay-linux-amd64"},
		{"linux", "arm64", "ponysay-linux-arm64"},
		{"windows", "amd64", "ponysay-windows-amd64.exe"},
		{"windows", "arm64", "ponysay-windows-arm64.exe"},
	}

	for _, tt := range tests {
		got := GetTargetAssetName(tt.goos, tt.goarch)
		if got != tt.expected {
			t.Errorf("GetTargetAssetName(%q, %q) = %q; want %q", tt.goos, tt.goarch, got, tt.expected)
		}
	}
}

func TestFetchLatestReleaseMock(t *testing.T) {
	mockRelease := Release{
		TagName: "v0.0.1",
		Name:    "v0.0.1",
		Assets: []ReleaseAsset{
			{Name: "ponysay-darwin-arm64", BrowserDownloadURL: "https://example.com/ponysay-darwin-arm64", Size: 1000},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockRelease)
	}))
	defer ts.Close()

	// Direct URL fetch override test if needed or testing JSON parsing directly
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("failed to connect to mock server: %v", err)
	}
	defer resp.Body.Close()

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		t.Fatalf("failed to decode release: %v", err)
	}

	if rel.TagName != "v0.0.1" {
		t.Errorf("got TagName %q; want v0.0.1", rel.TagName)
	}
	if len(rel.Assets) != 1 || rel.Assets[0].Name != "ponysay-darwin-arm64" {
		t.Errorf("unexpected assets: %+v", rel.Assets)
	}
}
