---
title: config
description: View or edit configuration
---

View or edit configuration

## Usage

```bash
bookmark config
```

## Description

newConfigCmd creates the config command for viewing or editing configuration.

The config command opens the configuration file in your configured editor.
If the config file doesn't exist, it will be created with default values.

Subcommands:
  - init: Generate a default config file

Examples:

	# Open config in editor
	bookmark config

	# Initialize new config file
	bookmark config init

	# Force overwrite existing config
	bookmark config init --force

## Source

See [config.go](https://github.com/imdevan/bookmark//blob/main/cmd/bookmark/config.go) for implementation details.
