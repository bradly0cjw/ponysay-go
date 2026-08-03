# ponysay-go installer script for Windows (PowerShell)
param(
    [switch]$Terminal
)
$ErrorActionPreference = 'Stop'

$repo = "bradly0cjw/ponysay-go"

# 1. Detect Architecture
$arch = $env:PROCESSOR_ARCHITECTURE
switch -regex ($arch) {
    'AMD64|x86_64' { $goarch = 'amd64' }
    'ARM64|aarch64' { $goarch = 'arm64' }
    default {
        Write-Error "Unsupported architecture: $arch"
        exit 1
    }
}

$assetName = "ponysay-windows-$goarch.exe"
$downloadUrl = "https://github.com/$repo/releases/latest/download/$assetName"

Write-Host "Detected Windows OS (Arch: $goarch)"
Write-Host "Downloading $assetName..."

# 2. Determine Install Directory
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if ($isAdmin) {
    $installDir = "$env:ProgramFiles\ponysay"
    $pathTarget = "Machine"
} else {
    $installDir = "$env:LocalAppData\Programs\ponysay"
    $pathTarget = "User"
}

if (-not (Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

$ponysayExe = Join-Path -Path $installDir -ChildPath "ponysay.exe"
$ponythinkExe = Join-Path -Path $installDir -ChildPath "ponythink.exe"

# 3. Download Binary
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-WebRequest -Uri $downloadUrl -OutFile $ponysayExe -UseBasicParsing

# 4. Create copy for ponythink
Copy-Item -Path $ponysayExe -Destination $ponythinkExe -Force

Write-Host "Successfully installed ponysay.exe & ponythink.exe to $installDir!"

# 5. Check and update PATH environment variable
$currentPath = [Environment]::GetEnvironmentVariable("Path", $pathTarget)
$pathEntries = $currentPath -split ';' | ForEach-Object { $_.TrimEnd('\') }
if ($pathEntries -notcontains $installDir.TrimEnd('\')) {
    $newPath = "$currentPath;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, $pathTarget)
    $env:Path = "$env:Path;$installDir"
    Write-Host "Added $installDir to your $pathTarget PATH."
    Write-Host "Please restart your terminal/PowerShell window for PATH changes to take effect."
}

Write-Host ""
Write-Host "Try running:"
Write-Host "    ponysay `"I am just the cutest pony!`""

# 6. Terminal startup hook (-Terminal flag)
if ($Terminal) {
    $marker = '# ponysay-go terminal greeting'
    $profilePath = $PROFILE.CurrentUserCurrentHost

    if ($profilePath -and (Test-Path -Path $profilePath) -and (Select-String -Path $profilePath -Pattern ([regex]::Escape($marker)) -Quiet)) {
        Write-Host ""
        Write-Host "Terminal hook already present in $profilePath, skipping."
    } else {
        Write-Host ""
        Write-Host "Adding ponysay startup hook to $profilePath..."
        $profileDir = Split-Path -Path $profilePath -Parent
        if (-not (Test-Path -Path $profileDir)) {
            New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
        }
        $hookBlock = @"

$marker
if (Get-Command ponysay -ErrorAction SilentlyContinue) { ponysay -q }
"@
        Add-Content -Path $profilePath -Value $hookBlock
        Write-Host "Done! A random pony quote will greet you on every new PowerShell session."
        Write-Host "To remove it later, delete the 'ponysay-go terminal greeting' block from $profilePath."
    }
}
