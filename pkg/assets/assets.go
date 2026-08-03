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

//go:embed balloons extraponies extrattyponies ponies ponyquotes ttyponies share
var embeddedFS embed.FS

// AssetManager provides access to ponies, balloons, and quotes.
// Supports a hybrid model: user local directories take precedence,
// with embedded binary assets as default fallbacks.
type AssetManager struct {
	mu             sync.RWMutex
	rnd            *rand.Rand
	aliasToQuote   map[string]string
	quoteFiles     map[string][]string
	ponyAliases    map[string]map[string]bool
	casePonyMap    map[string]string // lowerCase -> canonical clean pony name
	ucsMap         map[string]string
	reverseUCSMap  map[string]string
	ponyWidths     map[string]int
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
		casePonyMap:    make(map[string]string),
		ucsMap:         make(map[string]string),
		reverseUCSMap:  make(map[string]string),
		ponyWidths:     make(map[string]int),
		customDirs:     getSearchDirectories(),
		customPonies:   make(map[string]string),
		customBalloons: make(map[string]string),
	}
	am.initCustomIndex()
	am.initQuotes()
	am.initCasePonyMap()
	am.initUCSMap()
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
	// 0. Parse ponyquotes/ponies alias mappings
	if data, err := embeddedFS.ReadFile("ponyquotes/ponies"); err == nil {
		am.parsePoniesAliasFile(string(data))
	}
	for _, baseDir := range am.customDirs {
		aliasPath := filepath.Join(baseDir, "ponyquotes", "ponies")
		if data, err := os.ReadFile(aliasPath); err == nil {
			am.parsePoniesAliasFile(string(data))
		}
	}

	// 1. Index quote files in ponies/ (legacy/embedded if any)
	if entries, err := embeddedFS.ReadDir("ponies"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".quote") {
				name := entry.Name()
				basePony := strings.TrimSuffix(name, ".quote")
				relPath := path.Join("ponies", name)
				for _, p := range strings.Split(basePony, "+") {
					if p != "" {
						am.quoteFiles[p] = append(am.quoteFiles[p], "embed:"+relPath)
						if _, exists := am.aliasToQuote[p]; !exists {
							am.aliasToQuote[p] = p
						}
					}
				}
			}
		}
	}

	// 2. Index quote files in custom ponies dirs
	for _, baseDir := range am.customDirs {
		pDir := filepath.Join(baseDir, "ponies")
		if entries, err := os.ReadDir(pDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".quote") {
					name := entry.Name()
					basePony := strings.TrimSuffix(name, ".quote")
					absPath := filepath.Join(pDir, name)
					for _, p := range strings.Split(basePony, "+") {
						if p != "" {
							am.quoteFiles[p] = append(am.quoteFiles[p], "file:"+absPath)
						}
					}
				}
			}
		}
	}

	// 3. Index embedded quote files
	if entries, err := embeddedFS.ReadDir("ponyquotes"); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || entry.Name() == "ponies" {
				continue
			}
			name := entry.Name()
			if dotIdx := strings.Index(name, "."); dotIdx > 0 {
				basePony := name[:dotIdx]
				relPath := path.Join("ponyquotes", name)
				for _, p := range strings.Split(basePony, "+") {
					if p != "" {
						am.quoteFiles[p] = append(am.quoteFiles[p], "embed:"+relPath)
						if _, exists := am.aliasToQuote[p]; !exists {
							am.aliasToQuote[p] = p
						}
					}
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
					for _, p := range strings.Split(basePony, "+") {
						if p != "" {
							am.quoteFiles[p] = append(am.quoteFiles[p], "file:"+absPath)
						}
					}
				}
			}
		}
	}

	// 5. Index MASTER metadata tags from .pony files to set quote relationships
	ponySubdirs := []string{"ponies", "extraponies", "ttyponies", "extrattyponies"}
	for _, subDir := range ponySubdirs {
		if entries, err := embeddedFS.ReadDir(subDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pony") {
					cleanName := strings.ToLower(strings.TrimSuffix(entry.Name(), ".pony"))
					relPath := path.Join(subDir, entry.Name())
					if data, err := embeddedFS.ReadFile(relPath); err == nil {
						if master := parseMasterTag(string(data)); master != "" {
							am.aliasToQuote[cleanName] = master
						}
					}
				}
			}
		}
	}

	for _, baseDir := range am.customDirs {
		for _, subDir := range ponySubdirs {
			diskDir := filepath.Join(baseDir, subDir)
			if entries, err := os.ReadDir(diskDir); err == nil {
				for _, entry := range entries {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pony") {
						cleanName := strings.ToLower(strings.TrimSuffix(entry.Name(), ".pony"))
						absPath := filepath.Join(diskDir, entry.Name())
						if data, err := os.ReadFile(absPath); err == nil {
							if master := parseMasterTag(string(data)); master != "" {
								am.aliasToQuote[cleanName] = master
							}
						}
					}
				}
			}
		}
	}
}

