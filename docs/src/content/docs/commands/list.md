---
title: list
description: List all bookmarks
---

List all bookmarks

## Usage

```bash
bookmark list
```

## Description

newListCmd creates the list command for displaying all bookmarks.

The list command shows all bookmarks in a formatted table with:
  - Alias: The bookmark name
  - Path: The directory path
  - Description: Optional bookmark description

The output is formatted with proper alignment for easy reading.

Examples:

	# List all bookmarks
	bookmark list

	# Use with custom config
	bookmark list -c ~/.config/bookmark/custom.toml

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `-c, --config` | string | config file path |

## Source

See [list_cmd.go](https://github.com/imdevan/bookmark//blob/main/cmd/bookmark/list_cmd.go) for implementation details.
