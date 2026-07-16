#!/usr/bin/env bash
# Rebuilds frontend/tailwind.css from src/input.css, scanning index.html and
# js/*.js for used utility classes. Run this after adding/removing Tailwind
# classes in the markup or JS templates.
#
# Uses the Tailwind standalone CLI (a single binary, no Node/npm/CDN needed
# at runtime — only this one-time download to build). Internet is required
# here, at dev-time, the same way `go mod download` or `apk add` need it;
# the built tailwind.css is committed and served locally, so the running
# app itself never talks to the internet.
set -euo pipefail
cd "$(dirname "$0")/.."

CLI=./.tailwindcss-cli
if [ ! -x "$CLI" ]; then
  os="linux"; arch="x64"
  case "$(uname -s)" in
    Darwin) os="macos" ;;
    MINGW*|MSYS*|CYGWIN*) os="windows" ;;
  esac
  case "$(uname -m)" in
    arm64|aarch64) arch="arm64" ;;
  esac
  ext=""; [ "$os" = "windows" ] && ext=".exe"
  url="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-${os}-${arch}${ext}"
  echo "Downloading Tailwind CLI from $url"
  curl -sL -o "$CLI$ext" "$url"
  chmod +x "$CLI$ext" 2>/dev/null || true
  CLI="$CLI$ext"
fi

"$CLI" -i ./src/input.css -o ./tailwind.css --minify
echo "Wrote frontend/tailwind.css"
