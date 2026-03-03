---
title: completion
description: Generate shell completion scripts
---

Generate shell completion scripts

## Usage

```bash
bookmark completion [bash|zsh|fish|powershell]
```

## Description

newCompletionCmd creates the completion command for generating shell completion scripts.

The completion command generates shell completion scripts for various shells.
This enables tab-completion for bookmark commands and aliases.

Supported shells:
  - bash
  - zsh
  - fish
  - powershell

Examples:

	# Generate bash completion
	bookmark completion bash > /etc/bash_completion.d/bookmark

	# Generate zsh completion
	bookmark completion zsh > ~/.zsh/completion/_bookmark

	# Generate fish completion
	bookmark completion fish > ~/.config/fish/completions/bookmark.fish

	# Generate powershell completion
	bookmark completion powershell > bookmark.ps1

## Source

See [completion.go](https://github.com/imdevan/bookmark//blob/main/cmd/bookmark/completion.go) for implementation details.
