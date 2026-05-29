# Production build: portable exe + Windows NSIS installer.
# Run from the repository root: .\scripts\build-installer.ps1

$ErrorActionPreference = "Stop"

function Require-Command {
    param(
        [string]$Name,
        [string]$InstallHint
    )
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        Write-Error "$Name not found in PATH.`n$InstallHint"
    }
}

$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

Write-Host "==> Go Code Visualizer — production build" -ForegroundColor Cyan
Write-Host "    repo: $repoRoot"

Require-Command "go" "Install Go 1.22+: https://go.dev/dl/"
Require-Command "node" "Install Node.js 18+: https://nodejs.org/"
Require-Command "wails" "Install Wails CLI: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
Require-Command "makensis" "Install NSIS: winget install NSIS.NSIS"

Write-Host "==> Building frontend + app + NSIS installer..." -ForegroundColor Cyan
wails build -clean -nsis -trimpath
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

$binDir = Join-Path $repoRoot "build\bin"
$portable = Join-Path $binDir "visualisator.exe"
$installer = Join-Path $binDir "visualisator-amd64-installer.exe"

Write-Host ""
Write-Host "==> Build complete" -ForegroundColor Green
if (Test-Path $portable) {
    Write-Host "    Portable:  $portable"
} else {
    Write-Warning "Portable exe not found at $portable"
}
if (Test-Path $installer) {
    Write-Host "    Installer: $installer"
} else {
    Write-Warning "Installer not found at $installer (is makensis on PATH?)"
}
