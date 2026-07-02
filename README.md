# Gopher-Ops-Toolkit (gopt)

[![asciicast](https://asciinema.org/a/4GIzPOqI6gEDSRj7.svg)](https://asciinema.org/a/4GIzPOqI6gEDSRj7)

> 一个用 Go 编写的轻量级 SRE 网络健康探测 CLI 工具。  
> 支持 TCP/HTTP/DNS 多协议探测、并发探测引擎、XDG 标准配置管理、Docker 容器化部署。

---

## 特性

- **多协议探测** — TCP 端口探测（握手耗时分解 + Banner 读取）、HTTP 探测（自定义 Method/Header/关键字匹配）、DNS 探测（A/AAAA/MX/NS/CNAME/TXT）
- **并发探测引擎** — goroutine + channel 信号量限流，context 熔断机制，支持可配置轮次与平均延迟统计
- **四级配置链** — CLI 参数 > 环境变量 > XDG 标准路径 > 编译时内嵌默认配置（`//go:embed`）
- **灵活输出** — 表格/JSON 输出、日志文件保存（JSON Lines 格式）、窄终端自动回流行输出
- **容器化部署** — Docker 多阶段构建（golang → scratch），CGO_ENABLED=0 纯静态编译，最终镜像 <10MB

---

## 快速开始

### 安装

```bash
# 本地编译
go build -o gopt .

# 或 Docker 构建
docker build -t gopt .
```

### 一键探测

```bash
./gopt health --all
```

首次运行自动创建 `~/.config/gopt/config.yaml`，编辑即可自定义探测目标。

---

## 使用指南

### 全部探测

```bash
gopt health --all
```

从配置文件读取所有目标，逐一探测并汇总输出。

### TCP 端口探测

```bash
gopt health --type tcp --host example.com --port 3306
```

自动进行 DNS 解析耗时 + TCP 握手耗时分解，并尝试读取服务 Banner。

### HTTP 服务探测

```bash
# 基础用法
gopt health --type http --url http://example.com/health

# 自定义方法和请求头
gopt health --type http --url http://example.com/api \
    --method POST -H "Content-Type: application/json"

# 关键字检测
gopt health --type http --url http://example.com/health \
    --keyword "healthy"
```

### DNS 解析探测

```bash
# A 记录（默认）
gopt health --type dns --domain baidu.com

# 指定记录类型
gopt health --type dns --domain google.com --record-type MX
gopt health --type dns --domain google.com --record-type NS

# 使用指定 DNS 服务器
gopt health --type dns --domain baidu.com --dns-server 8.8.8.8
```

### 并发探测

```bash
# 5 轮并发取平均延迟（信号量限流 + 自动熔断）
gopt health --type tcp --host example.com --port 80 --rounds 5
```

### 输出与保存

```bash
gopt health --all --output json             # JSON 格式输出
gopt health --all --save                    # 保存结果到日志文件
gopt health --all --save --log-dir /var/log/ops  # 指定日志目录
```

---

## 配置

### 配置优先级（高 → 低）

```
① --config 显式指定
② GOPT_CONFIG 环境变量
③ ~/.config/gopt/config.yaml（XDG 标准路径）
④ configs/config.yaml（本地开发路径）
⑤ 编译时嵌入的默认配置（兜底）
```

### 配置文件示例

```yaml
output: table
save: false
log-dir: ./logs

health:
    timeout: 5s
    rounds: 1
    targets:
        tcp:  "localhost:8080"
        http: "http://localhost:8080"
        dns:  "baidu.com"
```

环境变量覆盖（`GOPT_` 前缀）：

```bash
GOPT_OUTPUT=json GOPT_LOG_DIR=/var/log/gopt gopt health --all
```

---

## Docker

```bash
# 构建
docker build -t gopt .

# 运行（配置已嵌入二进制，开箱即用）
docker run --rm gopt health --all

# 环境变量覆盖
docker run --rm \
  -e GOPT_OUTPUT=json \
  -e GOPT_TIMEOUT=2s \
  gopt health --type tcp --host example.com --port 80

# 挂载自定义配置
docker run --rm \
  -v /path/to/prod.yaml:/root/.config/gopt/config.yaml \
  gopt health --all
```

镜像基于多阶段构建（`golang:1.26-alpine` → `scratch`），最终大小约 **10MB**。

---

## 项目结构

```
├── main.go              # 入口
├── cmd/                 # Cobra 命令
│   ├── root.go          # 根命令 + 配置初始化
│   └── health.go        # 健康检查命令
├── internal/
│   ├── config/          # 配置管理（XDG + //go:embed）
│   ├── output/          # 输出格式化（表格/JSON/文件）
│   └── prober/          # 探测引擎
│       ├── prober.go    # Prober 接口 + MultiRoundProbe（并发 + 熔断）
│       ├── result.go    # Result 结构体 + 错误码
│       ├── tcp.go       # TCP 探针
│       ├── http.go      # HTTP 探针
│       └── dns.go       # DNS 探针
├── configs/             # 开发配置文件
└── Dockerfile           # 多阶段构建
```

---

## 技术栈

| 模块 | 技术 |
|------|------|
| 命令行框架 | Cobra + Viper |
| 并发控制 | goroutine + channel + context |
| 配置管理 | XDG 标准路径 + `//go:embed` |
| 输出格式 | tablewriter + JSON Lines |
| 容器化 | Docker 多阶段构建（scratch） |
| CI | GitHub Actions（lint + test + build） |

---

## 测试

```bash
# 运行全部测试
go test ./... -short -timeout 30s
```

| 包 | 状态 |
|------|------|
| `internal/prober` | ✅ 通过 |
| `internal/executor` | ✅ 通过 |
| `tests` | ✅ 通过 |

### 基准测试

```
BenchmarkMultiRoundProbe_Serial-12               5.2ms/op    0 allocs/op
BenchmarkMultiRoundProbe_Concurrent10-12         10.5ms/op   44 allocs/op
BenchmarkMultiRoundProbe_Concurrent50-12         52.4ms/op   260 allocs/op
BenchmarkMultiRoundProbe_ConcurrentClients-12    0.45ms/op   18 allocs/op
```

> `ConcurrentClients` 模拟多客户端并发请求，验证全局限流机制有效。

## 开发

```bash
# 安装 Git hooks（阻止误推送）
./scripts/setup-hooks.sh
```

