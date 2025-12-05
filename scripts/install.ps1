#!/usr/bin/env pwsh
# Extensible installer for magnet2torrent via npm tarball from GitHub Releases (Windows).
# Customize hooks to prompt for config or switch install strategies.

param(
  [string] $RepoOwner = $env:REPO_OWNER,
  [string] $RepoName = $env:REPO_NAME,
  [string] $Version = $env:VERSION,
  [string] $Prefix = $env:PREFIX,
  [switch] $DryRun,
  [bool] $RegisterMagnet = $true
)

if (-not $RepoOwner) { $RepoOwner = "alacasse" }
if (-not $RepoName) { $RepoName = "magnet2torrent" }
if ($env:REGISTER_MAGNET) { $RegisterMagnet = -not ($env:REGISTER_MAGNET -eq "0") }

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Usage {
  Write-Host "Usage: pwsh scripts/install.ps1 -RepoOwner you -RepoName repo -Version 0.1.0 [-Prefix C:\path\to\npm] [-DryRun] [-RegisterMagnet \$true|\$false]"
  Write-Host ""
  Write-Host "Env vars: REGISTER_MAGNET=0 to skip magnet handler registration"
}

function Ensure-Version {
  param([string] $V)
  if (-not $V) {
    Write-Error "Version is required (e.g., -Version 0.1.0 or set VERSION env var)."
    Usage
    exit 1
  }
}

function Tarball-Url {
  param([string] $Owner, [string] $Name, [string] $V)
  "https://github.com/$Owner/$Name/releases/download/v$V/magnet2torrent-$V.tgz"
}

function Run {
  param([string] $Cmd)
  if ($DryRun) {
    Write-Host "[dry-run] $Cmd"
  } else {
    iex $Cmd
  }
}

function Pre-InstallConfig {
  # Hook for future interactive config collection.
}

function Register-MagnetWindows {
  if (-not $RegisterMagnet) {
    Write-Host "Skipping magnet handler registration (RegisterMagnet=$RegisterMagnet)."
    return
  }

  $registerScript = Join-Path $ScriptDir "register-magnet-windows.ps1"
  if (-not (Test-Path $registerScript)) {
    Write-Warning "Registration script not found at $registerScript"
    return
  }

  Write-Host "Registering magnet handler..."
  $cmd = "pwsh `"$registerScript`""
  if ($DryRun) { $cmd += " -DryRun" }
  Run $cmd
}

function Install-Windows {
  param([string] $Url, [string] $MaybePrefix)
  Write-Host "Installing magnet2torrent $Version for Windows (npm global)..."
  if ($MaybePrefix) {
    Run "npm install -g --prefix `"$MaybePrefix`" `"$Url`""
  } else {
    Run "npm install -g `"$Url`""
  }
  Register-MagnetWindows
}

function Post-InstallMessage {
  param([string] $MaybePrefix)
  Write-Host "magnet2torrent installed. Verify with: magnet2torrent --help"
  if ($MaybePrefix) {
    Write-Host "If not on PATH, add: setx PATH `"$MaybePrefix\bin;%PATH%`""
  }
  if (-not $RegisterMagnet) {
    Write-Host "Magnet handler registration skipped. Run scripts/register-magnet-windows.ps1 later if needed."
  }
}

Ensure-Version -V $Version
Pre-InstallConfig
$url = Tarball-Url -Owner $RepoOwner -Name $RepoName -V $Version
Install-Windows -Url $url -MaybePrefix $Prefix
Post-InstallMessage -MaybePrefix $Prefix
