[CmdletBinding()]
param([string]$Prefix="$env:ProgramFiles\SEGO",[string]$Data="$env:ProgramData\SEGO")
$ErrorActionPreference='Stop'
function Test-Admin{$p=New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent());$p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)}
if(-not(Test-Admin)){$args=@('-NoProfile','-ExecutionPolicy','Bypass','-File',('"'+$PSCommandPath+'"'),'-Prefix',('"'+$Prefix+'"'),'-Data',('"'+$Data+'"'));$p=Start-Process powershell.exe -Verb RunAs -ArgumentList $args -Wait -PassThru;exit $p.ExitCode}
$env:Path=[Environment]::GetEnvironmentVariable('Path','Machine')+';'+[Environment]::GetEnvironmentVariable('Path','User')
$root=Split-Path -Parent $PSScriptRoot;$bin=Join-Path $Prefix 'bin';$config=Join-Path $Data 'sego.json';$updater=Join-Path $bin 'sego-updater.exe'
Write-Host "`nSEGO Configuring automatic release updates" -ForegroundColor Cyan
Push-Location $root
try{& go build -trimpath -o $updater ./cmd/updater;if($LASTEXITCODE-ne 0){throw 'Failed to build SEGO updater.'}}finally{Pop-Location}
$action=New-ScheduledTaskAction -Execute $updater -Argument "--config `"$config`""
$triggers=@((New-ScheduledTaskTrigger -AtStartup),(New-ScheduledTaskTrigger -Daily -At 4am))
$principal=New-ScheduledTaskPrincipal -UserId 'SYSTEM' -LogonType ServiceAccount -RunLevel Highest
Register-ScheduledTask -TaskName 'SEGO Update' -Action $action -Trigger $triggers -Principal $principal -Description 'SEGO GitHub Release update checker' -Force|Out-Null
Write-Host '[OK] Automatic GitHub Release checks enabled.' -ForegroundColor Green
