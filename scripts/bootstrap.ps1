# SEGO Server beginner installer bootstrap.
# Run in PowerShell:
#   irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
$Repo = 'nadeemmhdm/sego-server'
$Tmp = Join-Path $env:TEMP ('sego-install-' + [guid]::NewGuid())

function Step([string]$Text) { Write-Host "`nSEGO $Text" -ForegroundColor Cyan }
function Fail([string]$Text) { throw "SEGO install error: $Text" }

Step 'Server installer'
Write-Host 'This will automatically download SEGO from GitHub, install required tools,'
Write-Host 'create automatic startup, and guide you through the first admin account.'
$answer = Read-Host 'Continue? [Y/n]'
if ($answer -match '^(n|no)$') { Write-Host 'Installation cancelled.'; return }

New-Item -ItemType Directory -Force -Path $Tmp | Out-Null
try {
    Step 'Downloading the current GitHub version'
    $meta = Invoke-RestMethod -Headers @{ Accept = 'application/vnd.github+json' } -Uri "https://api.github.com/repos/$Repo/commits/main"
    $commit = [string]$meta.sha
    if ($commit -notmatch '^[0-9a-f]{40}$') { Fail 'Could not resolve the GitHub main commit.' }

    $zip = Join-Path $Tmp 'sego.zip'
    Invoke-WebRequest -UseBasicParsing -Uri "https://github.com/$Repo/archive/$commit.zip" -OutFile $zip
    Expand-Archive -Path $zip -DestinationPath $Tmp -Force
    $installer = Join-Path $Tmp "sego-server-$commit\scripts\install.ps1"
    if (-not (Test-Path $installer)) { Fail 'Downloaded installer is missing.' }

    Step 'Starting guided setup'
    & $installer -Profile auto -Yes
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}
finally {
    Remove-Item $Tmp -Recurse -Force -ErrorAction SilentlyContinue
}
