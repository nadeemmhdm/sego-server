[CmdletBinding()]
param(
    [ValidateSet('low-resource','standard')]
    [string]$Profile = 'low-resource',
    [string]$Prefix = "$env:ProgramFiles\SEGO",
    [string]$Data = "$env:ProgramData\SEGO",
    [switch]$Yes
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Test-Administrator {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($identity)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

if ($PSVersionTable.PSVersion.Major -lt 5) { throw 'PowerShell 5.1 or later is required.' }
if ([Environment]::OSVersion.Platform -ne [PlatformID]::Win32NT) { throw 'Use scripts/install.sh on Linux.' }
if (-not (Test-Administrator)) { throw 'Run PowerShell as Administrator for a system-wide SEGO installation.' }
if (-not (Get-Command go -ErrorAction SilentlyContinue)) { throw 'Go 1.23+ is required when installing from a source checkout.' }

$Config = Join-Path $Data 'sego.json'
Write-Host 'SEGO Server install plan' -ForegroundColor Cyan
Write-Host "  profile : $Profile"
Write-Host "  binaries: $Prefix\bin"
Write-Host "  data    : $Data"
Write-Host "  config  : $Config"
Write-Host '  startup : manual in the MVP source installer (native Windows service integration is tracked separately)'
Write-Host ''
Write-Host 'The installer never formats or repartitions disks. Existing user data is preserved.'

if (-not $Yes) {
    $answer = Read-Host 'Continue? [y/N]'
    if ($answer -notin @('y','Y','yes','YES','Yes')) { exit 0 }
}

$Bin = Join-Path $Prefix 'bin'
New-Item -ItemType Directory -Force -Path $Bin, $Data | Out-Null

Write-Host '[1/5] Building SEGO binaries'
go build -trimpath -o (Join-Path $Bin 'sego-server.exe') ./cmd/server
if ($LASTEXITCODE -ne 0) { throw 'Failed to build sego-server.exe' }
go build -trimpath -o (Join-Path $Bin 'sego.exe') ./cmd/sego
if ($LASTEXITCODE -ne 0) { throw 'Failed to build sego.exe' }

Write-Host '[2/5] Creating configuration'
& (Join-Path $Bin 'sego.exe') --config $Config mode set $Profile
if ($LASTEXITCODE -ne 0) { throw 'Failed to initialize SEGO configuration' }

Write-Host '[3/5] Running diagnostics'
& (Join-Path $Bin 'sego.exe') --config $Config doctor
if ($LASTEXITCODE -ne 0) { throw 'SEGO diagnostics failed' }

Write-Host '[4/5] Installation files ready'
Write-Host 'Native Windows service registration is not marked complete in this MVP; do not treat this source installer as a finished production service installer.' -ForegroundColor Yellow

Write-Host '[5/5] Installation complete' -ForegroundColor Green
Write-Host "Create the first administrator: & '$Bin\sego.exe' --config '$Config' admin create"
Write-Host "Start SEGO manually:            & '$Bin\sego-server.exe' --config '$Config'"
Write-Host "Set low-resource mode:          & '$Bin\sego.exe' --config '$Config' mode set low-resource"
Write-Host "Set standard mode:              & '$Bin\sego.exe' --config '$Config' mode set standard"
