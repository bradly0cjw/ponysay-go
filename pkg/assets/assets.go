package assets

import (
	"embed"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

//go:embed all:assets
var embeddedFS embed.FS

// AssetManager provides access to ponies, balloons, and quotes.
// Supports a hybrid model: user local directories take precedence,
// with embedded binary assets as default fallbacks.
type AssetManager struct {
	rnd            *rand.Rand
	mu             sync.Mutex
	aliasToQuote   map[string]string
	quoteFiles     map[string][]string
	ponyAliases    map[string]map[string]bool
	customDirs     []string
	customPonies   map[string]string // key: "subDir/cleanName", val: diskPath
	customBalloons map[string]string // key: "cleanName+ext", val: diskPath
	initialized    bool
}

func NewAssetManager() *AssetManager {
	am := &AssetManager{
		rnd:            rand.New(rand.NewSource(time.Now().UnixNano())),
		aliasToQuote:   make(map[string]string),
		quoteFiles:     make(map[string][]string),
		ponyAliases:    make(map[string]map[string]bool),
		customDirs:     getSearchDirectories(),
		customPonies:   make(map[string]string),
		customBalloons: make(map[string]string),
	}
	am.initCustomIndex()
	am.initQuotes()
	return am
}

func getSearchDirectories() []string {
	var candidateDirs []string
	if cwd, err := os.Getwd(); err == nil {
		candidateDirs = append(candidateDirs, cwd)
	}

	if userConfig, err := os.UserConfigDir(); err == nil {
		candidateDirs = append(candidateDirs, filepath.Join(userConfig, "ponysay"))
	}

	if home, err := os.UserHomeDir(); err == nil {
		candidateDirs = append(candidateDirs, filepath.Join(home, ".ponysay"))
	}

	candidateDirs = append(candidateDirs, "/usr/share/ponysay", "/usr/local/share/ponysay")

	var validDirs []string
	subdirs := []string{"ponies", "extraponies", "ttyponies", "extrattyponies", "balloons", "ponyquotes"}

	for _, dir := range candidateDirs {
		st, err := os.Stat(dir)
		if err != nil || !st.IsDir() {
			continue
		}
		hasAssets := false
		for _, sub := range subdirs {
			if subSt, subErr := os.Stat(filepath.Join(dir, sub)); subErr == nil && subSt.IsDir() {
				hasAssets = true
				break
			}
		}
		if hasAssets {
			validDirs = append(validDirs, dir)
		}
	}

	return validDirs
}

func (am *AssetManager) initCustomIndex() {
	if len(am.customDirs) == 0 {
		return
	}
	ponySubdirs := []string{"ponies", "extraponies", "ttyponies", "extrattyponies"}
	for _, baseDir := range am.customDirs {
		for _, subDir := range ponySubdirs {
			diskDir := filepath.Join(baseDir, subDir)
			entries, err := os.ReadDir(diskDir)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pony") {
					cleanName := strings.TrimSuffix(entry.Name(), ".pony")
					key := subDir + "/" + cleanName
					if _, exists := am.customPonies[key]; !exists {
						am.customPonies[key] = filepath.Join(diskDir, entry.Name())
					}
				}
			}
		}

		bDir := filepath.Join(baseDir, "balloons")
		bEntries, err := os.ReadDir(bDir)
		if err == nil {
			for _, entry := range bEntries {
				if !entry.IsDir() {
					if _, exists := am.customBalloons[entry.Name()]; !exists {
						am.customBalloons[entry.Name()] = filepath.Join(bDir, entry.Name())
					}
				}
			}
		}
	}
}

func (am *AssetManager) initQuotes() {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.initialized {
		return
	}
	am.initialized = true

	// 1. Read alias mapping from embedded assets
	if poniesMapData, err := embeddedFS.ReadFile("assets/ponyquotes/ponies"); err == nil {
		am.parsePoniesAliasFile(string(poniesMapData))
	}

	// 2. Read custom ponies mapping if present on local FS
	for _, baseDir := range am.customDirs {
		pPath := filepath.Join(baseDir, "ponyquotes", "ponies")
		if data, err := os.ReadFile(pPath); err == nil {
			am.parsePoniesAliasFile(string(data))
		}
	}

	// 3. Index embedded quote files
	if entries, err := embeddedFS.ReadDir("assets/ponyquotes"); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || entry.Name() == "ponies" {
				continue
			}
			name := entry.Name()
			if dotIdx := strings.Index(name, "."); dotIdx > 0 {
				basePony := name[:dotIdx]
				relPath := path.Join("assets/ponyquotes", name)
				am.quoteFiles[basePony] = append(am.quoteFiles[basePony], "embed:"+relPath)
				if _, exists := am.aliasToQuote[basePony]; !exists {
					am.aliasToQuote[basePony] = basePony
				}
			}
		}
	}

	// 4. Index local FS custom quote files
	for _, baseDir := range am.customDirs {
		qDir := filepath.Join(baseDir, "ponyquotes")
		if entries, err := os.ReadDir(qDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() || entry.Name() == "ponies" {
					continue
				}
				name := entry.Name()
				if dotIdx := strings.Index(name, "."); dotIdx > 0 {
					basePony := name[:dotIdx]
					absPath := filepath.Join(qDir, name)
					am.quoteFiles[basePony] = append(am.quoteFiles[basePony], "file:"+absPath)
				}
			}
		}
	}
}

