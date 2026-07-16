#Requires -Version 5.1
<#
.SYNOPSIS
  Reproduces this PC's Claude Code CLI environment (native binary + global
  settings: model, language, effort level, theme, enabled plugins/marketplaces)
  on another Windows PC.

.NOTES
  Deliberately does NOT touch API keys/OAuth tokens — those are per-machine
  secrets and must be set up by logging in on the new machine (see README.md).
#>

$ErrorActionPreference = "Stop"

Write-Host "== Claude Code CLI environment setup ==" -ForegroundColor Cyan

# 1. Install the Claude Code CLI (native binary), if not already present.
$claudeCmd = Get-Command claude -ErrorAction SilentlyContinue
if ($claudeCmd) {
    Write-Host "claude CLI already installed at $($claudeCmd.Source) - skipping install." -ForegroundColor Yellow
} else {
    Write-Host "Installing Claude Code CLI..."
    irm https://claude.sh | iex
}

# 2. Write global settings.json (model/language/effort/theme/plugins), without
#    clobbering an existing one that might have its own apiKeyHelper etc.
$claudeConfigDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { Join-Path $env:USERPROFILE ".claude" }
New-Item -ItemType Directory -Force -Path $claudeConfigDir | Out-Null
$settingsPath = Join-Path $claudeConfigDir "settings.json"
$templatePath = Join-Path $PSScriptRoot "settings.template.json"

if (Test-Path $settingsPath) {
    $backupPath = "$settingsPath.bak.$(Get-Date -Format yyyyMMdd-HHmmss)"
    Copy-Item $settingsPath $backupPath
    Write-Host "Existing settings.json found - backed up to $backupPath" -ForegroundColor Yellow
    Write-Host "Not overwriting automatically (it may have machine-specific keys)." -ForegroundColor Yellow
    Write-Host "Diff it against $templatePath and merge the fields you want (enabledPlugins, extraKnownMarketplaces, model, language, effortLevel, theme, statusLine)." -ForegroundColor Yellow
} else {
    Copy-Item $templatePath $settingsPath
    Write-Host "Wrote $settingsPath from the template." -ForegroundColor Green
}

Write-Host ""
Write-Host "== Done. Manual steps remaining: ==" -ForegroundColor Cyan
Write-Host "1. Run 'claude' and log in (this provisions ~/.claude/.credentials.json - never scripted/shared)."
Write-Host "2. On first run, Claude Code will fetch the marketplaces/plugins listed in settings.json automatically."
Write-Host "3. If you use API-key billing instead of a Claude subscription login, set ANTHROPIC_API_KEY yourself"
Write-Host "   (or your own apiKeyHelper in settings.json) - do not copy another machine's key."
