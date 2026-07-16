<#
.SYNOPSIS
  Run this before switching to the other PC in a synced Claude Code setup
  (see SYNC.md). Warns if Claude Code is still running here, and reminds
  you to confirm Syncthing shows "Up to Date" before switching.
#>

$running = Get-Process -Name "claude" -ErrorAction SilentlyContinue
if ($running) {
    Write-Host "⚠ claude.exe is still running on this PC (PID(s): $($running.Id -join ', '))." -ForegroundColor Red
    Write-Host "  Close it first, so its state finishes writing to disk before Syncthing syncs." -ForegroundColor Red
} else {
    Write-Host "✓ claude.exe is not running on this PC." -ForegroundColor Green
}

Write-Host ""
Write-Host "Before switching to the other PC:" -ForegroundColor Cyan
Write-Host "  1. Confirm claude.exe is closed here (see above)."
Write-Host "  2. Open Syncthing's web UI (default http://localhost:8384) and confirm this"
Write-Host "     folder shows 'Up to Date' on BOTH devices, not 'Syncing' or 'Out of Sync'."
Write-Host "  3. Only then start Claude Code on the other PC."