func parseMasterTag(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "$$$" {
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "$$$" {
				break
			}
			if strings.HasPrefix(line, "MASTER:") {
				val := strings.TrimSpace(strings.TrimPrefix(line, "MASTER:"))
				return strings.ToLower(val)
			}
		}
	}
	return ""
}

func (am *AssetManager) initCasePonyMap() {
	allPonies := am.listPoniesLocked(true, true)
	for _, p := range allPonies {
		am.casePonyMap[strings.ToLower(p)] = p
	}
}

func (am *AssetManager) parsePoniesAliasFile(content string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimRight(line, "\r"))
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
// GetPonyFile finds and reads the content of a pony file by name.
// includeStandard: search standard pony directories ("ponies", "ttyponies")
// includeExtra: search extra pony directories ("extraponies", "extrattyponies")
func (am *AssetManager) GetPonyFile(name string, includeStandard bool, includeExtra bool) (string, string, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.getPonyFileLocked(name, includeStandard, includeExtra)
}

func (am *AssetManager) getPonyFileLocked(name string, includeStandard bool, includeExtra bool) (string, string, error) {
	if name == "" {
		return am.getRandomPonyFileLocked(includeStandard, includeExtra)
	}

	name = am.RemapUCS(name)
	cleanName := strings.TrimSuffix(name, ".pony")

	var dirs []string
	var excludedDirs []string

	if includeStandard {
		dirs = append(dirs, "ponies", "ttyponies")
	} else {
		excludedDirs = append(excludedDirs, "ponies", "ttyponies")
	}

	if includeExtra {
		dirs = append(dirs, "extraponies", "extrattyponies")
	} else {
		excludedDirs = append(excludedDirs, "extraponies", "extrattyponies")
	}

	// 1. Check custom ponies index in allowed dirs
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

	// 2. Fall back to embedded binary assets in allowed dirs
	for _, subDir := range dirs {
		embedPath := path.Join(subDir, cleanName+".pony")
		data, err := embeddedFS.ReadFile(embedPath)
		if err == nil {
			return cleanName, string(data), nil
		}
	}

	// 3. Check if exact match exists in EXCLUDED dirs
	if len(am.customPonies) > 0 {
		for _, subDir := range excludedDirs {
			key := subDir + "/" + cleanName
			if _, ok := am.customPonies[key]; ok {
				return "", "", fmt.Errorf("I have never heard of anypony named %s", name)
			}
		}
	}
	for _, subDir := range excludedDirs {
		embedPath := path.Join(subDir, cleanName+".pony")
		if _, err := embeddedFS.ReadFile(embedPath); err == nil {
			return "", "", fmt.Errorf("I have never heard of anypony named %s", name)
		}
	}

	// 4. Try case-insensitive lookup within allowed ponies
	availablePonies := am.listPoniesLocked(includeStandard, includeExtra)
	lowerClean := strings.ToLower(cleanName)
	for _, p := range availablePonies {
		if strings.ToLower(p) == lowerClean {
			return am.getPonyFileLocked(p, includeStandard, includeExtra)
		}
	}

	// 5. Try fuzzy spell correction against allowed ponies
	corrector := NewSpelloCorrecter()
	bestMatches, _ := corrector.Correct(cleanName, availablePonies)
	if len(bestMatches) > 0 {
		chosen := bestMatches[am.rnd.Intn(len(bestMatches))]
		return am.getPonyFileLocked(chosen, includeStandard, includeExtra)
	}

	return "", "", fmt.Errorf("I have never heard of anypony named %s", name)
}

