#!/usr/bin/env bash
# Initialize Homebrew tap repository
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_TOML="${ROOT_DIR}/internal/package/package.toml"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

# Read package metadata
NAME="$(parse_toml_key "${PACKAGE_TOML}" "name")"
DESCRIPTION="$(parse_toml_key "${PACKAGE_TOML}" "description")"
HOMEPAGE="$(parse_toml_key "${PACKAGE_TOML}" "homepage")"
REPO_URL="$(parse_toml_key "${PACKAGE_TOML}" "repository")"

# Extract GitHub username from repository URL
GITHUB_USER="$(echo "${REPO_URL}" | sed -E 's|https://github.com/([^/]+)/.*|\1|')"

TAP_NAME="homebrew-${NAME}"
TAP_DIR="${ROOT_DIR}/../${TAP_NAME}"

echo "🍺 Initializing Homebrew tap repository..."
echo "   Tap name: ${TAP_NAME}"
echo "   Location: ${TAP_DIR}"

# Create tap directory if it doesn't exist
if [[ -d "${TAP_DIR}" ]]; then
  echo "⚠️  Tap directory already exists: ${TAP_DIR}"
  read -p "Do you want to reinitialize it? (y/N) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
  fi
  rm -rf "${TAP_DIR}"
fi

mkdir -p "${TAP_DIR}/Formula"

# Initialize git repository
cd "${TAP_DIR}"
git init
git branch -M main

# Create README
cat >"${TAP_DIR}/README.md" <<EOF
# Homebrew Tap for ${NAME}

This is the official Homebrew tap for [${NAME}](${HOMEPAGE}).

## Installation

\`\`\`bash
brew tap ${GITHUB_USER}/${NAME}
brew install ${NAME}
\`\`\`

## Updating

\`\`\`bash
brew update
brew upgrade ${NAME}
\`\`\`

## Uninstall

\`\`\`bash
brew uninstall ${NAME}
brew untap ${GITHUB_USER}/${NAME}
\`\`\`
EOF

# Create initial formula template
cat >"${TAP_DIR}/Formula/${NAME}.rb" <<EOF
class $(echo "${NAME}" | sed 's/.*/\u&/') < Formula
  desc "${DESCRIPTION}"
  homepage "${HOMEPAGE}"
  url "${REPO_URL}archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/${NAME}"
  end

  test do
    assert_match "v0.1.0", shell_output("#{bin}/${NAME} --version")
  end
end
EOF

# Create .gitignore
cat >"${TAP_DIR}/.gitignore" <<EOF
.DS_Store
*.swp
*.swo
*~
EOF

# Initial commit
git add .
git commit -m "Initial commit: Homebrew tap for ${NAME}"

echo ""
echo "✅ Homebrew tap initialized at: ${TAP_DIR}"
echo ""
echo "Next steps:"
echo "1. Create a GitHub repository: https://github.com/new"
echo "   Repository name: ${TAP_NAME}"
echo "2. Push the tap:"
echo "   cd ${TAP_DIR}"
echo "   git remote add origin git@github.com:${GITHUB_USER}/${TAP_NAME}.git"
echo "   git push -u origin main"
echo "3. Update the formula with actual release SHA256 using:"
echo "   just update-homebrew-formula VERSION"
