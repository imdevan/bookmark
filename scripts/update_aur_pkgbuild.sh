#!/usr/bin/env bash
# Update AUR PKGBUILD with new version
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
AUTHOR="$(parse_toml_key "${PACKAGE_TOML}" "author")"

AUR_DIR="${ROOT_DIR}/../aur-${NAME}"
PKGBUILD_PATH="${AUR_DIR}/PKGBUILD"

if [[ ! -d "${AUR_DIR}" ]]; then
	echo "❌ AUR repository not found at: ${AUR_DIR}"
	echo "Run 'just init-aur-repo' first"
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

# Update PKGBUILD
cat >"${PKGBUILD_PATH}" <<EOF
# Maintainer: ${AUTHOR}
pkgname=${NAME}
pkgver=${VERSION}
pkgrel=1
pkgdesc="${DESCRIPTION}"
arch=('x86_64' 'aarch64')
url="${HOMEPAGE}"
license=('MIT')
depends=()
makedepends=('go')
source=("\${pkgname}-\${pkgver}.tar.gz::${TARBALL_URL}")
sha256sums=('${SHA256}')

build() {
  cd "\${pkgname}-\${pkgver}"
  export CGO_ENABLED=0
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  go build -ldflags="-s -w" -o \${pkgname} ./cmd/\${pkgname}
}

package() {
  cd "\${pkgname}-\${pkgver}"
  install -Dm755 \${pkgname} "\${pkgdir}/usr/bin/\${pkgname}"
  install -Dm644 LICENSE "\${pkgdir}/usr/share/licenses/\${pkgname}/LICENSE"
}
EOF

# Generate .SRCINFO
cd "${AUR_DIR}"
if command -v makepkg &>/dev/null; then
	makepkg --printsrcinfo >.SRCINFO
	echo "✅ Generated .SRCINFO"
else
	echo "⚠️  makepkg not found, skipping .SRCINFO generation"
	echo "   You'll need to run 'makepkg --printsrcinfo > .SRCINFO' manually"
fi

echo "✅ Updated PKGBUILD: ${PKGBUILD_PATH}"
echo ""
echo "Next steps:"
echo "1. Test the package locally:"
echo "   cd ${AUR_DIR}"
echo "   makepkg -si"
echo "2. Commit and push:"
echo "   git add PKGBUILD .SRCINFO"
echo "   git commit -m \"Update ${NAME} to v${VERSION}\""
echo "   git push"