// GetRandomPonyFile picks a random pony.
func (am *AssetManager) GetRandomPonyFile(includeStandard bool, includeExtra bool) (string, string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	return am.getRandomPonyFileLocked(includeStandard, includeExtra)
}

func (am *AssetManager) getRandomPonyFileLocked(includeStandard bool, includeExtra bool) (string, string, error) {
	var checkDirs []string
	if includeStandard {
		checkDirs = append(checkDirs, "ponies")
	}
	if includeExtra {
		checkDirs = append(checkDirs, "extraponies")
	}

	// 1. Check for best.pony fallback
	for _, subDir := range checkDirs {
		if diskPath, ok := am.customPonies[subDir+"/best"]; ok {
			if data, err := os.ReadFile(diskPath); err == nil {
				return "best", string(data), nil
			}
		}
		if data, err := embeddedFS.ReadFile(subDir + "/best.pony"); err == nil {
			return "best", string(data), nil
		}
	}

	ponies := am.listPoniesLocked(includeStandard, includeExtra)
	if len(ponies) == 0 {
		return "", "", fmt.Errorf("no ponies available")
	}

	fittingPonies := am.FilterFittingPonies(ponies, includeStandard, includeExtra)
	chosen := fittingPonies[am.rnd.Intn(len(fittingPonies))]
	return am.getPonyFileLocked(chosen, includeStandard, includeExtra)
}

// ListPonies returns sorted list of available pony names (combining local FS & embedded).
func (am *AssetManager) ListPonies(includeStandard bool, includeExtra bool) []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.listPoniesLocked(includeStandard, includeExtra)
}

