# Agent Guide (`AGENTS.md`)

This guide provides context, architectural background, build/test commands, and guidelines for AI agents working on `ponysay-go`.

---

## 1. Project Overview

`ponysay-go` is a zero-dependency, cross-platform Go port of [ponysay](https://github.com/erkin/ponysay) (a cowsay reimplementation for ponies).

### Key Characteristics & Philosophy
- **Zero External Runtime Dependencies**: Compiles to a single self-contained binary with embedded assets (`//go:embed`). No Python runtime, `coreutils`, or `stty` commands required.
- **Sub-5ms Startup Time**: Optimized for instant shell startup hooks (`fortune | ponysay` or `ponysay -q` in `.bashrc` / `.zshrc`).
- **Cross-Platform**: Full native support for Linux, macOS, FreeBSD, and Windows (Cmd, PowerShell, Windows Terminal).
- **Dual Binary Identity**: Operates as `ponysay` (speech balloons) or `ponythink` (thought balloons) based on `os.Args[0]` (or symlink name).
- **100% Feature Parity with Python `ponysay`**: Implements speech/thought balloons, quotes database, balloon styles, word wrapping, metadata info, fuzzy spell correction, UCS accent remapping, terminal width auto-filtering, and a built-in self-updater (`ponysay update`).

---

## 2. Directory Structure & Architecture

```
ponysay-go/
├── cmd/
│   └── ponysay/
│       ├── main.go               # Main CLI entry point, flag parsing, command routing
│       └── main_test.go          # End-to-end CLI integration test suite
├── pkg/
│   ├── assets/
│   │   ├── assets.go             # Embedded & disk asset loading, hybrid resolution logic
│   │   ├── assets_test.go        # Unit tests for asset manager
│   │   ├── balloons/             # Embedded balloon style templates (.balloon files)
│   │   ├── dimension.go          # Pony artwork width calculation & terminal width filtering
│   │   ├── extraponies/          # Embedded extra pony artwork files (.pony)
│   │   ├── extrattyponies/       # Embedded TTY extra pony files (.ttypony)
│   │   ├── ponies/               # Embedded standard MLP pony artwork files (.pony)
│   │   ├── ponyquotes/           # Embedded quote files (.quote)
│   │   ├── spellocorrecter.go    # Phonetic Levenshtein fuzzy search for pony name typos
│   │   ├── spellocorrecter_test.go # Unit tests for spell correction
│   │   ├── ttyponies/            # Embedded TTY pony files (.ttypony)
│   │   ├── ucs.go                # Unicode accent remapping (PONYSAY_UCS_ME)
│   │   └── ucs_test.go           # Unit tests for UCS remapping
│   ├── balloon/
│   │   ├── balloon.go            # Balloon style parser, text word-wrapping, border layout
│   │   └── balloon_test.go       # Balloon parsing & text wrapping tests
│   ├── color/
│   │   ├── color.go              # ANSI escape sequence stripping & unicode display width
│   │   └── color_test.go         # ANSI & display width unit tests
│   ├── pony/
│   │   ├── pony.go               # Pony file parser, color rendering, balloon tail linking
│   │   └── pony_test.go          # Pony parser & rendering tests
│   ├── term/
│   │   └── term.go               # Native terminal width/height syscall detection
│   └── update/
│       ├── update.go             # Self-updater querying GitHub Releases API
│       └── update_test.go        # Unit tests for self-updater
├── docs/
│   └── MANUAL.md                 # Comprehensive user manual & CLI flag reference
├── go.mod                        # Go module configuration (Go 1.26+)
├── go.sum                        # Dependency checksums
├── install.sh / install.ps1      # Install scripts for Unix & Windows
├── uninstall.sh / uninstall.ps1  # Uninstall scripts for Unix & Windows
└── README.md                     # Main repository documentation
```

---

## 3. Core Architectural Concepts

### Asset Resolution Order
Assets (ponies, extra ponies, balloons, quotes) are resolved using a hybrid lookup model:
1. **Local Disk Directories**: Checks user/system directories first (`./ponies/`, `~/.config/ponysay/ponies/`, `~/.ponysay/ponies/`, `/usr/share/ponysay/ponies/`).
2. **Embedded Binary Assets**: Fallback to assets embedded in the binary via Go's `//go:embed` filesystem (`embed.FS`).

### Name Resolution & Fuzzy Matching
- **`SpelloCorrecter`**: Calculates weighted Levenshtein distance with phonetic substitution rules (e.g. `k`↔`c`, `s`↔`z`, `o`↔`u`) when matching pony names. Typo threshold defaults to edit distance $\le 5$ (configurable via `PONYSAY_TYPO_LIMIT`).
- **UCS Accent Remapping**: When `PONYSAY_UCS_ME=1` environment variable is set, accented pony names (e.g. `mjölna`, `bifröst`, `piñacolada`) are mapped to ASCII filenames.
- **Terminal Width Auto-Filtering**: When selecting random ponies, ponies exceeding current terminal width are automatically excluded.

### Speech vs. Thought Mode
- Executable binary inspects `filepath.Base(os.Args[0])`.
- If named `ponythink` (or `ponythink.exe`), default balloon type is set to `thought` (cloud bubble).
- If named `ponysay` (or `ponysay.exe`), default balloon type is set to `speech` (standard speech bubble).
- Command-line flags can explicitly override balloon types.

---

## 4. Development & Testing Instructions

### Run Unit & Integration Tests
Always run the full test suite before committing changes:
```bash
go test -v ./...
```

To run tests in a specific package:
```bash
go test -v ./pkg/assets
go test -v ./cmd/ponysay
```

### Build Binary
To build the executable locally:
```bash
go build -o ponysay ./cmd/ponysay
```

To test `ponythink` behavior locally via symlink:
```bash
ln -s ponysay ponythink
./ponythink "What am I thinking?"
```

### Verify Code Formatting
Use standard Go formatting tools:
```bash
go fmt ./...
go vet ./...
```

---

## 5. Coding Standards & Git Conventions

### Code Guidelines
1. **Idiomatic Go**: Follow standard Go naming conventions and error handling patterns. Return descriptive errors rather than panicking.
2. **Zero External Runtime Dependencies**: Keep minimal third-party dependencies. Prefer Go standard library (`os`, `io`, `embed`, `net/http`, `fmt`, etc.) or official `golang.org/x/` packages.
3. **Preserve CLI Flag Parity**: When modifying CLI flag parsing in [cmd/ponysay/main.go](cmd/ponysay/main.go), preserve backward compatibility with traditional Python `ponysay` flags (`-f`, `+f`, `-F`, `-q`, `-b`, `-W`, `-l`, `-i`, `+i`, `-o`, `-v`, `-h`, `update`).
4. **Thread Safety & Asset Caching**: The asset manager ([pkg/assets/assets.go](pkg/assets/assets.go)) uses `sync.Once` for embedded asset loading to guarantee thread safety during concurrent asset resolution. Maintain thread-safe patterns when extending asset management.
5. **Cross-Platform Compatibility**: Do not use hardcoded Unix paths or platform-specific syscalls outside of [pkg/term/term.go](pkg/term/term.go). Use `filepath.Join` and `os.UserConfigDir()` for file system paths.

### Git Commits & Versioning
1. **Conventional Commits**: All commit messages MUST strictly adhere to the [Conventional Commits](https://www.conventionalcommits.org/) specification using lower-case types and concise descriptions:
   - `feat: add feature X`
   - `fix: resolve issue Y`
   - `docs: update documentation`
   - `refactor: simplify asset resolution logic`
   - `test: add unit test for balloon parsing`
   - `chore: update dependencies`
2. **Semantic Versioning Tags**: Git release tags MUST strictly follow [Semantic Versioning (SemVer)](https://semver.org/) prefixed with `v` (e.g. `v1.0.0`, `v1.1.2`, `v2.0.0-rc.1`).

---

## 6. Key Source References

- Entry point & CLI flags: [cmd/ponysay/main.go](cmd/ponysay/main.go)
- CLI Integration tests: [cmd/ponysay/main_test.go](cmd/ponysay/main_test.go)
- Asset engine & fuzzy search: [pkg/assets/assets.go](pkg/assets/assets.go) & [pkg/assets/spellocorrecter.go](pkg/assets/spellocorrecter.go)
- Balloon parser & text wrapping: [pkg/balloon/balloon.go](pkg/balloon/balloon.go)
- Pony renderer & ANSI color parser: [pkg/pony/pony.go](pkg/pony/pony.go)
- Self-updater: [pkg/update/update.go](pkg/update/update.go)
- Full User Manual: [docs/MANUAL.md](docs/MANUAL.md)
