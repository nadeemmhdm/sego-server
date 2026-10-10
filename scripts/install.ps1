[CmdletBinding()]
param(
    [ValidateSet('auto','low-resource','standard')][string]$Profile = 'auto',
    [string]$Prefix = "$env:ProgramFiles\SEGO",
    [string]$Data = "$env:ProgramData\SEGO",
    [switch]$Yes
)
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

function Step([string]$Text) { Write-Host "`nSEGO $Text" -ForegroundColor Cyan }
function Test-Administrator {
    $p = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    return $p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}
function Refresh-Path {
    $env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')
}
function Ensure-WingetTool([string]$Command,[string]$PackageId) {
    if (Get-Command $Command -ErrorAction SilentlyContinue) { return }
    if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
        throw "$Command is required, but Windows Package Manager (winget) is unavailable. Update 'App Installer' from Microsoft Store and run the same SEGO install command again."
    }
    Step "Installing required tool: $Command"
    & winget install --id $PackageId -e --accept-source-agreements --accept-package-agreements --silent
    if ($LASTEXITCODE -ne 0) { throw "Failed to install $Command." }
    Refresh-Path
    if (-not (Get-Command $Command -ErrorAction SilentlyContinue)) {
        throw "$Command was installed but Windows has not refreshed PATH yet. Close PowerShell, open it again, and run the same SEGO install command."
    }
}
function Get-GoVersion {
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) { return [version]'0.0' }
    $v = (& go env GOVERSION 2>$null) -replace '^go',''
    try { return [version]$v } catch { return [version]'0.0' }
}
function Ensure-Go {
    if ((Get-GoVersion) -ge [version]'1.24') { return }
    if (-not (Get-Command winget -ErrorAction SilentlyContinue)) { throw 'Go 1.24+ is required and winget is unavailable.' }
    Step 'Installing/updating Go'
    if (Get-Command go -ErrorAction SilentlyContinue) {
        & winget upgrade --id GoLang.Go -e --accept-source-agreements --accept-package-agreements --silent
        if ($LASTEXITCODE -ne 0) { & winget install --id GoLang.Go -e --accept-source-agreements --accept-package-agreements --silent }
    } else {
        & winget install --id GoLang.Go -e --accept-source-agreements --accept-package-agreements --silent
    }
    Refresh-Path
    if ((Get-GoVersion) -lt [version]'1.24') { throw 'Go 1.24+ could not be installed.' }
}
function Add-MachinePath([string]$PathToAdd) {
    $machine = [Environment]::GetEnvironmentVariable('Path','Machine')
    $parts = @($machine -split ';' | Where-Object { $_ })
    if ($parts -notcontains $PathToAdd) {
        [Environment]::SetEnvironmentVariable('Path', (($parts + $PathToAdd) -join ';'), 'Machine')
    }
    if (($env:Path -split ';') -notcontains $PathToAdd) { $env:Path += ';' + $PathToAdd }
}

if ($PSVersionTable.PSVersion.Major -lt 5) { throw 'PowerShell 5.1 or later is required.' }
if ([Environment]::OSVersion.Platform -ne [PlatformID]::Win32NT) { throw 'Use the Linux installer on Linux.' }

if (-not (Test-Administrator)) {
    Step 'Administrator permission is required once to install the background service'
    $elevateArgs = @(
        '-NoProfile','-ExecutionPolicy','Bypass','-File',('"' + $PSCommandPath + '"'),
        '-Profile',$Profile,'-Prefix',('"' + $Prefix + '"'),'-Data',('"' + $Data + '"')
    )
    if ($Yes) { $elevateArgs += '-Yes' }
    $p = Start-Process powershell.exe -Verb RunAs -ArgumentList $elevateArgs -Wait -PassThru
    exit $p.ExitCode
}

if ($Profile -eq 'auto') {
    try {
        $ram = (Get-CimInstance Win32_ComputerSystem).TotalPhysicalMemory
        if ($ram -lt 6GB) { $Profile = 'low-resource' } else { $Profile = 'standard' }
    } catch { $Profile = 'low-resource' }
}

$RepoRoot = Split-Path -Parent $PSScriptRoot
$Config = Join-Path $Data 'sego.json'
$Bin = Join-Path $Prefix 'bin'
$WebRoot = Join-Path $Data 'web'

Write-Host "`nSEGO is ready to install." -ForegroundColor Cyan
Write-Host "  Mode      : $Profile"
Write-Host '  Dashboard : http://127.0.0.1:8088'
Write-Host "  Data      : $Data"
Write-Host '  Startup   : automatic background startup with restart recovery'
Write-Host '  Updates   : automatic GitHub Release checks and verified package updates'
Write-Host 'No disk will be formatted or repartitioned. Existing SEGO data is preserved.'
if (-not $Yes) {
    $a = Read-Host 'Continue? [Y/n]'
    if ($a -match '^(n|no)$') { exit 0 }
}

Ensure-Go
Ensure-WingetTool 'node' 'OpenJS.NodeJS.LTS'
Ensure-WingetTool 'npm' 'OpenJS.NodeJS.LTS'
$nodeMajor = [int]((& node -p 'process.versions.node.split(".")[0]') | Select-Object -First 1)
if ($nodeMajor -lt 18) { throw "Node.js 18 or newer is required. Installed version: $(& node --version)" }