func (am *AssetManager) parsePoniesAliasFile(content string) {
	lines := strings.Split(content, "\n")
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

// GetPonyFile finds and reads the content of a pony file by name.
// First checks local user folders, then falls back to embedded binary assets.
func (am *AssetManager) GetPonyFile(name string, allowsNonMLP bool) (string, string, error) {
	if name == "" {
		return am.GetRandomPonyFile(allowsNonMLP)
	}

	cleanName := strings.TrimSuffix(name, ".pony")

	dirs := []string{"ponies", "ttyponies"}
	if allowsNonMLP {
		dirs = append([]string{"extraponies", "extrattyponies"}, dirs...)
	} else {
		dirs = append(dirs, "extraponies", "extrattyponies")
	}

	// 1. Check custom ponies index (if any custom directories exist)
	if len(am.customPonies) > 0 {
		for _, subDir := range dirs {
			key := subDir + "/" + cleanName
			if diskPath, ok := am.customPonies[key]; ok {
				data, err := os.ReadFile(diskPath)
				if err == nil {
					return cleanName, string(data), nil
				}
			}
		}
	}

	// 2. Fall back to embedded binary assets
	for _, subDir := range dirs {
		embedPath := path.Join("assets", subDir, cleanName+".pony")
		data, err := embeddedFS.ReadFile(embedPath)
		if err == nil {
			return cleanName, string(data), nil
		}
	}

	// 3. Try case-insensitive alias match
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

// ListPonies returns sorted list of available pony names (combining local FS & embedded).
func (am *AssetManager) ListPonies(allowsNonMLP bool, includeExtra bool) []string {
	seen := make(map[string]bool)

	dirs := []string{"ponies"}
	if allowsNonMLP || includeExtra {
		dirs = append(dirs, "extraponies")
	}

	// Custom FS entries (from in-memory customPonies index)
	for key := range am.customPonies {
		parts := strings.SplitN(key, "/", 2)
		if len(parts) == 2 {
			subDir, cleanName := parts[0], parts[1]
			for _, d := range dirs {
				if d == subDir {
					seen[cleanName] = true
				}
			}
		}
	}

	// Embedded FS directories
	for _, subDir := range dirs {
		embedDir := path.Join("assets", subDir)
		entries, err := embeddedFS.ReadDir(embedDir)
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
		name = "cowsay"
	}
	cleanName := strings.TrimSuffix(name, ext)
	fileName := cleanName + ext

	// 1. Check custom balloons index
	if diskPath, ok := am.customBalloons[fileName]; ok {
		if data, err := os.ReadFile(diskPath); err == nil {
			return string(data), nil
		}
	}

	// 2. Check embedded assets
	embedPath := path.Join("assets/balloons", fileName)
	data, err := embeddedFS.ReadFile(embedPath)
	if err == nil {
		return string(data), nil
	}

	data, err = embeddedFS.ReadFile(path.Join("assets/balloons", "cowsay"+ext))
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
	seen := make(map[string]bool)

	// Custom balloons
	for fileName := range am.customBalloons {
		if strings.HasSuffix(fileName, ext) {
			seen[strings.TrimSuffix(fileName, ext)] = true
		}
	}

	// Embedded balloons
	entries, err := embeddedFS.ReadDir("assets/balloons")
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
				seen[strings.TrimSuffix(e.Name(), ext)] = true
			}
		}
	}

	var list []string
	for k := range seen {
		list = append(list, k)
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

// GetPonyQuote selects a quote and corresponding pony name.
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

	chosenSpec := files[am.rnd.Intn(len(files))]
	var data []byte
	var err error

	if strings.HasPrefix(chosenSpec, "embed:") {
		data, err = embeddedFS.ReadFile(strings.TrimPrefix(chosenSpec, "embed:"))
	} else if strings.HasPrefix(chosenSpec, "file:") {
		data, err = os.ReadFile(strings.TrimPrefix(chosenSpec, "file:"))
	}

	if err != nil {
		return targetPony, "Zecora! Help me, I am mute!", nil
	}

	quoteText := strings.TrimSpace(string(data))
	if quoteText == "" {
		quoteText = "Zecora! Help me, I am mute!"
	}

	return targetPony, quoteText, nil
}

