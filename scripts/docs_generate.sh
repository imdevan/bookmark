#!/usr/bin/env bash
set -euo pipefail

# Generate API documentation from Go packages using gomarkdoc
# Usage: ./docs_generate.sh [--dev]
#   --dev: Use '/' as base for local development

PACKAGE_FILE="internal/package/package.toml"
DOCS_API_DIR="docs/src/content/docs/api"
DOCS_CONFIG="docs/config.mjs"
DOCS_SIDEBAR="docs/sidebar.mjs"
CMD_DIR="cmd/bookmark"

# Parse package.toml
parse_toml() {
  local key=$1
  grep "^$key = " "$PACKAGE_FILE" | sed 's/^[^=]*= *"\(.*\)"$/\1/'
}

echo "📦 Reading package metadata..."
PROJECT_NAME=$(parse_toml "name")
DESCRIPTION=$(parse_toml "description")
DOCS_SITE=$(parse_toml "docs_site")
DOCS_BASE=$(parse_toml "docs_base")
REPOSITORY=$(parse_toml "repository")

# Use defaults if repository is empty
if [ -z "$REPOSITORY" ]; then
  REPOSITORY="https://github.com/yourusername/${PROJECT_NAME}"
fi

echo "🔧 Updating docs config..."

# Update docs/config.mjs with values from package.toml
if [ -f "$DOCS_CONFIG" ]; then
  cat >"$DOCS_CONFIG" <<EOF
const stage = process.env.NODE_ENV || "dev"
const isProduction = stage === "production"

export default {
  url: isProduction ? "$DOCS_SITE" : "http://localhost:4321",
  basePath:  isProduction ? "$DOCS_BASE" : "/",
  github: "$REPOSITORY",
  githubDocs: "$REPOSITORY",
  title: "$PROJECT_NAME",
  description: "$DESCRIPTION",
}
EOF
  echo "  ✓ Updated config.mjs with package metadata"
fi

echo "🔧 Generating sidebar configuration..."

