[CmdletBinding()]
param(
    [string]$Prefix = "$env:ProgramFiles\SEGO",
    [string]$Data = "$env:ProgramData\SEGO",
    [switch]$PurgeData
)
$ErrorActionPreference = 'Stop'
if (Test-Path $Prefix) { Remove-Item -Recurse -Force $Prefix }
if ($PurgeData -and (Test-Path $Data)) {
    Remove-Item -Recurse -Force $Data
    Write-Host 'SEGO binaries, configuration and data removed.'
} else {
    Write-Host 'SEGO binaries removed. Configuration and user data were preserved.'
}
