# ponysay-go

Cross-platform Go port of [ponysay](https://github.com/erkin/ponysay) (cowsay reimplementation for ponies).

> **Disclaimer**: This is a personal project created strictly for fun. The maintainer is not responsible for maintaining, updating, or providing ongoing support for this repository.

---

## Features

- Zero External Dependencies: Single self-contained binary (`//go:embed` bundles all 470+ pony art files, balloon styles, and quotes). No Python environment or `coreutils`/`stty` required.
- Blazing Fast: Sub-5ms startup time for shell startup scripts (`fortune | ponysay`).
- True Cross-Platform: Native support for Windows (Command Prompt, PowerShell, Windows Terminal), macOS, Linux, and FreeBSD.
- Hybrid Asset Lookup: Built-in embedded assets work out-of-the-box, with optional runtime overrides from local user directories (`~/.config/ponysay/ponies`).
- Full Feature Parity: Supports speech balloons (`ponysay`), thought balloons (`ponythink`), quote databases (`-q`), custom balloon borders (`-b`), word wrapping (`-W`), and metadata info (`-i`).

---

## Installation

### Option 1: Quick Install Script (Recommended)

Run the automated installer script, which detects your OS and architecture automatically:

```bash
curl -fsSL https://raw.githubusercontent.com/bradly0cjw/ponysay-go/mane/install.sh | bash
```

- **With Root / Admin Privileges** (`sudo` or write access to `/usr/local/bin`): Installs `ponysay` and `ponythink` globally to `/usr/local/bin`.
- **Without Root Privileges** (Non-admin): Installs `ponysay` and `ponythink` locally to `~/.local/bin`.

---

### Option 2: Manual Installation via Prebuilt Release

Download the latest prebuilt binary for your platform from [GitHub Releases](https://github.com/bradly0cjw/ponysay-go/releases):

- **macOS**: `ponysay-darwin-amd64` / `ponysay-darwin-arm64`
- **Linux**: `ponysay-linux-amd64` / `ponysay-linux-arm64`
- **Windows**: `ponysay-windows-amd64.exe` / `ponysay-windows-arm64.exe`

#### A. Global System Installation (With Root / Admin Privileges)

```bash
# 1. Download binary (e.g. Linux amd64 or macOS arm64)
curl -L -o ponysay https://github.com/bradly0cjw/ponysay-go/releases/latest/download/ponysay-linux-amd64

# 2. Make executable and move to /usr/local/bin
chmod +x ponysay
sudo mv ponysay /usr/local/bin/ponysay

# 3. Create global symlink for ponythink
sudo ln -sf /usr/local/bin/ponysay /usr/local/bin/ponythink
```

#### B. User Local Installation (Without Root / Non-Admin)

If you don't have `root` or `sudo` access:

```bash
# 1. Create local binary directory
mkdir -p ~/.local/bin

# 2. Download binary (e.g. Linux amd64 or macOS arm64)
curl -L -o ~/.local/bin/ponysay https://github.com/bradly0cjw/ponysay-go/releases/latest/download/ponysay-linux-amd64
chmod +x ~/.local/bin/ponysay

# 3. Create symlink for ponythink
ln -sf ~/.local/bin/ponysay ~/.local/bin/ponythink

# 4. Ensure ~/.local/bin is in your PATH (add to ~/.bashrc or ~/.zshrc if needed)
export PATH="$HOME/.local/bin:$PATH"
```

#### C. Windows Installation

Download `ponysay-windows-amd64.exe`, rename it to `ponysay.exe`, and move it to a folder in your System `PATH` (such as `C:\Windows` or a custom tools directory). Create a copy named `ponythink.exe` for thought balloons.

---

### Option 3: Build from Source

```bash
git clone https://github.com/bradly0cjw/ponysay-go.git
cd ponysay-go
go build -o ponysay ./cmd/ponysay
ln -s ponysay ponythink
```

---

### Option 4: Install via Go CLI

```bash
go install github.com/bradly0cjw/ponysay-go/cmd/ponysay@latest
```

---

## Basic Usage & Command Examples

### Speech Balloon (ponysay)

```bash
ponysay "I am just the cutest pony!"
```

### Thought Balloon (ponythink)

```bash
ponythink "Hmm... is Golang fast?"
```

### Specific Pony Selection (-f)

```bash
ponysay -f pinkie "Partay!~"
ponysay -f derpy "I brought muffins!"
ponysay +f cow "Moo!"
```

### Pony Quotes (-q)

```bash
# Print a quote from Pinkie Pie
ponysay -q pinkie

# Print a random quote from any pony with show quotes
ponysay -q

# Randomly select a quote from Pinkie Pie or Rarity
ponysay -q pinkie -q rarity
```

### Shell Startup Hook (Pony Fortune)

Add this to your `~/.bashrc` or `~/.zshrc`:

```bash
fortune | ponysay
```

### Custom Balloon Styles (-b)

```bash
ponysay -b cowsay "Default Cowsay style"
ponysay -b unicode "Unicode box border style"
ponysay -b ascii "ASCII border style"
```

### Listing Ponies & Styles

```bash
ponysay -l           # List MLP ponies
ponysay +l          # List extra/non-MLP ponies
ponysay -A          # List all ponies
ponysay -L          # List ponies with alternative names
ponysay -B          # List balloon styles
ponysay --quoters   # List ponies that have quotes
ponysay -l --onelist # Single line list output
```

### Metadata & Pony Art Only

```bash
ponysay -f derpy -i  # Print metadata info of Derpy
ponysay -f derpy -o  # Print artwork only without speech balloon
```

### Updating ponysay

```bash
# Fetch and install the latest release directly from GitHub
ponysay update

# Alternatively using ponythink or flags
ponythink update
ponysay -u
ponysay --update
```

---

## Custom Assets & User Overrides (Hybrid Model)

`ponysay-go` uses a hybrid asset resolution engine:
1. **Local Disk Lookup**: Checks user directories first for custom artwork:
   - `./ponies/`, `./balloons/`, `./ponyquotes/` (Current directory)
   - `~/.config/ponysay/ponies/` (User config folder)
   - `/usr/share/ponysay/` (System path)
2. **Embedded Fallback**: If an asset is not found on disk, it falls back to the embedded binary assets.

### Adding Custom Ponies

To add your own custom `.pony` file without recompiling:

```bash
mkdir -p ~/.config/ponysay/ponies
cp mycustompony.pony ~/.config/ponysay/ponies/
ponysay -f mycustompony "Hello world!"
```

---

## Command Line Flag Reference

| Flag | Short | Description |
| :--- | :--- | :--- |
| `--file PONY` | `-f` | Select a pony by name or file |
| `--file PONY` | `+f` | Select a non-MLP pony |
| `--any-pony PONY` | `-F` | Select any pony (MLP or non-MLP) |
| `--files PONY...` | `--f` | Variadic selection among ponies |
| `--quote [PONY]` | `-q` | Select a pony quote |
| `--quotes [PONY]` | `--q` | Variadic quote selection |
| `--bubble STYLE` | `-b` | Select balloon style (`cowsay`, `unicode`, `ascii`, etc.) |
| `--wrap COLUMN` | `-W` | Specify maximum wrapping width |
| `--list` | `-l` | List pony names |
| `--symlist` | `-L` | List pony names with alternative names |
| `--all` | `-A` | List all pony names |
| `--bubblelist` | `-B` | List balloon styles |
| `--quoters` | | List ponies that have show quotes |
| `--onelist` | | Format list output in a single line |
| `--compress` | `-c` | Compress empty lines in message |
| `--info` | `-i` | Print pony metadata info |
| `++info` | `+i` | Print pony metadata info formatted with colors |
| `--pony-only` | `-o` | Print only the pony artwork |
| `--256-colours` | `-X` | Enable 256 color mode |
| `--tty-colours` | `-V` | Enable TTY 16 color mode |
| `--kms-colours` | `-K` | Enable KMS color mode |
| `--update` / `update` | `-u` | Update to the latest release from GitHub |
| `--version` | `-v` | Print version information |
| `--help` | `-h` | Display help menu |

---

## Compatibility Matrix

| Feature / Flag | Original Python Implementation | Go Port (`ponysay-go`) | Status |
| :--- | :--- | :--- | :--- |
| **Native Windows Support** | Requires WSL or Cygwin | Native `.exe` (cmd, PowerShell, Windows Terminal) | Supported |
| **Dependencies** | Python 3, `coreutils` (`stty`), `setup.py` | Single self-contained binary (`//go:embed`) | Fully Self-Contained |
| **Pony Selection** (`-f`, `+f`, `-F`) | Supported | Supported | 100% Compatible |
| **Variadic Selection** (`--f`, `++f`, `--F`) | Supported | Supported | 100% Compatible |
| **Show Pony Quotes** (`-q`, `--quote`, `--quotes`) | Supported | Supported (290+ pony quote files indexed) | 100% Compatible |
| **Quoters Listing** (`--quoters`) | Supported | Supported | 100% Compatible |
| **Pony Listing** (`-l`, `+l`, `-A`) | Supported | Supported | 100% Compatible |
| **Alias Listing** (`-L`, `+L`, `+A`) | Supported | Supported | 100% Compatible |
| **One-Line Listing** (`--onelist`, `++onelist`) | Supported | Supported | 100% Compatible |
| **Balloon Listing** (`-B`, `--bubblelist`) | Supported | Supported | 100% Compatible |
| **Thought Balloon** (`ponythink`) | Supported | Supported (binary alias / symlink) | 100% Compatible |
| **Balloon Styles** (`-b cowsay`, `unicode`, `ascii`) | Supported | Supported (Default: `cowsay`) | 100% Compatible |
| **Word Wrapping** (`-W`, `--wrap`) | Supported | Supported (`mattn/go-runewidth` ANSI-aware) | 100% Compatible |
| **Message Compression** (`-c`, `--compress`) | Supported | Supported | 100% Compatible |
| **Pony Only Output** (`-o`, `--pony-only`) | Supported | Supported | 100% Compatible |
| **Pony Metadata Info** (`-i`, `+i`, `--info`) | Supported | Supported | 100% Compatible |
| **Balloon & Link Colors** (`+c`, `--colour-*`) | Supported | Supported | 100% Compatible |
| **Color Modes** (`-X` 256, `-V` TTY 16, `-K` KMS) | Supported | Supported | 100% Compatible |
| **Terminal Width Detection** | Shell call `stty size` | Native OS syscalls via `golang.org/x/term` | Native & Cross-Platform |

---

## Running Tests

Execute all Go unit tests and integration tests:

```bash
go test -v ./...
```

---

## Disclaimer

This is a personal project created strictly for fun. The maintainer is not responsible for maintaining, updating, or fixing issues with this repository.
