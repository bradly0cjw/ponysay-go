# ponysay-go

High-performance, zero-dependency, cross-platform Go port of ponysay (cowsay reimplementation for ponies).

---

## Compatibility & Feature Support Matrix

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

## Installation

### Option 1: Build from Source

```bash
git clone https://github.com/cypone/ponysay-go.git
cd ponysay-go
go build -o ponysay ./cmd/ponysay
ln -s ponysay ponythink
```

### Option 2: Go Install

```bash
go install ponysay-go/cmd/ponysay@latest
```

---

## Usage Examples

### Basic Usage

```bash
ponysay "I am just the cutest pony!"
```

### Specific Pony

```bash
ponysay -f pinkie "Partay!~"
ponysay -f derpy "I brought muffins!"
```

### Pony Fortune (Shell Startup Hook)

Add this to your `~/.bashrc` or `~/.zshrc`:

```bash
fortune | ponysay
```

### Show Quotes (-q)

```bash
ponysay -q pinkie
ponysay -q
```

### Thought Bubble (ponythink)

```bash
ponythink "Hmm... is Golang fast?"
```

### Listing Options

```bash
ponysay -l           # List MLP ponies
ponysay +l          # List extra/non-MLP ponies
ponysay -A          # List all ponies
ponysay -L          # List ponies with alternative names
ponysay -B          # List balloon styles
ponysay --quoters   # List ponies that have quotes
```

---

## Running Unit Tests

```bash
go test -v ./...
```

---

## License

GPL-3.0 License.
