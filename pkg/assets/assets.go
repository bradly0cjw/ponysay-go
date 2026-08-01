package assets

import (
	"embed"
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

//go:embed all:assets
var embeddedFS embed.FS

// AssetManager provides access to ponies, balloons, and quotes.
type AssetManager struct {
	rnd          *rand.Rand
	mu           sync.Mutex
	aliasToQuote map[string]string            // e.g. "derpysit" -> "derpy"
	quoteFiles   map[string][]string          // e.g. "derpy" -> ["assets/ponyquotes/derpy.0", "assets/ponyquotes/derpy.1", ...]
	ponyAliases  map[string]map[string]bool   // e.g. "derpy" -> {"derpysit": true, "derpystand": true, ...}
	initialized  bool
}

func NewAssetManager() *AssetManager {
	am := &AssetManager{
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
		aliasToQuote: make(map[string]string),
		quoteFiles:   make(map[string][]string),
		ponyAliases:  make(map[string]map[string]bool),
	}
	am.initQuotes()
	return am
}

func (am *AssetManager) initQuotes() {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.initialized {
		return
	}
	am.initialized = true

	// Read alias mapping file `assets/ponyquotes/ponies`
	poniesMapData, err := embeddedFS.ReadFile("assets/ponyquotes/ponies")
	if err == nil {
		lines := strings.Split(string(poniesMapData), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, "+")
			if len(parts) == 0 {
				continue
			}
			baseName := parts[0]
			if _, ok := am.ponyAliases[baseName]; !ok {
				am.ponyAliases[baseName] = make(map[string]bool)
			}
			for _, alias := range parts {
				am.aliasToQuote[alias] = baseName
				am.ponyAliases[baseName][alias] = true
			}
		}
	}

	// Index all quote files `assets/ponyquotes/<pony>.<num>`
	entries, err := embeddedFS.ReadDir("assets/ponyquotes")
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || entry.Name() == "ponies" {
				continue
			}
			name := entry.Name()
			dotIdx := strings.Index(name, ".")
			if dotIdx > 0 {
				basePony := name[:dotIdx]
				relPath := filepath.Join("assets/ponyquotes", name)
				am.quoteFiles[basePony] = append(am.quoteFiles[basePony], relPath)
				if _, exists := am.aliasToQuote[basePony]; !exists {
					am.aliasToQuote[basePony] = basePony
				}
			}
		}
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

	// Try case-insensitive or partial alias match
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
	am.mu.Lock()
	defer am.mu.Unlock()

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

// ListPoniesWithAliases returns pony names formatted with alternative names (aliases).
func (am *AssetManager) ListPoniesWithAliases(allowsNonMLP bool, includeExtra bool) []string {
	basePonies := am.ListPonies(allowsNonMLP, includeExtra)
	var result []string

	for _, p := range basePonies {
		if aliasesMap, ok := am.ponyAliases[p]; ok && len(aliasesMap) > 1 {
			var altList []string
			for alt := range aliasesMap {
				if alt != p {
					altList = append(altList, alt)
				}
			}
			sort.Strings(altList)
			result = append(result, fmt.Sprintf("%s (%s)", p, strings.Join(altList, ", ")))
		} else {
			result = append(result, p)
		}
	}
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

// ListQuoters returns sorted list of all ponies that have quotes.
func (am *AssetManager) ListQuoters() []string {
	am.mu.Lock()
	defer am.mu.Unlock()

	var quoters []string
	for ponyName, files := range am.quoteFiles {
		if len(files) > 0 {
			quoters = append(quoters, ponyName)
		}
	}
	sort.Strings(quoters)
	return quoters
}

// GetPonyQuote selects a quote and corresponding pony name for given target pony choices.
// If choices is empty, picks a random pony from all quoters.
func (am *AssetManager) GetPonyQuote(choices []string) (string, string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	var targetPony string
	var baseQuoteKey string

	if len(choices) > 0 {
		targetPony = choices[am.rnd.Intn(len(choices))]
		cleanTarget := strings.ToLower(strings.TrimSuffix(targetPony, ".pony"))
		if base, ok := am.aliasToQuote[cleanTarget]; ok {
			baseQuoteKey = base
		} else {
			baseQuoteKey = cleanTarget
		}
	} else {
		// Pick random quoter
		quoters := make([]string, 0, len(am.quoteFiles))
		for q, files := range am.quoteFiles {
			if len(files) > 0 {
				quoters = append(quoters, q)
			}
		}
		if len(quoters) == 0 {
			return "derpy", "Zecora! Help me, I am mute!", nil
		}
		baseQuoteKey = quoters[am.rnd.Intn(len(quoters))]
		targetPony = baseQuoteKey
	}

	files, ok := am.quoteFiles[baseQuoteKey]
	if !ok || len(files) == 0 {
		return targetPony, "Zecora! Help me, I am mute!", nil
	}

	chosenFile := files[am.rnd.Intn(len(files))]
	data, err := embeddedFS.ReadFile(chosenFile)
	if err != nil {
		return targetPony, "Zecora! Help me, I am mute!", nil
	}

	quoteText := strings.TrimSpace(string(data))
	if quoteText == "" {
		quoteText = "Zecora! Help me, I am mute!"
	}

	return targetPony, quoteText, nil
}
