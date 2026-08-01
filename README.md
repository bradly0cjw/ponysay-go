# ponysay-go

> A high-performance, zero-dependency, cross-platform Go port of [`ponysay`](https://github.com/erkin/ponysay) (cowsay reimplementation for ponies).

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux%20%7C%20FreeBSD-blue)
![License](https://img.shields.io/badge/license-GPLv3-green)

---

## ✨ Features

- 🚀 **Zero External Dependencies**: Single self-contained binary (`//go:embed` bundles all ~470+ pony art files, balloon styles, and quotes). No Python environment or `coreutils`/`stty` required.
- ⚡ **Blazing Fast**: Sub-5ms startup time — perfect for shell startup scripts (`fortune | ponysay`).
- 💻 **True Cross-Platform**: Native support for **Windows** (Command Prompt, PowerShell, Windows Terminal), **macOS**, **Linux**, and **FreeBSD** (`GOOS/GOARCH` support).
- 🎨 **Full Color Support**: ANSI, 256-color, and 24-bit TrueColor display.
- 💬 **Pony Quotes & Balloons**: Full support for speech balloons (`ponysay`), thought balloons (`ponythink`), quote databases (`-q`), custom balloon borders (`-b`), and word wrapping (`-W`).

---

## 📦 Installation

### Option 1: Via `go install` (Recommended)

```bash
go install ponysay-go/cmd/ponysay@latest
```

### Option 2: Build from Source

```bash
git clone https://github.com/cypone/ponysay-go.git
cd ponysay-go
go build -o ponysay ./cmd/ponysay
ln -s ponysay ponythink
```

---

## 🚀 Usage

### Basic Usage

```bash
ponysay "I am just the cutest pony!"
```

### Specific Pony

```bash
ponysay -f pinkie "Partay!~"
ponysay -f derpy "I brought muffins!"
```

### Pony Fortune (Shell Startup)

Add this to your `~/.bashrc` or `~/.zshrc`:

```bash
fortune | ponysay
```

### Pony Quotes (`-q`)

```bash
ponysay -q pinkie
ponysay -q
```

### Thought Bubble (`ponythink`)

```bash
ponythink "Hmm... is Golang fast?"
```

### List Available Ponies & Balloons

```bash
ponysay -l           # List MLP ponies
ponysay +l          # List extra/non-MLP ponies
ponysay -A          # List all ponies
ponysay -B          # List balloon styles
```

---

## 🛠️ Command-Line Options

| Option | Description |
| :--- | :--- |
| `-f`, `--file PONY` | Select a pony by name or file |
| `+f PONY` | Select a non-MLP pony |
| `-F PONY` | Select any pony (MLP or non-MLP) |
| `-q`, `--quote [PONY]` | Select a pony quote |
| `-b`, `--bubble STYLE` | Select balloon style (`unicode`, `cowsay`, `ascii`, `round`, etc.) |
| `-W`, `--wrap COLUMN` | Specify maximum wrapping width |
| `-l`, `--list` | List pony names |
| `+l` | List non-MLP pony names |
| `-A`, `--all` | List all pony names |
| `-B`, `--bubblelist` | List balloon styles |
| `-o`, `--pony-only` | Print only the pony artwork |
| `-v`, `--version` | Print program version |
| `-h`, `--help` | Display help menu |

---

## 🤝 Credits

Based on the original [`ponysay`](https://github.com/erkin/ponysay) created by Erkin Batu Altunbaş, Mattias "maandree" Andrée, and contributors. Artwork created by respective authors listed within each pony file metadata.

---

## 📜 License

GPL-3.0 License.
