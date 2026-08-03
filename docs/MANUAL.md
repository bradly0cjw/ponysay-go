# ponysay-go Manual

`ponysay-go` is a cross-platform Go port of [ponysay](https://github.com/erkin/ponysay) (a cowsay reimplementation for ponies). This document serves as the comprehensive user and reference manual for `ponysay-go`.

---

## Table of Contents

- [ponysay-go Manual](#ponysay-go-manual)
  - [Table of Contents](#table-of-contents)
  - [Synopsis](#synopsis)
  - [Execution Modes (`ponysay` vs `ponythink`)](#execution-modes-ponysay-vs-ponythink)
  - [Command Line Options Reference](#command-line-options-reference)
    - [Pony Selection Options](#pony-selection-options)
    - [Quote Database Options](#quote-database-options)
    - [Balloon \& Formatting Options](#balloon--formatting-options)
    - [Color Engine Options](#color-engine-options)
    - [Listing \& Discovery Options](#listing--discovery-options)
    - [Information \& Artwork Display](#information--artwork-display)
    - [Maintenance \& System Commands](#maintenance--system-commands)
  - [Input Handling \& Message Fallbacks](#input-handling--message-fallbacks)
  - [Smart Search \& Fuzzy Spell Correction (`SpelloCorrecter`)](#smart-search--fuzzy-spell-correction-spellocorrecter)
  - [Unicode Accent Remapping (`PONYSAY_UCS_ME`)](#unicode-accent-remapping-ponysay_ucs_me)
  - [Asset Resolution Engine \& Directory Hierarchy](#asset-resolution-engine--directory-hierarchy)
  - [Terminal Width Auto-Fitting](#terminal-width-auto-fitting)
  - [Environment Variables](#environment-variables)
  - [Usage Examples](#usage-examples)
    - [Interactive Speech \& Thought](#interactive-speech--thought)
    - [Random Quotes \& Quoter List](#random-quotes--quoter-list)
    - [Shell Startup Integration](#shell-startup-integration)
    - [Custom Ponies](#custom-ponies)
  - [Troubleshooting \& FAQ](#troubleshooting--faq)

---

## Synopsis

```bash
ponysay [OPTIONS] [message]
ponythink [OPTIONS] [message]
```

---

## Execution Modes (`ponysay` vs `ponythink`)

`ponysay-go` changes speech balloon behavior based on how the executable is invoked:
- **`ponysay`**: Displays speech balloons with standard speech stems/pointers (`\`, `/`).
- **`ponythink`**: Displays thought balloons with thought bubbles (`o`, `.`).

When installing, `ponythink` can be created as a symbolic link (`ln -s ponysay ponythink` on Unix) or copy (`copy ponysay.exe ponythink.exe` on Windows).

---

## Command Line Options Reference

### Pony Selection Options

| Flag | Short / Alternate | Description |
| :--- | :--- | :--- |
| `--file PONY` | `-f PONY`, `-file PONY`, `-pony PONY` | Select a specific standard MLP:FiM pony. |
| `--file PONY` | `+f PONY`, `++file PONY`, `++pony PONY` | Select a non-MLP extra pony. |
| `--file PONY` | `-F PONY`, `+F PONY`, `--any-ponies` | Select any pony (standard MLP or extra). |
| `--files PONY...` | `--f`, `++f`, `--F` | Select randomly from a variadic list of ponies. |

### Quote Database Options

| Flag | Short / Alternate | Description |
| :--- | :--- | :--- |
| `--quote [PONY]` | `-q [PONY]`, `+q [PONY]`, `-quote`, `--quote` | Print a quote from specified pony, or a random quote if no pony is specified. |
| `--quotes [PONY...]` | `--q`, `++q`, `--quotes`, `++quotes` | Variadic quote selection from specified ponies. |

### Balloon & Formatting Options

| Flag | Short / Alternate | Description |
| :--- | :--- | :--- |
| `--bubble STYLE` | `-b STYLE`, `-bubble`, `-balloon`, `--balloon` | Select balloon border style (`cowsay`, `unicode`, `ascii`, `round`, etc.). |
| `--wrap COLUMN` | `-W COLUMN`, `-wrap` | Set custom maximum text wrapping column width. |
| `--compress` | `-c`, `-compress`, `--compact` | Compress consecutive blank lines in the message body. |
| `--colour-bubble COLOR` | `+c COLOR`, `--colour-balloon COLOR` | Override balloon border color code (ANSI escape). |
| `--colour-link COLOR` | `--colour-link COLOR` | Override speech/thought stem color code. |
| `--colour-msg COLOR` | `--colour-message COLOR` | Override message text color code. |

### Color Engine Options

| Flag | Short / Alternate | Description |
| :--- | :--- | :--- |
| `--256-colours` | `-X`, `--256colours` | Enable 256-color palette mode (default on modern terminals). |
| `--tty-colours` | `-V`, `--ttycolours` | Force 16-color TTY mode for legacy consoles. |
| `--kms-colours` | `-K` | Enable Kernel Mode Setting (KMS) Linux console colors. |

### Listing & Discovery Options

| Flag | Short / Alternate | Description |
| :--- | :--- | :--- |
| `--list` | `-l` | List standard MLP pony names. |
| `--symlist` | `-L`, `--altlist` | List standard MLP pony names alongside their aliases. |
| `+l` | `++list` | List extra (non-MLP) pony names. |
| `+L` | `++symlist`, `++altlist` | List extra (non-MLP) pony names alongside their aliases. |
| `--all` | `-A` | List all pony names (MLP and extra). |
| `+A` | `++all`, `++symall`, `++altall` | List all pony names alongside their aliases. |
| `--bubblelist` | `-B`, `--balloonlist` | List available balloon styles. |
| `--quoters` | | List all ponies that have quotes in the quote database. |
| `--onelist` | `++onelist`, `--Onelist` | Output listing results on a single line separated by spaces. |

### Information & Artwork Display

| Flag | Short / Alternate | Description |
| :--- | :--- | :--- |
| `--info` | `-i` | Display metadata info for the pony (name, author, license, etc.). |
| `++info` | `+i` | Display metadata info formatted with bold/colored headers. |
| `--pony-only` | `-o`, `--ponyonly` | Output only the pony artwork without a balloon or message. |

### Maintenance & System Commands

| Flag / Command | Short / Alternate | Description |
| :--- | :--- | :--- |
| `update` | `-u`, `--update` | Automatically check, download, and replace binary with latest GitHub release. |
| `--version` | `-v` | Display version and commit hash details. |
| `--help` | `-h` | Display the built-in help summary. |

---

## Input Handling & Message Fallbacks

`ponysay-go` resolves the message to display using the following precedence:

1. **Quote Mode (`-q`)**: If quote mode is active, the quote database supplies the text.
2. **Explicit Command-Line Arguments**: Any non-flag arguments passed to `ponysay` are joined into the message string.
3. **Standard Input Pipe (`stdin`)**: If no command-line text is provided and input is piped (e.g. `fortune | ponysay`), `ponysay` reads from `stdin`.
4. **Default Message**: If no input is piped or passed, and `--pony-only` is not set, `ponysay` defaults to:
   ```text
   I am just the cutest pony!
   ```

---

## Smart Search & Fuzzy Spell Correction (`SpelloCorrecter`)

`ponysay-go` features a weighted Levenshtein spell corrector (`SpelloCorrecter`) to handle typos and approximate pony names:

- **Phonetic & Substitution Weights**: Common substitutions (e.g. `k`↔`c`, `s`↔`z`, `o`↔`u`, `v`↔`w`, `-`↔`_`, case mismatches) carry reduced edit weights.
- **Distance Threshold (`PONYSAY_TYPO_LIMIT`)**: Candidates within distance $\le 5$ (or custom threshold) are automatically matched to the nearest valid pony file name.
- **Error Behavior**: If no pony name matches within the distance threshold, `ponysay-go` outputs:
  ```text
  Error: I have never heard of anypony named <input>
  ```
  and exits with status `1`.

---

## Unicode Accent Remapping (`PONYSAY_UCS_ME`)

Certain pony files have Unicode accented names (e.g. `mjölna`, `bifröst`, `piñacolada`). Setting the environment variable `PONYSAY_UCS_ME` remaps accented inputs to their ASCII file equivalents:

- **Accepted Values**: `1`, `yes`, `y`, `harder`, `h`, `2`
- **Example**: `PONYSAY_UCS_ME=1 ponysay -f mjölna` automatically resolves to `mjolna.pony`.

---

## Asset Resolution Engine & Directory Hierarchy

`ponysay-go` utilizes a hybrid asset resolution engine. When a pony, balloon, or quote file is requested, directories are searched in the following priority order:

1. **Current Directory**: `./ponies/`, `./extraponies/`, `./balloons/`
2. **User Config Directory**:
   - Linux/macOS: `~/.config/ponysay/ponies/`
   - Windows: `%LocalAppData%\ponysay\ponies\`
3. **User Home Directory**: `~/.ponysay/ponies/`
4. **System Directories**: `/usr/share/ponysay/ponies/`, `/usr/local/share/ponysay/ponies/`
5. **Embedded Assets**: Embedded fallback files compiled into the `ponysay-go` binary (`//go:embed`).

---

## Terminal Width Auto-Fitting

When selecting a random pony:
1. `ponysay-go` detects active terminal dimensions via native syscalls (`golang.org/x/term`).
2. It strips ANSI escape sequences to compute exact visual column width of each pony candidate.
3. Candidates exceeding terminal width are filtered out to prevent visual wrap artifacts.
4. If present in asset directories, `best.pony` serves as the preferred fallback when no pony selection flag is passed.

---

## Environment Variables

| Variable | Default | Description |
| :--- | :--- | :--- |
| `PONYSAY_TYPO_LIMIT` | `5` | Maximum edit distance threshold for fuzzy spell correction. |
| `PONYSAY_UCS_ME` | (unset) | Set to `1`, `yes`, `y`, `harder`, `h`, or `2` to enable Unicode name remapping to ASCII filenames. |
| `NO_COLOR` | (unset) | If set, disables ANSI color output in compatible environments. |
| `TERM` | `xterm-256color` | Terminal capability identifier used for color palette detection. |

---

## Usage Examples

### Interactive Speech & Thought
```bash
# Speech balloon
ponysay -f twilight "Books are wonderful!"

# Thought balloon
ponythink -f starlight "Should I try time travel again?"
```

### Random Quotes & Quoter List
```bash
# Random quote from any pony in database
ponysay -q

# Specific quote from Pinkie Pie
ponysay -q pinkie

# Specific quote from multiple ponies
ponysay --quotes twilight spike

# List all ponies that have quotes
ponysay --quoters # or ponysay --q
```

### Shell Startup Integration

Unix shell (macOS / Linux) startup example (`~/.bashrc` or `~/.zshrc`):
```bash
# Add to end of shell configuration file
if command -v ponysay >/dev/null 2>&1; then
    fortune | ponysay
    # or you just want a random quote on shell startup
    ponysay -q
fi
```
Windows PowerShell equivalent:
```powershell
# Add to end of PowerShell profile
if (Get-Command ponysay -ErrorAction SilentlyContinue) {
    ponysay -q
}
```

### Custom Ponies
```bash
# Place custom pony artwork in config folder
mkdir -p ~/.config/ponysay/ponies
cp mypony.pony ~/.config/ponysay/ponies/

# Render custom pony
ponysay -f mypony "Custom pony ready!"
```

---

## Troubleshooting & FAQ

**Q: `ponythink` command not found?**
- Ensure you created a symlink or copy:
  - Unix: `ln -s $(which ponysay) ~/.local/bin/ponythink`
  - Windows PowerShell: `Copy-Item (Get-Command ponysay).Path -Destination "%LocalAppData%\Programs\ponysay\ponythink.exe"`

**Q: Colors look incorrect on Linux TTY console?**
- Use the TTY 16-color mode flag: `ponysay -V "Hello world"` or KMS mode: `ponysay -K "Hello world"`.

**Q: How do I update `ponysay-go`?**
- Run `ponysay update` (or `ponysay -u`) to automatically fetch and replace the binary with the latest GitHub release.
