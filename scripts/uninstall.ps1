[CmdletBinding()]
param(
    [string]$Prefix = "$env:ProgramFiles\SEGO",
    [string]$Data = "$env:ProgramData\SEGO",
    [switch]$PurgeData
)
$ErrorActionPreference = 'Stop'
$task='SEGO Server'
if(Get-ScheduledTask -TaskName $task -ErrorAction SilentlyContinue){
    Stop-ScheduledTask -TaskName $task -ErrorAction SilentlyContinue
    Unregister-ScheduledTask -TaskName $task -Confirm:$false
}
if(Test-Path $Prefix){Remove-Item -Recurse -Force $Prefix}
if($PurgeData -and (Test-Path $Data)){
    Remove-Item -Recurse -Force $Data
    Write-Host 'SEGO binaries, configuration and data removed.'
}else{
    Write-Host 'SEGO binaries and startup task removed. Configuration and user data were preserved.'
}