New-Item -ItemType Directory -Force -Path $Bin,$Data,$WebRoot | Out-Null
Push-Location $RepoRoot
try {
    Step '1/8 Resolving verified Go dependencies'
    & go mod tidy
    if ($LASTEXITCODE -ne 0) { throw 'go mod tidy failed.' }

    Step '2/8 Building the dashboard'
    Push-Location web
    try {
        & npm install --no-audit --no-fund
        if ($LASTEXITCODE -ne 0) { throw 'npm install failed.' }
        & npm run build
        if ($LASTEXITCODE -ne 0) { throw 'dashboard build failed.' }
    } finally { Pop-Location }

    Step '3/8 Building SEGO'
    & go build -trimpath -o (Join-Path $Bin 'sego-server.exe') ./cmd/server
    if ($LASTEXITCODE -ne 0) { throw 'Failed to build sego-server.exe.' }
    & go build -trimpath -o (Join-Path $Bin 'sego.exe') ./cmd/sego
    if ($LASTEXITCODE -ne 0) { throw 'Failed to build sego.exe.' }
    & go build -trimpath -o (Join-Path $Bin 'sego-updater.exe') ./cmd/updater
    if ($LASTEXITCODE -ne 0) { throw 'Failed to build sego-updater.exe.' }

    Step '4/8 Installing dashboard and configuration'
    if (Test-Path $WebRoot) { Remove-Item $WebRoot -Recurse -Force }
    New-Item -ItemType Directory -Force -Path $WebRoot | Out-Null
    Copy-Item (Join-Path $RepoRoot 'web\dist\*') $WebRoot -Recurse -Force
    $oldData = $env:SEGO_DATA_DIR
    $env:SEGO_DATA_DIR = $Data
    try {
        & (Join-Path $Bin 'sego.exe') --config $Config mode $Profile
        if ($LASTEXITCODE -ne 0) { throw 'Failed to initialize SEGO configuration.' }
    } finally { $env:SEGO_DATA_DIR = $oldData }

    Step '5/8 Protecting server data'
    & icacls $Data /inheritance:r /grant:r '*S-1-5-18:(OI)(CI)F' '*S-1-5-32-544:(OI)(CI)F' | Out-Null

    Step '6/8 Enabling automatic background startup and release updates'
    $principal = New-ScheduledTaskPrincipal -UserId 'SYSTEM' -LogonType ServiceAccount -RunLevel Highest
    $serverExe = Join-Path $Bin 'sego-server.exe'
    $serverAction = New-ScheduledTaskAction -Execute $serverExe -Argument "--config `"$Config`""
    $serverTrigger = New-ScheduledTaskTrigger -AtStartup
    $serverSettings = New-ScheduledTaskSettingsSet -StartWhenAvailable -RestartCount 10 -RestartInterval (New-TimeSpan -Minutes 1) -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -ExecutionTimeLimit ([TimeSpan]::Zero) -MultipleInstances IgnoreNew
    Register-ScheduledTask -TaskName 'SEGO Server' -Action $serverAction -Trigger $serverTrigger -Principal $principal -Settings $serverSettings -Description 'SEGO Server private cloud background service' -Force | Out-Null
    Start-ScheduledTask -TaskName 'SEGO Server'

    $updateExe = Join-Path $Bin 'sego-updater.exe'
    $updateAction = New-ScheduledTaskAction -Execute $updateExe -Argument "--config `"$Config`""
    $updateTriggers = @((New-ScheduledTaskTrigger -AtStartup),(New-ScheduledTaskTrigger -Daily -At 4am))
    $updateSettings = New-ScheduledTaskSettingsSet -StartWhenAvailable -RestartCount 5 -RestartInterval (New-TimeSpan -Minutes 5) -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -ExecutionTimeLimit (New-TimeSpan -Minutes 30) -MultipleInstances IgnoreNew
    Register-ScheduledTask -TaskName 'SEGO Update' -Action $updateAction -Trigger $updateTriggers -Principal $principal -Settings $updateSettings -Description 'SEGO GitHub Release update checker and installer' -Force | Out-Null
    Add-MachinePath $Bin

    Step '7/8 Creating your administrator account'
    & (Join-Path $Bin 'sego.exe') --config $Config setup
    if ($LASTEXITCODE -ne 0) { throw 'Administrator setup failed.' }

    Step '8/8 Checking SEGO'
    & (Join-Path $Bin 'sego.exe') --config $Config doctor
    if ($LASTEXITCODE -ne 0) { throw 'SEGO diagnostics failed.' }

    Write-Host "`nSEGO Server is installed." -ForegroundColor Green
    Write-Host 'Open SEGO:    sego'
    Write-Host 'Check SEGO:   sego status'
    Write-Host 'Remote once:  sego remote'
    Write-Host 'Update now:   sego update'
    Write-Host 'Update check: sego update check'
    Write-Host 'Dashboard:    http://127.0.0.1:8088'
    Write-Host 'After a reboot SEGO starts in the background automatically.'
    Write-Host 'If a tunnel was configured once, cloudflared reconnects automatically too.'
    try { Start-Process 'http://127.0.0.1:8088' } catch {}
}
finally { Pop-Location }
