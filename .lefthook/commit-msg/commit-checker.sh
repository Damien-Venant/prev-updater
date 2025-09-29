#!/bin/bash

WriteMessage() {
  echo "$1" >&2
}

COMMIT_MSG=$(head -n 1 $1)
PATTERN="(fix|feat|refacto|style|docs|chore|perf|test|build|ci)!?: .{1,}"

if ! echo "$COMMIT_MSG" | grep -Eq "$PATTERN"; then
  WriteMessage "Bad commit format"
  WriteMessage "<type>: [subject]"
  exit 1
fi
exit 0
