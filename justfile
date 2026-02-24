set shell := ["zsh", "-cu"]

build:
	go build -o bin/bookmark ./cmd/bookmark
	@size=$(stat -c %s bin/go-cli-template 2>/dev/null || stat -f %z bin/go-cli-template 2>/dev/null); \
	echo "Build size: $(awk "BEGIN {printf \"%.2f MB\", $size/1048576}")"

build-run:
	go build -o bin/bookmark ./cmd/bookmark && ./bin/bookmark

watch:
	@rg --files | entr -r sh -c 'sleep 0.5; go build -o bin/bookmark ./cmd/bookmark'

dev-build:
	go build -gcflags "all=-N -l" -o bin/bookmark ./cmd/bookmark

cross-platform:
	./scripts/build.sh

build-aur:
	./scripts/build_aur.sh

install:
	install -m 0755 bin/bookmark /usr/local/bin/bookmark

test:
	go test ./...

test-verbose:
	go test -v ./...

sync:
	./scripts/sync.sh

clean:
	rm -rf bin

# Documentation tasks
docs-init:
	@echo "📦 Installing documentation dependencies..."
	cd docs && bun install

docs-generate:
	@echo "📝 Generating API documentation from Go packages..."
	./scripts/docs_generate.sh

docs-dev:
	@echo "🚀 Starting documentation development server..."
	@just docs-generate
	cd docs && bun run dev

docs-build:
	@echo "🏗️  Building documentation site..."
	@just docs-generate
	cd docs && NODE_ENV=production bun run build

docs-preview:
	@echo "👀 Previewing built documentation..."
	cd docs && bun run preview

docs-clean:
	@echo "🧹 Cleaning documentation build artifacts..."
	rm -rf docs/dist docs/.astro docs/node_modules docs/src/content/docs/api
