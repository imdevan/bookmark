#!/usr/bin/env bash
# Update Homebrew formula with new version
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_TOML="${ROOT_DIR}/internal/package/package.toml"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [[ -z "${VERSION}" ]]; then
	echo "Usage: $0 VERSION"
	echo "Example: $0 0.2.0"
	exit 1
fi

# Remove 'v' prefix if present
VERSION="${VERSION#v}"

# Read package metadata
NAME="$(parse_toml_key "${PACKAGE_TOML}" "name")"
REPO_URL="$(parse_toml_key "${PACKAGE_TOML}" "repository")"
DESCRIPTION="$(parse_toml_key "${PACKAGE_TOML}" "description")"
HOMEPAGE="$(parse_toml_key "${PACKAGE_TOML}" "homepage")"

TAP_DIR="${ROOT_DIR}/../homebrew-${NAME}"
FORMULA_PATH="${TAP_DIR}/Formula/${NAME}.rb"

if [[ ! -d "${TAP_DIR}" ]]; then
	echo "❌ Homebrew tap not found at: ${TAP_DIR}"
	echo "Run 'just init-homebrew-tap' first"
	exit 1
fi

# Download tarball and calculate SHA256
TARBALL_URL="${REPO_URL}archive/refs/tags/v${VERSION}.tar.gz"
echo "📥 Downloading release tarball..."

if ! SHA256=$(download_and_hash "${TARBALL_URL}"); then
	echo "❌ Failed to download: ${TARBALL_URL}"
	exit 1
fi

echo "✅ SHA256: ${SHA256}"

# Update formula
CLASS_NAME="$(echo "${NAME}" | sed 's/.*/\u&/')"

cat >"${FORMULA_PATH}" <<EOF
class ${CLASS_NAME} < Formula
  desc "${DESCRIPTION}"
  homepage "${HOMEPAGE}"
  url "${TARBALL_URL}"
  sha256 "${SHA256}"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/${NAME}"
  end

  test do
    assert_match "v${VERSION}", shell_output("#{bin}/${NAME} --version")
  end
end
EOF

echo "✅ Updated formula: ${FORMULA_PATH}"
echo ""
echo "Next steps:"
echo "1. Test the formula locally:"
echo "   brew install --build-from-source ${FORMULA_PATH}"
echo "2. Commit and push:"
echo "   cd ${TAP_DIR}"
echo "   git add Formula/${NAME}.rb"
echo "   git commit -m \"Update ${NAME} to v${VERSION}\""
echo "   git push"
