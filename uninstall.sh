#!/bin/sh
set -e

# ponysay-go uninstaller script (Linux & macOS)

echo "Uninstalling ponysay & ponythink..."

REMOVED=0

# Standard binary targets
TARGETS="
/usr/local/bin/ponysay
/usr/local/bin/ponythink
${HOME}/.local/bin/ponysay
${HOME}/.local/bin/ponythink
/usr/bin/ponysay
/usr/bin/ponythink
"

# Add dynamic lookup if command exists in PATH
if command -v ponysay >/dev/null 2>&1; then
    P_PATH="$(command -v ponysay)"
    TARGETS="${TARGETS} ${P_PATH}"
fi
if command -v ponythink >/dev/null 2>&1; then
    PT_PATH="$(command -v ponythink)"
    TARGETS="${TARGETS} ${PT_PATH}"
fi

for file in $TARGETS; do
    if [ -f "$file" ] || [ -L "$file" ] || [ -e "$file" ]; then
        dir="$(dirname "$file")"
        if [ -w "$dir" ] || [ "$(id -u)" -eq 0 ]; then
            rm -f "$file"
            echo "Removed: $file"
            REMOVED=$((REMOVED + 1))
        else
            echo "Warning: Permission denied to remove $file. Try running with sudo: sudo ./uninstall.sh" >&2
        fi
    fi
done

# Clean up user config folder if present
if [ -d "${HOME}/.config/ponysay" ]; then
    rm -rf "${HOME}/.config/ponysay"
    echo "Removed user configuration directory: ${HOME}/.config/ponysay"
fi

# Clean up terminal startup hooks from shell RC files
MARKER="# ponysay-go terminal greeting"
RC_FILES="${HOME}/.bashrc ${HOME}/.zshrc ${HOME}/.config/fish/config.fish"
for rc in $RC_FILES; do
    if [ -f "$rc" ] && grep -qF "$MARKER" "$rc" 2>/dev/null; then
        # Remove the ponysay greeting block (marker line + following non-empty lines until blank line or EOF)
        if [ "$(basename "$rc")" = "config.fish" ]; then
            # Fish: remove from marker through 'end'
            sed -i.bak "/$MARKER/,/^end$/d" "$rc" && rm -f "${rc}.bak"
        else
            # Bash/Zsh: remove the marker line and the if-then-fi one-liner
            sed -i.bak "/$MARKER/{N;d;}" "$rc" && rm -f "${rc}.bak"
        fi
        # Remove any trailing blank line left behind
        sed -i.bak -e :a -e '/^\n*$/{$d;N;ba' -e '}' "$rc" 2>/dev/null && rm -f "${rc}.bak"
        echo "Removed terminal startup hook from $rc"
    fi
done

if [ "$REMOVED" -gt 0 ]; then
    echo "Successfully uninstalled ponysay & ponythink!"
else
    echo "No ponysay installations were found."
fi
