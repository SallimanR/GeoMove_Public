#!/bin/bash
# Installs git hooks from scripts/ into .git/hooks/
set -e

HOOKS_DIR="$(cd "$(dirname "$0")" && pwd)"
GIT_HOOKS_DIR="$(git rev-parse --git-common-dir)/hooks"

for script in "$HOOKS_DIR"/prepare-commit-msg; do
	name=$(basename "$script")
	cp "$script" "$GIT_HOOKS_DIR/$name"
	chmod +x "$GIT_HOOKS_DIR/$name"
	echo "[ok] installed $name"
done