# Detect commands from cmd directory
COMMANDS=""
if [ -d "$CMD_DIR" ]; then
  for cmd_file in "$CMD_DIR"/*.go; do
    # Skip test files, main.go, and root.go
    if [[ "$cmd_file" == *"_test.go" ]] || [[ "$cmd_file" == *"/main.go" ]] || [[ "$cmd_file" == *"/root.go" ]]; then
      continue
    fi
    
    # Extract command name from filename (e.g., config.go -> config, delete_cmd.go -> delete)
    cmd_name=$(basename "$cmd_file" .go | sed 's/_cmd$//')
    
    # Convert underscores to spaces for display (e.g., config_init -> config init)
    cmd_display=$(echo "$cmd_name" | sed 's/_/ /g')
    
    # Convert underscores to hyphens for URL (e.g., config_init -> config-init)
    cmd_url=$(echo "$cmd_name" | sed 's/_/-/g')
    
    COMMANDS="${COMMANDS}            { label: '${cmd_display}', link: '/commands/${cmd_url}' },
"
  done
fi

# Generate sidebar.mjs
cat >"$DOCS_SIDEBAR" <<EOF
export default [
  {
    label: '${PROJECT_NAME}',
    link: '/',
  },
  {
    label: 'Install',
    link: '/install',
  },
  {
    label: 'Commands',
    items: [
      { label: '${PROJECT_NAME}', link: '/commands/${PROJECT_NAME}' },
${COMMANDS}    ],
  },
  {
    label: 'Configuration',
    link: '/configuration',
  },
]
EOF

echo "  ✓ Generated sidebar.mjs with detected commands"

echo "📝 Generating content pages..."

DOCS_CONTENT_DIR="docs/src/content/docs"

# Generate index page from README.md
if [ -f "README.md" ]; then
  {
    echo "---"
    echo "title: ${PROJECT_NAME}"
    echo "description: ${DESCRIPTION}"
    echo "---"
    echo ""
    # Skip the first heading from README and output the rest
    sed '1{/^# /d;}' README.md
  } >"${DOCS_CONTENT_DIR}/index.md"
  echo "  ✓ Generated index.md from README.md"
fi

# Generate install page from INSTALL.md
if [ -f "INSTALL.md" ]; then
  {
    echo "---"
    echo "title: Install"
    echo "description: Installation instructions for ${PROJECT_NAME}"
    echo "---"
    echo ""
    # Skip the first heading from INSTALL.md and output the rest
    sed '1{/^# /d;}' INSTALL.md
  } >"${DOCS_CONTENT_DIR}/install.md"
  echo "  ✓ Generated install.md from INSTALL.md"
fi

# Generate configuration page placeholder
cat >"${DOCS_CONTENT_DIR}/configuration.md" <<EOF
---
title: Configuration
description: Configuration options for ${PROJECT_NAME}
---

Configuration file location: \`\$XDG_CONFIG_HOME/${PROJECT_NAME}/config.toml\`

See the [config API documentation](/api/config) for available configuration options.
EOF

echo "  ✓ Generated configuration.md"

# Create commands directory
mkdir -p "${DOCS_CONTENT_DIR}/commands"

# Generate root command page from root.go
cat >"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md" <<EOF
---
title: ${PROJECT_NAME}
description: Root command for ${PROJECT_NAME}
---

The root command for ${PROJECT_NAME}.

See [root.go](${REPOSITORY}/blob/main/cmd/bookmark/root.go) for implementation details.
EOF

echo "  ✓ Generated commands/${PROJECT_NAME}.md"

echo "🔧 Checking for gomarkdoc..."
if ! command -v gomarkdoc &>/dev/null; then
  echo "📦 Installing gomarkdoc..."
  go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
fi

echo "🧹 Cleaning old API docs..."
rm -rf "$DOCS_API_DIR"
mkdir -p "$DOCS_API_DIR"

echo "📝 Generating API documentation..."

# Generate docs for each internal package
for pkg in internal/*/; do
  pkg_name=$(basename "$pkg")

  # Skip test utilities and adapters subdirectories
  if [[ "$pkg_name" == "testutil" ]]; then
    continue
  fi

  echo "  - Processing $pkg_name..."

  # Generate to temp file first
  gomarkdoc --output "/tmp/${pkg_name}.md" "./$pkg" 2>/dev/null || {
    echo "    ⚠️  No exported symbols in $pkg_name"
    continue
  }

  # Add frontmatter and content (skip HTML comment and blank lines)
  {
    echo "---"
    echo "title: ${pkg_name}"
    echo "description: API documentation for the ${pkg_name} package"
    echo "---"
    echo ""
    # Skip HTML comment and any frontmatter that gomarkdoc added
    sed -n '/^# /,$p' "/tmp/${pkg_name}.md"
  } >"$DOCS_API_DIR/${pkg_name}.md"
done

# Generate docs for adapters
echo "  - Processing adapters..."
mkdir -p "$DOCS_API_DIR/adapters"

for adapter in internal/adapters/*/; do
  adapter_name=$(basename "$adapter")
  echo "    - Processing adapters/$adapter_name..."

  # Generate to temp file first
  gomarkdoc --output "/tmp/adapter_${adapter_name}.md" "./$adapter" 2>/dev/null || {
    echo "      ⚠️  No exported symbols in $adapter_name"
    continue
  }

  # Add frontmatter and content (skip HTML comment and blank lines)
  {
    echo "---"
    echo "title: adapters/${adapter_name}"
    echo "description: API documentation for the ${adapter_name} adapter"
    echo "---"
    echo ""
    # Skip HTML comment and any frontmatter that gomarkdoc added
    sed -n '/^# /,$p' "/tmp/adapter_${adapter_name}.md"
  } >"$DOCS_API_DIR/adapters/${adapter_name}.md"
done

echo "✅ API documentation generated successfully!"
echo "📁 Output: $DOCS_API_DIR"