func (am *AssetManager) listPoniesLocked(includeStandard bool, includeExtra bool) []string {
	seen := make(map[string]bool)

	var dirs []string
	if includeStandard {
		dirs = append(dirs, "ponies")
	}
	if includeExtra {
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
		entries, err := embeddedFS.ReadDir(subDir)
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

// ListPoniesFormatted returns pony names formatted with bold ANSI codes for ponies that have quotes.
func (am *AssetManager) ListPoniesFormatted(includeStandard bool, includeExtra bool) []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	basePonies := am.listPoniesLocked(includeStandard, includeExtra)
	quoters := make(map[string]bool)
	for q := range am.quoteFiles {
		quoters[q] = true
	}

	var result []string
	for _, p := range basePonies {
		if quoters[p] {
			result = append(result, "\x1b[1m"+p+"\x1b[0m")
		} else {
			result = append(result, p)
		}
	}
	return result
}

var symlinkMap = map[string][]string{
	// ponies/
	"airheart":       {"jetstream"},
	"blinkie":        {"limestone"},
	"blueberry":      {"berrydreams"},
	"blues":          {"noteworthy"},
	"bonbon":         {"sweetiedrops"},
	"boxxy":          {"craftycrate"},
	"caballeron":     {"drcaballeron"},
	"carrot":         {"carrottop", "goldenharvest"},
	"carrotcake":     {"carecake"},
	"clyde":          {"igneousrock"},
	"coldheart":      {"snowheart"},
	"colgate":        {"minuette"},
	"flashsentry":    {"brad"},
	"fleurdelis":     {"fleurdislee"},
	"grace":          {"manewitz"},
	"hairytipper":    {"dancefever"},
	"highscore":      {"buttonmash"},
	"horsemd":        {"doctop"},
	"hughjelly":      {"hughbertjellius"},
	"inky":           {"marble"},
	"lilyvalley":     {"lily"},
	"lotus":          {"lotusblossom"},
	"lovemelody":     {"venus"},
	"lyra":           {"harpass", "heartstrings"},
	"lyrabonbon":     {"bonbonlyra"},
	"manticore":      {"mannyroar"},
	"maybelle":       {"mabel"},
	"misspommel":     {"cocopommel"},
	"mrsparkle":      {"nightlight"},
	"mrssparkle":     {"twilightvelvet"},
	"oinkoinkoink":   {"pinkieoink"},
	"perrypierce":    {"perry"},
	"pokeypierce":    {"royalpin"},
	"powderrouge":    {"sindy"},
	"prettyvision":   {"elsie"},
	"quickfix":       {"clockwork", "epona"},
	"raindrops":      {"sunshowerraindrops"},
	"rara":           {"countess"},
	"raritysdad":     {"hondoflanks", "magnum"},
	"raritysmom":     {"bettybouffant", "cookiecrumbles", "pearl"},
	"ravenunicorn":   {"raven"},
	"rose":           {"roseluck"},
	"ruby":           {"berrypinch"},
	"snowflake":      {"bulkbiceps", "horsepower"},
	"sparkler":       {"amethyststar"},
	"stormyflare":    {"spitfiresmom"},
	"sue":            {"cloudyquartz"},
	"timeturner":     {"drhooves"},
	"trixie":         {"lulamoon", "trixielulamoon"},
	"vinyl":          {"djpon-3"},
	"violet":         {"royalribbon"},
	"waltercoltchak": {"walter"},

	// extraponies/
	"barbara":          {"barbra"},
	"internetexplorer": {"ie"},
}

// ListPoniesWithAliases returns pony names formatted with alternative names (aliases).
func (am *AssetManager) ListPoniesWithAliases(includeStandard bool, includeExtra bool) []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	basePonies := am.listPoniesLocked(includeStandard, includeExtra)
	quoters := make(map[string]bool)
	for q := range am.quoteFiles {
		quoters[q] = true
	}

	var result []string
	for _, p := range basePonies {
		displayPony := p
		if quoters[p] {
			displayPony = "\x1b[1m" + p + "\x1b[0m"
		}

		if syms, ok := symlinkMap[p]; ok && len(syms) > 0 {
			var altList []string
			for _, alt := range syms {
				if quoters[alt] {
					altList = append(altList, "\x1b[1m"+alt+"\x1b[0m")
				} else {
					altList = append(altList, alt)
				}
			}
			sort.Strings(altList)
			if len(altList) > 0 {
				result = append(result, fmt.Sprintf("%s (%s)", displayPony, strings.Join(altList, " ")))
			} else {
				result = append(result, displayPony)
			}
		} else {
			result = append(result, displayPony)
		}
	}
	return result
}

// PonyGroup represents a collection of ponies originating from a specific directory source.
type PonyGroup struct {
	DirectoryPath string
	Ponies        []string
}

func (am *AssetManager) applyUCSInList(groupPonies []string, dynamicSymlinks map[string][]string) []string {
	envVal := strings.ToLower(os.Getenv("PONYSAY_UCS_ME"))
	ucsConf := 0
	if envVal == "yes" || envVal == "y" || envVal == "1" {
		ucsConf = 1
	} else if envVal == "harder" || envVal == "h" || envVal == "2" {
		ucsConf = 2
	}

	if ucsConf == 0 {
		return groupPonies
	}

	if ucsConf == 1 {
		var additions []string
		for _, p := range groupPonies {
			if ucs, ok := am.reverseUCSMap[p]; ok {
				additions = append(additions, ucs)
				if dynamicSymlinks != nil {
					dynamicSymlinks[p] = append(dynamicSymlinks[p], ucs)
				}
			}
		}
		groupPonies = append(groupPonies, additions...)
	} else if ucsConf == 2 {
		for i, p := range groupPonies {
			if ucs, ok := am.reverseUCSMap[p]; ok {
				groupPonies[i] = ucs
			}
		}
	}

	return groupPonies
}

