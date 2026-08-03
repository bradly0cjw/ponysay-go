package assets

import (
	"os"
	"path/filepath"
	"strings"
)

func (am *AssetManager) initUCSMap() {
	am.ucsMap = make(map[string]string)
	am.reverseUCSMap = make(map[string]string)

	parseMap := func(content string) {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(strings.TrimRight(line, "\r"))
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, "→")
			if len(parts) == 2 {
				ucs := strings.TrimSpace(parts[0])
				ascii := strings.TrimSpace(parts[1])
				am.ucsMap[ucs] = ascii
				am.reverseUCSMap[ascii] = ucs
			}
		}
	}

	// 1. Embedded ucsmap
	if data, err := embeddedFS.ReadFile("share/ucsmap"); err == nil {
		parseMap(string(data))
	}

	// 2. Custom disk ucsmap
	for _, baseDir := range am.customDirs {
		for _, uPath := range []string{filepath.Join(baseDir, "share", "ucsmap"), filepath.Join(baseDir, "ucsmap")} {
			if data, err := os.ReadFile(uPath); err == nil {
				parseMap(string(data))
			}
		}
	}
}

// RemapUCS converts Unicode pony names to ASCII names if PONYSAY_UCS_ME is set.
func (am *AssetManager) RemapUCS(name string) string {
	envVal := strings.ToLower(os.Getenv("PONYSAY_UCS_ME"))
	if envVal == "yes" || envVal == "y" || envVal == "1" || envVal == "harder" || envVal == "h" || envVal == "2" {
		am.mu.RLock()
		defer am.mu.RUnlock()
		if ascii, ok := am.ucsMap[name]; ok {
			return ascii
		}
	}
	return name
}
