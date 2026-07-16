#!/usr/bin/env bash
# Run this before switching to the other PC in a synced Claude Code setup
# (see SYNC.md). Warns if Claude Code is still running here, and reminds
# you to confirm Syncthing shows "Up to Date" before switching.
set -uo pipefail

if pgrep -x "claude" >/dev/null 2>&1; then
  echo "⚠ claude is still running on this machine (PID(s): $(pgrep -x claude | tr '\n' ' '))."
  echo "  Close it first, so its state finishes writing to disk before Syncthing syncs."
else
  echo "✓ claude is not running on this machine."
fi

echo ""
echo "Before switching to the other machine:"
echo "  1. Confirm claude is closed here (see above)."
echo "  2. Open Syncthing's web UI (default http://localhost:8384) and confirm this"
echo "     folder shows 'Up to Date' on BOTH devices, not 'Syncing' or 'Out of Sync'."
echo "  3. Only then start Claude Code on the other machine."
