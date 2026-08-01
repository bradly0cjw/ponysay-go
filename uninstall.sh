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

if [ "$REMOVED" -gt 0 ]; then
    echo "Successfully uninstalled ponysay & ponythink!"
else
    echo "No ponysay installations were found."
fi
