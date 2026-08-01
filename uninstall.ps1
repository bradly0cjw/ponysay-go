# ponysay-go uninstaller script for Windows (PowerShell)
$ErrorActionPreference = 'Continue'

Write-Host "Uninstalling ponysay & ponythink..."

$locations = @(
    "$env:ProgramFiles\ponysay",
    "$env:LocalAppData\Programs\ponysay"
)

$removedAny = $false

foreach ($dir in $locations) {
    if (Test-Path -Path $dir) {
        Remove-Item -Path $dir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "Removed directory: $dir"
        $removedAny = $true
    }
}

# Clean PATH for User and Machine environment scopes
foreach ($target in @("User", "Machine")) {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", $target)
    if ($currentPath) {
        $pathEntries = $currentPath -split ';'
        $filteredEntries = $pathEntries | Where-Object { 
            $_.TrimEnd('\') -notlike "*\ponysay"
        }
        if ($pathEntries.Count -ne $filteredEntries.Count) {
            $newPath = $filteredEntries -join ';'
            [Environment]::SetEnvironmentVariable("Path", $newPath, $target)
            Write-Host "Removed ponysay from $target PATH environment variable."
            $removedAny = $true
        }
    }
}

# Remove config directory if present
$configDir = Join-Path -Path $env:APPDATA -ChildPath "ponysay"
if (Test-Path -Path $configDir) {
    Remove-Item -Path $configDir -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Removed user configuration directory: $configDir"
}

if ($removedAny) {
    Write-Host "Successfully uninstalled ponysay & ponythink!"
    Write-Host "Please restart your terminal/PowerShell window for PATH changes to take effect."
} else {
    Write-Host "No ponysay installations were found."
}
