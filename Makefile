# headroom-eval-cli — build matrix

GOOS ?= linux
GOARCH ?= amd64

# Standard Go optimizations
build-go:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="-s -w" \
		-trimpath \
		-o headroom-eval

# With UPX compression
build-upx: build-go
	upx --best headroom-eval -o headroom-eval-upx

# With mold linker + mimalloc (requires CGO + mold installed)
build-mold:
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="-s -w -linkmode external -extldflags '-fuse-ld=mold'" \
		-tags=mimalloc \
		-trimpath \
		-o headroom-eval-mold

# All sizes
sizes:
	@echo "go:     $$(wc -c < headroom-eval 2>/dev/null || echo 0)"
	@echo "upx:    $$(wc -c < headroom-eval-upx 2>/dev/null || echo 0)"
	@echo "mold:   $$(wc -c < headroom-eval-mold 2>/dev/null || echo 0)"

# Default — smallest
all: build-upx sizes
