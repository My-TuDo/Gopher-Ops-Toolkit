BINARY   = gopt
MODULE   = github.com/My-TuDo/Gopher-Ops-Toolkit
COMMIT   = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILT_AT = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  = -ldflags="-s -w \
	-X $(MODULE)/cmd.commitHash=$(COMMIT) \
	-X $(MODULE)/cmd.buildTime=$(BUILT_AT)"

DIST_DIR   = dist
GO        ?= go
CGO       ?= CGO_ENABLED=0

.PHONY: all build test lint clean install release docker

all: build

# ─── 编译 ──────────────────────────────────────────────
build:
	$(CGO) $(GO) build $(LDFLAGS) -o $(BINARY) .

build/linux/amd64:
	GOOS=linux GOARCH=amd64 $(CGO) $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY)-linux-amd64 .

build/linux/arm64:
	GOOS=linux GOARCH=arm64 $(CGO) $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY)-linux-arm64 .

build/darwin/amd64:
	GOOS=darwin GOARCH=amd64 $(CGO) $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY)-darwin-amd64 .

build/darwin/arm64:
	GOOS=darwin GOARCH=arm64 $(CGO) $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY)-darwin-arm64 .

build/windows/amd64:
	GOOS=windows GOARCH=amd64 $(CGO) $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY)-windows-amd64.exe .

# ─── 测试 ──────────────────────────────────────────────
test:
	$(GO) test ./... -v -count=1 -timeout 120s

test/short:
	$(GO) test ./... -short -count=1 -timeout 60s

test/cover:
	$(GO) test ./... -coverprofile=coverage.out -timeout 120s
	$(GO) tool cover -func=coverage.out | tail -1

# ─── 代码检查 ──────────────────────────────────────────
lint:
	go vet ./...

# ─── 清理 ──────────────────────────────────────────────
clean:
	rm -rf $(BINARY) $(DIST_DIR) coverage.out
	-docker rmi $(BINARY):local 2>/dev/null

# ─── 安装 ──────────────────────────────────────────────
install: build
	install -m 0755 $(BINARY) /usr/local/bin/$(BINARY)

# ─── 交叉编译 + 打包 ──────────────────────────────────
release: clean \
	build/linux/amd64 build/linux/arm64 \
	build/darwin/amd64 build/darwin/arm64 \
	build/windows/amd64
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  生成 SHA256 校验和"
	@echo "═══════════════════════════════════════════"
	cd $(DIST_DIR) && sha256sum $(BINARY)-* > checksums.txt
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  Release 产物清单"
	@echo "═══════════════════════════════════════════"
	ls -lh $(DIST_DIR)/*

# ─── Docker ────────────────────────────────────────────
docker:
	docker build -t $(BINARY):local .

.PHONY: build/linux/amd64 build/linux/arm64 \
	build/darwin/amd64 build/darwin/arm64 \
	build/windows/amd64
