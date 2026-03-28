#!/usr/bin/env bash
# Auto-format Go files with gofumpt after edits
FILE="$CLAUDE_TOOL_INPUT_FILE_PATH"
case "$FILE" in
  *.go) gofumpt -w "$FILE" 2>/dev/null ;;
esac
exit 0
