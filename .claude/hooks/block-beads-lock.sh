#!/usr/bin/env bash
# Block edits to beads runtime lock files
FILE="$CLAUDE_TOOL_INPUT_FILE_PATH"
if echo "$FILE" | grep -q '\.beads/.*\.lock$'; then
  echo "Blocked: do not edit beads lock files ($FILE)" >&2
  exit 2
fi
exit 0
