# ponysay-go

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](LICENCE)
[![Platform Support](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows%20%7C%20FreeBSD-brightgreen)](#installation)
[![Release](https://img.shields.io/github/v/release/bradly0cjw/ponysay-go?color=orange)](https://github.com/bradly0cjw/ponysay-go/releases)

Cross-platform Go port of [ponysay](https://github.com/erkin/ponysay) (cowsay reimplementation for ponies).

> [!NOTE]
> Personal project created strictly for fun. Maintainer is not responsible for ongoing support.

---

## Features

- **Zero External Dependencies**: Single self-contained binary with embedded assets (`//go:embed`). No Python runtime or `coreutils`/`stty` commands required.
- **Blazing Fast**: Sub-5ms startup time for instantaneous shell startup hooks (`fortune | ponysay`).
- **True Cross-Platform**: Native binaries for Linux, macOS, FreeBSD, and Windows (Cmd, PowerShell, Windows Terminal).
- **Hybrid Asset Engine**: Seamless resolution order combining embedded fallback assets with local disk overrides (`~/.config/ponysay/ponies`).
- **Smart Resolution**: Phonetic-weighted Levenshtein fuzzy search (`SpelloCorrecter`), Unicode accent remapping (`PONYSAY_UCS_ME`), terminal width auto-filtering, and `best.pony` fallback.
- **Full Feature Parity**: Speech (`ponysay`), thought (`ponythink`), quote database (`-q`), custom balloon borders (`-b`), word wrapping (`-W`), metadata info (`-i`), and built-in self-updater (`ponysay update`).

---

## Screenshots

| macOS | Linux | Windows |
| :---: | :---: | :---: |
| ![macOS](docs/img/macos-1.png) | ![Linux](docs/img/linux-1.png) | ![Windows](docs/img/windows-1.png) |

<details>
<summary>More Screenshots (Quote Mode)</summary>

| macOS (Quote) | Linux (Quote) | Windows (Quote) |
| :---: | :---: | :---: |
| ![macOS Quote](docs/img/macos-2.png) | ![Linux Quote](docs/img/linux-2.png) | ![Windows Quote](docs/img/windows-2.png) |

</details>

---

## Installation

### Quick Install (Recommended)

**Linux & macOS (Bash / Zsh):**
```bash
curl -fsSL https://raw.githubusercontent.com/bradly0cjw/ponysay-go/mane/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/bradly0cjw/ponysay-go/mane/install.ps1 | iex
```

*Detects OS and architecture automatically. Installs globally with admin/root privileges, or to user directory (`~/.local/bin` / `%LocalAppData%\Programs\ponysay`) without admin/root.*

#### Terminal Greeting (Optional)

Add `--terminal` to automatically hook `ponysay -q` into your shell startup, so a random pony quote greets you on every new terminal session:

**Linux & macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/bradly0cjw/ponysay-go/mane/install.sh | bash -s -- --terminal
```

**Windows (PowerShell):**
```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/bradly0cjw/ponysay-go/mane/install.ps1))) -Terminal
```

*Appends to `~/.bashrc`, `~/.zshrc`, `~/.config/fish/config.fish`, or your PowerShell `$PROFILE` depending on your shell. Idempotent — safe to re-run without duplicating the hook. The uninstall script automatically cleans it up.*

<details>
<summary>Manual & Alternative Installation Options</summary>

#### Prebuilt Release Binary
Download the latest binary for your platform from [GitHub Releases](https://github.com/bradly0cjw/ponysay-go/releases):
- **Linux / macOS**: Place binary in `/usr/local/bin` (global) or `~/.local/bin` (user). Symlink `ponythink` -> `ponysay`.
- **Windows**: Place `ponysay.exe` in Program Files or `%LocalAppData%\Programs\ponysay`, copy as `ponythink.exe`, and add folder to `PATH`.

#### Go CLI
```bash
go install github.com/bradly0cjw/ponysay-go/cmd/ponysay@latest
```

#### Build from Source
```bash
git clone https://github.com/bradly0cjw/ponysay-go.git
cd ponysay-go
go build -o ponysay ./cmd/ponysay
ln -s ponysay ponythink
```
</details>

### Uninstallation

<details>
<summary>Uninstall Instructions</summary>

**Linux & macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/bradly0cjw/ponysay-go/mane/uninstall.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/bradly0cjw/ponysay-go/mane/uninstall.ps1 | iex
```
</details>

---

## Quick Usage

```bash
# Speech & Thought Balloons
ponysay "I am just the cutest pony!"
ponythink "Hmm... is Golang fast?"

# Pony Selection & Quotes
ponysay -f pinkie "Partay!~"       # Select specific pony
ponysay +f cow "Moo!"              # Select extra/non-MLP pony
ponysay -q pinkie                  # Print quote from Pinkie Pie
ponysay -q                         # Print random pony quote

# Shell Startup Hook (put inside ~/.bashrc or ~/.zshrc)
fortune | ponysay                 # Pipe fortune output into ponysay for random pony messages on shell startup
ponysay -q                        # Print random pony quote on shell startup

# Balloon Styles & Formatting
ponysay -b unicode "Box border"   # Styles: cowsay, unicode, ascii, round, etc.
ponysay -f derpy -o               # Artwork only (no balloon)
ponysay -f derpy -i               # Print pony metadata info
ponysay -f derpy +i               # Print metadata with color highlights

# Listing Options
ponysay -l                        # List MLP ponies (-A for all, +l for extra)
ponysay -B                        # List balloon styles
ponysay --quoters                 # List ponies that have quotes
ponysay -l --onelist              # Print pony list in a single line

# Self Update
ponysay update                    # Download and install latest release from GitHub
```

> [!TIP]
> For an exhaustive command line options reference, environment variable configurations, and advanced features, read the full [ponysay-go Manual](docs/MANUAL.md).

---

## Custom Assets

Assets are resolved using a hybrid lookup order:
1. **Local Disk**: `./ponies/`, `~/.config/ponysay/ponies/`, `~/.ponysay/ponies/`, and `/usr/share/ponysay/ponies/`
2. **Embedded Fallback**: Embedded binary assets compiled into `ponysay-go`.

To add a custom pony file without recompiling:
```bash
mkdir -p ~/.config/ponysay/ponies
cp mycustompony.pony ~/.config/ponysay/ponies/
ponysay -f mycustompony "Hello world!"
```

---

## Fuzzy Search & Variant Resolution

`ponysay-go` includes full support for pony name resolution, spell correction, and environment overrides:

- **Fuzzy Spell Correction (`SpelloCorrecter`)**: If a pony name has a typo (e.g. `ponysay -f fluter-shy`), `ponysay-go` calculates weighted Levenshtein edit distance with phonetic weights (`k`↔`c`, `s`↔`z`, `o`↔`u`, etc.). Typo distances $\le 5$ (configurable via `PONYSAY_TYPO_LIMIT`) automatically match the closest pony.
- **Unicode Remapping (`PONYSAY_UCS_ME`)**: Set `PONYSAY_UCS_ME=1` to map Unicode accented pony names (e.g. `mjölna`, `bifröst`, `piñacolada`) to their ASCII filenames.
- **Terminal Width Filtering**: When selecting a random pony, `ponysay-go` automatically filters out ponies whose artwork width exceeds your current terminal size.
- **`best.pony` Fallback**: Automatically uses `best.pony` in asset directories when no pony selection flag is provided (if present).

---

## CLI Options Overview

| Flag / Command | Short / Alias | Description |
| :--- | :--- | :--- |
| `--file PONY` | `-f`, `+f`, `-F` | Select MLP (`-f`), extra (`+f`), or any (`-F`) pony |
| `--files PONY...` | `--f`, `++f`, `--F` | Variadic selection among multiple ponies |
| `--quote [PONY]` | `-q`, `+q`, `--q`, `--quotes` | Select pony quote (specific pony or random) |
| `--bubble STYLE` | `-b`, `--balloon` | Select balloon style (`cowsay`, `unicode`, `ascii`, `round`, etc.) |
| `--wrap COLUMN` | `-W COLUMN` | Specify maximum wrapping column width |
| `--compress` | `-c`, `--compact` | Compress empty lines in message text |
| `--list`, `--all` | `-l`, `+l`, `-A` | List MLP (`-l`), extra (`+l`), or all (`-A`) ponies |
| `--symlist`, `--bubblelist` | `-L`, `+L`, `-B`, `--quoters` | List MLP aliases (`-L`), extra aliases (`+L`), styles (`-B`), or quoters |
| `--onelist` | `++onelist`, `--Onelist` | Print listing output formatted on a single line |
| `--info`, `--pony-only` | `-i`, `+i`, `-o` | Standard info (`-i`), color info (`+i`), artwork only (`-o`) |
| `--256-colours` | `-X`, `-V`, `-K` | Color modes: 256 (`-X`), TTY 16 (`-V`), KMS (`-K`) |
| `update` | `-u`, `--update` | Update binary to latest release from GitHub |
| `--version`, `--help` | `-v`, `-h` | Display version (`-v`) or help menu (`-h`) |

---

## Python vs Go Comparison

| Feature | Original Python `ponysay` | Go Port (`ponysay-go`) |
| :--- | :--- | :--- |
| **Platform Support** | Linux / macOS (Windows via WSL/Cygwin) | Native Windows (`.exe`), macOS, Linux, FreeBSD |
| **Dependencies** | Python 3, `coreutils` (`stty`), `setup.py` | Single self-contained binary (zero dependencies) |
| **Startup Speed** | ~100ms+ (Python interpreter overhead) | Sub-5ms native execution |
| **Terminal Detection**| Shell call to `stty size` | Native OS syscalls (`golang.org/x/term`) |
| **Asset Engine** | Filesystem lookup | Embedded fallback (`//go:embed`) + Local overrides |
| **Self Updater** | Manual git pull / package manager | Built-in `ponysay update` command |
| **Feature Parity** | Speech, thought, quotes, styles, wrapping | 100% complete feature parity |

---

## Development

```bash
# Clone repository
git clone https://github.com/bradly0cjw/ponysay-go.git
cd ponysay-go

# Run test suite
go test -v ./...

# Build local binary
go build -o ponysay ./cmd/ponysay
```

---

## Credits & License

Reimplementation of the original Python **[ponysay](https://github.com/erkin/ponysay)** created by [Erkin Batu Altunbaş](https://github.com/erkin), [Mattias Andrée](https://github.com/maandree), [Pablo Lezaeta](https://github.com/jristz), and community contributors. Artwork sourced from [Browser Ponies](https://panzi.github.io/Browser-Ponies/) and Desktop Ponies.

See [CREDITS](CREDITS) for asset attributions and [LICENCE](LICENCE) for license details.
