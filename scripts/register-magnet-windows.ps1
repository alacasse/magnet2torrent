#!/usr/bin/env pwsh
# Register magnet2torrent as the handler for magnet: links on Windows (per-user registry).

param(
  [string] $AppName = $(if ($env:APP_NAME) { $env:APP_NAME } else { "magnet2torrent" }),
  [switch] $DryRun
)

function Usage {
  Write-Host "Usage: pwsh scripts/register-magnet-windows.ps1 [-AppName magnet2torrent] [-DryRun]"
}

function Resolve-BinaryPath {
  param([string] $Name)
  $cmd = Get-Command $Name -ErrorAction SilentlyContinue
  if (-not $cmd) {
    Write-Error "Could not find '$Name' in PATH; install it first."
    exit 1
  }
  $cmd.Path
}

function Set-RegValue {
  param(
    [string] $Path,
    [string] $Name,
    [string] $Value
  )
  if ($DryRun) {
    Write-Host "[dry-run] Set $Path [$Name] = $Value"
    return
  }
  if (-not (Test-Path $Path)) {
    New-Item -Path $Path -Force | Out-Null
  }
  New-ItemProperty -Path $Path -Name $Name -Value $Value -PropertyType String -Force | Out-Null
}

function Register-Protocol {
  param([string] $BinaryPath)

  $magnetKey = "HKCU:\Software\Classes\magnet"
  $commandKey = Join-Path $magnetKey "shell\open\command"
  $defaultIconKey = Join-Path $magnetKey "DefaultIcon"

  Set-RegValue -Path $magnetKey -Name "(default)" -Value "URL:Magnet Protocol"
  Set-RegValue -Path $magnetKey -Name "URL Protocol" -Value ""
  Set-RegValue -Path $defaultIconKey -Name "(default)" -Value "$BinaryPath,0"
  Set-RegValue -Path $commandKey -Name "(default)" -Value "`"$BinaryPath`" `"%1`""

  Write-Host "Registered magnet: handler to $BinaryPath (per-user under HKCU)"
  Write-Host "If Windows still prompts for a handler, choose magnet2torrent for magnet links in Settings > Apps > Default apps > Choose defaults by link type > magnet."
}

if ($args -contains "--help" -or $args -contains "-h") {
  Usage
  exit 0
}

$binary = Resolve-BinaryPath -Name $AppName
Register-Protocol -BinaryPath $binary
