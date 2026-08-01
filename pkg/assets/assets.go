package assets

import (
	"embed"
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//go:embed all:assets
var embeddedFS embed.FS

// AssetManager provides access to ponies, balloons, and quotes.
type AssetManager struct {
	rnd *rand.Rand
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetPonyFile finds and reads the content of a pony file by name.
func (am *AssetManager) GetPonyFile(name string, allowsNonMLP bool) (string, string, error) {
	if name == "" {
		return am.GetRandomPonyFile(allowsNonMLP)
	}

	cleanName := strings.TrimSuffix(name, ".pony")

	dirs := []string{"assets/ponies", "assets/ttyponies"}
	if allowsNonMLP {
		dirs = append([]string{"assets/extraponies", "assets/extrattyponies"}, dirs...)
	} else {
		dirs = append(dirs, "assets/extraponies", "assets/extrattyponies")
	}

	for _, dir := range dirs {
		path := filepath.Join(dir, cleanName+".pony")
		data, err := embeddedFS.ReadFile(path)
		if err == nil {
			return cleanName, string(data), nil
		}
	}

	allPonies := am.ListPonies(allowsNonMLP, true)
	for _, p := range allPonies {
		if strings.EqualFold(p, cleanName) {
			return am.GetPonyFile(p, allowsNonMLP)
		}
	}

	return "", "", fmt.Errorf("pony '%s' not found", name)
}

// GetRandomPonyFile picks a random pony.
func (am *AssetManager) GetRandomPonyFile(allowsNonMLP bool) (string, string, error) {
	ponies := am.ListPonies(allowsNonMLP, false)
	if len(ponies) == 0 {
		return "", "", fmt.Errorf("no ponies available")
	}
	chosen := ponies[am.rnd.Intn(len(ponies))]
	return am.GetPonyFile(chosen, allowsNonMLP)
}

// ListPonies returns sorted list of available pony names.
func (am *AssetManager) ListPonies(allowsNonMLP bool, includeExtra bool) []string {
	seen := make(map[string]bool)

	dirs := []string{"assets/ponies"}
	if allowsNonMLP || includeExtra {
		dirs = append(dirs, "assets/extraponies")
	}

	for _, dir := range dirs {
		entries, err := embeddedFS.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pony") {
				name := strings.TrimSuffix(entry.Name(), ".pony")
				seen[name] = true
			}
		}
	}

	var result []string
	for k := range seen {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

// GetBalloonContent returns the contents of a balloon style file.
func (am *AssetManager) GetBalloonContent(name string, isThink bool) (string, error) {
	ext := ".say"
	if isThink {
		ext = ".think"
	}

	if name == "" {
		name = "unicode"
	}
	cleanName := strings.TrimSuffix(name, ext)

	path := filepath.Join("assets/balloons", cleanName+ext)
	data, err := embeddedFS.ReadFile(path)
	if err == nil {
		return string(data), nil
	}

	data, err = embeddedFS.ReadFile(filepath.Join("assets/balloons", "unicode"+ext))
	if err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("balloon style '%s' not found", name)
}

// ListBalloons returns available balloon styles.
func (am *AssetManager) ListBalloons(isThink bool) []string {
	ext := ".say"
	if isThink {
		ext = ".think"
	}
	entries, err := embeddedFS.ReadDir("assets/balloons")
	if err != nil {
		return nil
	}
	var list []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
			list = append(list, strings.TrimSuffix(e.Name(), ext))
		}
	}
	sort.Strings(list)
	return list
}

// GetPonyQuote returns a random quote for a given pony.
func (am *AssetManager) GetPonyQuote(ponyName string) (string, string, error) {
	entries, err := embeddedFS.ReadDir("assets/ponyquotes")
	if err != nil || len(entries) == 0 {
		return "", "", fmt.Errorf("no quotes available")
	}

	var candidates []string
	if ponyName != "" {
		target := strings.ToLower(strings.TrimSuffix(ponyName, ".quote"))
		for _, e := range entries {
			if strings.EqualFold(strings.TrimSuffix(e.Name(), ".quote"), target) {
				candidates = append(candidates, e.Name())
			}
		}
	}

	if len(candidates) == 0 {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".quote") {
				candidates = append(candidates, e.Name())
			}
		}
	}

	if len(candidates) == 0 {
		return "", "", fmt.Errorf("no quotes found for %s", ponyName)
	}

	chosenFile := candidates[am.rnd.Intn(len(candidates))]
	data, err := embeddedFS.ReadFile(filepath.Join("assets/ponyquotes", chosenFile))
	if err != nil {
		return "", "", err
	}

	quotes := strings.Split(string(data), "\n%\n")
	var validQuotes []string
	for _, q := range quotes {
		trimmed := strings.TrimSpace(q)
		if trimmed != "" {
			validQuotes = append(validQuotes, trimmed)
		}
	}

	if len(validQuotes) == 0 {
		return "", "", fmt.Errorf("empty quote file %s", chosenFile)
	}

	quote := validQuotes[am.rnd.Intn(len(validQuotes))]
	pName := strings.TrimSuffix(chosenFile, ".quote")
	return pName, quote, nil
}
