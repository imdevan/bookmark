---
title: delete
description: Delete a bookmark
---

Delete a bookmark

## Usage

```bash
bookmark delete <alias>
```

## Description

newDeleteCmd creates the delete command for removing bookmarks.

The delete command removes a bookmark by its alias.
By default, it will prompt for confirmation before deleting.

Flags:
  - --force/-f: Skip confirmation prompt and delete immediately

Examples:

	# Delete with confirmation
	bookmark delete myproject

	# Force delete without confirmation
	bookmark delete myproject --force

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `-c, --config` | string | config file path |
| `-f, --force` | bool | skip confirmation |

## Source

See [delete_cmd.go](https://github.com/imdevan/bookmark//blob/main/cmd/bookmark/delete_cmd.go) for implementation details.
