---
title: config init
description: Generate a default config file
---

Generate a default config file

## Usage

```bash
bookmark init
```

## Description

newConfigInitCmd creates the config init command for generating a default config file.

The config init command creates a new configuration file with default values
at the standard XDG config location ($XDG_CONFIG_HOME/bookmark/config.toml).

Flags:
  - --force/-f: Overwrite existing config file
  - --editor/-e: Open the config file in your editor after creation

The generated config file includes commented examples for all available options.

Examples:

	# Generate default config
	bookmark config init

	# Overwrite existing config
	bookmark config init --force

	# Generate and open in editor
	bookmark config init --editor

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `-f, --force` | bool | overwrite existing config |
| `-e, --editor` | bool | open config in editor after creation |

## Source

See [config_init.go](https://github.com/imdevan/bookmark//blob/main/cmd/bookmark/config_init.go) for implementation details.