// GetPonyGroups returns ponies grouped by their source directory (with directory paths for headers).
func (am *AssetManager) GetPonyGroups(includeStandard bool, includeExtra bool, withAliases bool, formatted bool) []PonyGroup {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var targetSubdirs []string
	if includeStandard {
		targetSubdirs = append(targetSubdirs, "ponies")
	}
	if includeExtra {
		targetSubdirs = append(targetSubdirs, "extraponies")
	}

	quoters := make(map[string]bool)
	for q := range am.quoteFiles {
		quoters[q] = true
	}

	var groups []PonyGroup

	for _, subDir := range targetSubdirs {
		seenInGroup := make(map[string]bool)
		var displayDirPath string
		var groupPonies []string
		dynamicSymlinks := make(map[string][]string)

		// 1. Custom FS search
		for _, baseDir := range am.customDirs {
			dPath := filepath.Join(baseDir, subDir)
			if st, err := os.Stat(dPath); err == nil && st.IsDir() {
				entries, err := os.ReadDir(dPath)
				if err == nil {
					for _, entry := range entries {
						if strings.HasSuffix(entry.Name(), ".pony") {
							name := strings.TrimSuffix(entry.Name(), ".pony")
							if !seenInGroup[name] {
								seenInGroup[name] = true
								groupPonies = append(groupPonies, name)
							}

							fullPath := filepath.Join(dPath, entry.Name())
							if lst, err := os.Lstat(fullPath); err == nil && (lst.Mode()&os.ModeSymlink != 0) {
								if target, err := os.Readlink(fullPath); err == nil {
									targetName := strings.TrimSuffix(filepath.Base(target), ".pony")
									if targetName != "" && targetName != name {
										dynamicSymlinks[targetName] = append(dynamicSymlinks[targetName], name)
									}
								}
							}
						}
					}
				}
				if len(groupPonies) > 0 && displayDirPath == "" {
					displayDirPath = filepath.Clean(dPath) + "/"
				}
			}
		}

		// 2. Embedded FS fallback
		if len(groupPonies) == 0 {
			entries, err := embeddedFS.ReadDir(subDir)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pony") {
						name := strings.TrimSuffix(entry.Name(), ".pony")
						if !seenInGroup[name] {
							seenInGroup[name] = true
							groupPonies = append(groupPonies, name)
						}
					}
				}
			}
			if displayDirPath == "" {
				displayDirPath = "bundled:" + subDir + "/"
			}
		}

		if len(groupPonies) == 0 {
			continue
		}

		groupPonies = am.applyUCSInList(groupPonies, dynamicSymlinks)

		var finalItems []string

		if !withAliases {
			sort.Strings(groupPonies)
			for _, p := range groupPonies {
				displayPony := p
				if formatted && quoters[p] {
					displayPony = "\x1b[1m" + p + "\x1b[0m"
				}
				finalItems = append(finalItems, displayPony)
			}
		} else {
			aliasMap := make(map[string][]string)
			allAliases := make(map[string]bool)

			for target, syms := range dynamicSymlinks {
				for _, sym := range syms {
					aliasMap[target] = append(aliasMap[target], sym)
					allAliases[sym] = true
				}
			}

			for _, p := range groupPonies {
				if syms, ok := symlinkMap[p]; ok {
					for _, sym := range syms {
						aliasMap[p] = append(aliasMap[p], sym)
						allAliases[sym] = true
					}
				}
			}

			var topPonies []string
			for _, p := range groupPonies {
				if !allAliases[p] {
					topPonies = append(topPonies, p)
				}
			}
			sort.Strings(topPonies)

			for _, p := range topPonies {
				displayPony := p
				if formatted && quoters[p] {
					displayPony = "\x1b[1m" + p + "\x1b[0m"
				}

				if syms, ok := aliasMap[p]; ok && len(syms) > 0 {
					symSet := make(map[string]bool)
					var altList []string
					for _, alt := range syms {
						if !symSet[alt] {
							symSet[alt] = true
							if formatted && quoters[alt] {
								altList = append(altList, "\x1b[1m"+alt+"\x1b[0m")
							} else {
								altList = append(altList, alt)
							}
						}
					}
					sort.Strings(altList)
					if len(altList) > 0 {
						displayPony = fmt.Sprintf("%s (%s)", displayPony, strings.Join(altList, " "))
					}
				}
				finalItems = append(finalItems, displayPony)
			}
		}

		groups = append(groups, PonyGroup{
			DirectoryPath: displayDirPath,
			Ponies:        finalItems,
		})
	}

	return groups
}

