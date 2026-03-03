#!/usr/bin/env bash
# Interactive bookmark navigation function
# Usage: source bm.sh, then run: bm

bm() {
  local cmd=$(CLICOLOR_FORCE=1 bookmark -i)

  # If a command was returned, execute it
  if [[ -n "$cmd" ]]; then
    eval "$cmd"
  fi
}

# Test it
echo "bm function loaded. Run 'bm' to use interactive bookmark navigation."
