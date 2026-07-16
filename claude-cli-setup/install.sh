#!/usr/bin/env bash
# Reproduces this PC's Claude Code CLI environment (native binary + global
# settings: model, language, effort level, theme, enabled plugins/marketplaces)
# on another macOS/Linux machine.
#
# Deliberately does NOT touch API keys/OAuth tokens - those are per-machine
# secrets and must be set up by logging in on the new machine (see README.md).
set -euo pipefail

echo "== Claude Code CLI environment setup =="

# 1. Install the Claude Code CLI (native binary), if not already present.
if command -v claude >/dev/null 2>&1; then
  echo "claude CLI already installed at $(command -v claude) - skipping install."
else
  echo "Installing Claude Code CLI..."
  curl -fsSL https://claude.sh | bash
fi

# 2. Write global settings.json, without clobbering an existing one that
#    might have its own apiKeyHelper etc.
CLAUDE_DIR="${CLAUDE_CONFIG_DIR:-$HOME/.claude}"
mkdir -p "$CLAUDE_DIR"
SETTINGS_PATH="$CLAUDE_DIR/settings.json"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_PATH="$SCRIPT_DIR/settings.template.json"

if [ -f "$SETTINGS_PATH" ]; then
  BACKUP_PATH="$SETTINGS_PATH.bak.$(date +%Y%m%d-%H%M%S)"
  cp "$SETTINGS_PATH" "$BACKUP_PATH"
  echo "Existing settings.json found - backed up to $BACKUP_PATH"
  echo "Not overwriting automatically (it may have machine-specific keys)."
  echo "Diff it against $TEMPLATE_PATH and merge the fields you want (enabledPlugins, extraKnownMarketplaces, model, language, effortLevel, theme, statusLine)."
else
  cp "$TEMPLATE_PATH" "$SETTINGS_PATH"
  echo "Wrote $SETTINGS_PATH from the template."
fi

echo ""
echo "== Done. Manual steps remaining: =="
echo "1. Run 'claude' and log in (this provisions ~/.claude/.credentials.json - never scripted/shared)."
echo "2. On first run, Claude Code will fetch the marketplaces/plugins listed in settings.json automatically."
echo "3. If you use API-key billing instead of a Claude subscription login, set ANTHROPIC_API_KEY yourself"
echo "   (or your own apiKeyHelper in settings.json) - do not copy another machine's key."