// ListPoniesOneList returns a sorted list of unique pony names without formatting or grouping.
func (am *AssetManager) ListPoniesOneList(includeStandard bool, includeExtra bool) []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var targetSubdirs []string
	if includeStandard {
		targetSubdirs = append(targetSubdirs, "ponies")
	}
	if includeExtra {
		targetSubdirs = append(targetSubdirs, "extraponies")
	}

	seen := make(map[string]bool)
	var ponies []string

	for _, subDir := range targetSubdirs {
		for _, baseDir := range am.customDirs {
			dPath := filepath.Join(baseDir, subDir)
			if entries, err := os.ReadDir(dPath); err == nil {
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".pony") {
						name := strings.TrimSuffix(entry.Name(), ".pony")
						if !seen[name] {
							seen[name] = true
							ponies = append(ponies, name)
						}
					}
				}
			}
		}

		if entries, err := embeddedFS.ReadDir(subDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pony") {
					name := strings.TrimSuffix(entry.Name(), ".pony")
					if !seen[name] {
						seen[name] = true
						ponies = append(ponies, name)
					}
				}
			}
		}
	}

	ponies = am.applyUCSInList(ponies, nil)

	sort.Strings(ponies)
	var unique []string
	last := ""
	for _, p := range ponies {
		if p != last {
			unique = append(unique, p)
			last = p
		}
	}

	return unique
}

// ListQuoters returns sorted list of all ponies that have quotes.
func (am *AssetManager) ListQuoters(includeStandard bool, includeExtra bool) []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	return am.listQuotersLocked(includeStandard, includeExtra)
}

func (am *AssetManager) listQuotersLocked(includeStandard bool, includeExtra bool) []string {
	var targetSubdirs []string
	if includeStandard {
		targetSubdirs = append(targetSubdirs, "ponies")
	}
	if includeExtra {
		targetSubdirs = append(targetSubdirs, "extraponies")
	}

	availablePonies := make(map[string]bool)
	for _, subDir := range targetSubdirs {
		for _, baseDir := range am.customDirs {
			dPath := filepath.Join(baseDir, subDir)
			if entries, err := os.ReadDir(dPath); err == nil {
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".pony") {
						availablePonies[strings.TrimSuffix(entry.Name(), ".pony")] = true
					}
				}
			}
		}
		if entries, err := embeddedFS.ReadDir(subDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pony") {
					availablePonies[strings.TrimSuffix(entry.Name(), ".pony")] = true
				}
			}
		}
	}

	var quoters []string
	for ponyName, files := range am.quoteFiles {
		if len(files) > 0 && availablePonies[ponyName] {
			quoters = append(quoters, ponyName)
		}
	}
	sort.Strings(quoters)
	var unique []string
	last := ""
	for _, q := range quoters {
		if q != last {
			unique = append(unique, q)
			last = q
		}
	}
	return unique
}

