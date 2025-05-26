#!/bin/bash

WriteMessage() {
  echo "$1" >&2
}

COMMIT_MSG=$(head -n 1 $1)
PATTERN="^(feat|fix|refacto|docs|test|chore|build|ci|perf|refactor|style)(\([a-Z0-9\-]*\))?: .{1,}$"

if ! echo "$COMMIT_MSG" | grep -Eq "$PATTERN"; then
  WriteMessage "Bad commit format"
  WriteMessage "<type>: [subject]"
  exit 1
fi
exit 0