// GetBalloonContent returns the contents of a balloon style file.
func (am *AssetManager) GetBalloonContent(name string, isThink bool) (string, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

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
	embedPath := path.Join("balloons", fileName)
	data, err := embeddedFS.ReadFile(embedPath)
	if err == nil {
		return string(data), nil
	}

	data, err = embeddedFS.ReadFile(path.Join("balloons", "cowsay"+ext))
	if err == nil {
		return string(data), nil
	}

	availableBalloons := am.ListBalloons(isThink)
	corrector := NewSpelloCorrecter()
	bestMatches, _ := corrector.Correct(cleanName, availableBalloons)
	if len(bestMatches) > 0 {
		chosen := bestMatches[am.rnd.Intn(len(bestMatches))]
		return am.GetBalloonContent(chosen, isThink)
	}

	return "", fmt.Errorf("That balloon style %s does not exist", name)
}

// ListBalloons returns available balloon styles.
func (am *AssetManager) ListBalloons(isThink bool) []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

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
	entries, err := embeddedFS.ReadDir("balloons")
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

func (am *AssetManager) resolveQuotePonyLocked(name string) (string, string) {
	if name == "" {
		return "", ""
	}
	name = am.RemapUCS(name)
	cleanName := strings.ToLower(strings.TrimSuffix(name, ".pony"))

	hasQuotes := func(quoteKey string) bool {
		files, ok := am.quoteFiles[quoteKey]
		return ok && len(files) > 0
	}

	// 1. Direct match in aliasToQuote or quoteFiles
	if base, ok := am.aliasToQuote[cleanName]; ok && hasQuotes(base) {
		targetPony := cleanName
		if matched, ok := am.casePonyMap[cleanName]; ok {
			targetPony = matched
		}
		return targetPony, base
	}
	if hasQuotes(cleanName) {
		targetPony := cleanName
		if matched, ok := am.casePonyMap[cleanName]; ok {
			targetPony = matched
		}
		return targetPony, cleanName
	}

	// 2. Case-insensitive check
	if matchedName, ok := am.casePonyMap[cleanName]; ok {
		cleanMatched := strings.ToLower(matchedName)
		if base, ok := am.aliasToQuote[cleanMatched]; ok && hasQuotes(base) {
			return matchedName, base
		}
		if hasQuotes(cleanMatched) {
			return matchedName, cleanMatched
		}
	}

	// 3. Fuzzy search against all available ponies first (to preserve requested variant name)
	availablePonies := am.listPoniesLocked(true, true)
	if len(availablePonies) > 0 {
		corrector := NewSpelloCorrecter()
		bestMatches, _ := corrector.Correct(cleanName, availablePonies)
		if len(bestMatches) > 0 {
			chosenPony := bestMatches[am.rnd.Intn(len(bestMatches))]
			cleanChosen := strings.ToLower(chosenPony)
			if base, ok := am.aliasToQuote[cleanChosen]; ok && hasQuotes(base) {
				return chosenPony, base
			}
			if hasQuotes(cleanChosen) {
				return chosenPony, cleanChosen
			}
		}
	}

	// 4. Fuzzy search against all available quoters
	quoters := am.listQuotersLocked(true, true)
	if len(quoters) > 0 {
		corrector := NewSpelloCorrecter()
		bestMatches, _ := corrector.Correct(cleanName, quoters)
		if len(bestMatches) > 0 {
			chosenQuoter := bestMatches[am.rnd.Intn(len(bestMatches))]
			if base, ok := am.aliasToQuote[chosenQuoter]; ok && hasQuotes(base) {
				return chosenQuoter, base
			}
			return chosenQuoter, chosenQuoter
		}
	}

	return cleanName, cleanName
}

// GetPonyQuote selects a quote and corresponding pony name.
func (am *AssetManager) GetPonyQuote(choices []string) (string, string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	var targetPony string
	var baseQuoteKey string

	if len(choices) > 0 {
		chosenChoice := choices[am.rnd.Intn(len(choices))]
		targetPony, baseQuoteKey = am.resolveQuotePonyLocked(chosenChoice)
	} else {
		quoters := am.listQuotersLocked(true, true)
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
