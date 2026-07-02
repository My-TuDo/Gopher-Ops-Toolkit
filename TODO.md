# TODO

## 2026-06

### health 命令重构
- [x] 将默认探测项目由 `tcp` 改为 `all`
- [x] 新增 `--all` 参数，执行全部探测
- [x] 细化参数传递逻辑：
  - `--host` / `--port` — 仅在 `--type tcp` 时可用
  - `--url` / `--method` / `--header` — 仅在 `--type http` 时可用
  - `--domain` / `--record-type` — 仅在 `--type dns` 时可用
- [x] 删除旧的 `--target` 参数
- [x] 完善 TCP 探测（握手耗时分解 + Banner 读取）
- [x] 完善探测功能 (并发探测，取耗时平均数)

---

## 2026-07 Plan

### P0 — 核心功能补全（优先）

- [x] **HTTP 探针增强**
  - [x] 使用 `--method` 参数，支持 POST/PUT/HEAD 等
  - [x] 使用 `--header` 参数，支持自定义请求头
  - [x] 增加响应体检测（根据关键字判断健康）
  - [x] 记录 HTTP 状态码到 Detail（如 `200 OK`）
  - [x] 记录响应大小

- [x] **health.go 重构：switch 职责分离**
  - 将 switch 缩减为：**参数校验 + 返回 (probeInstance, target) 二元组**
  - 不再在 switch 内部执行 MultiRoundProbe 和输出
  - switch 之后只保留一行通用调用：
    ```go
    res := prober.MultiRoundProbe(probeInstance, target, timeout, rounds)
    fmt.Println(res.String())
    ```
  - 提取 `buildProber()` 工厂函数（可选，switch 不挪出去也能达到分离效果）

- [x] **DNS 探针增强**
  - 给 `DNSProber` 加字段（`RecordType`, `DNSServer`），类似 HTTP 的 struct 方案
  - 根据 `--record-type` 查询不同记录类型：
    - `A` → `LookupHost`（已有）
    - `MX` → `LookupMX`
    - `NS` → `LookupNS`
    - `CNAME` → `LookupCNAME`
    - `TXT` → `LookupTXT`
  - 用 `--dns-server` 替代现有的 `@` 语法，消除分隔符冲突

- [x] **修复 MultiRoundProbe 重复调用问题**
  - 当前 `firstRes := p.Probe(target, timeout)` 额外多跑了一次探测
  - 应该用第一次成功的结果或者传 Name 进去

### P1 — 质量提升

- [ ] **Result.String() 增强**
  - 显示 Latency（响应时间）
  - 显示 Detail 信息
  - 示例：`TCP探测 | 目标: localhost:8080 | 状态: 健康 | 延迟: 12ms | 端口开放 (DNS: 0ms, 握手: 2ms)`

- [x] **代码清理**
  - 修复 `reusltCh` 拼写为 `resultCh`
  - 修复 `timeout` 描述末尾多余的反引号

### P2 — 新功能探索

- [ ] **Ping 探针**（TCP Connect 模拟）
  - 作为前置连通性检查
  - all 模式下自动从 tcp target 推导主机

- [x] **output 功能完善**
  - [x] 表格输出（默认表格展示）
  - [x] 支持 `--output json` 输出 JSON 格式到终端
  - [x] 支持 `--save` 将结果保存到日志文件（JSON Lines 格式，按日期命名）
  - [x] 日志文件目录可配置（默认 `./logs/`）
  - [x] 窄终端自动回退行输出

- [ ] **健康阈值可配置**
  - 允许用户设置延迟阈值（如 `--max-latency 500ms`）
  - 超阈值即使通也标记为"不健康"

---

## 2026-07 面试复盘优化项

以下优化点来自大厂面试技术拷打复盘，按优先级排列：

### P0 — 必须修复（影响正确性 & 健壮性）✅

- [x] **tcp.go: DNS 多 IP Fallback**
  - 遍历 ips 切片，逐 IP 重试 DialContext，成功后 break
  - 每次重试前检查 `ctx.Err()`，超时后停止重试

- [x] **result.go: MultiRoundProbe 并发限流 + 熔断**
  - 信号量 channel（`maxConcurrentProbes=5`）控制并发 goroutine 上限
  - 叠加熔断：atomic 计数失败次数，达到阈值时 `context.CancelFunc` 通知退出
  - `select` 多路复用令牌获取与 `<-ctx.Done()` 实现优雅取消

- [x] **http.go: HTTP Client 复用 & 内存安全**
  - 包级 `defaultHTTPClient` 复用 Transport 连接池
  - `io.LimitReader(resp.Body, 10MB)` 限制 body 读取，防止大响应体 OOM
  - `defer` 中 `io.Copy(io.Discard, resp.Body)` 耗尽 body 后归还连接到池

- [x] **config.go: os.Stat 权限判断 BUG**
  - 用 `os.IsNotExist(err)` 区分「不存在」和「其他错误」
  - 权限不足时 stderr 告警 + fallback 到下一级路径

- [x] **health.go & 所有 cmd: 干掉散落 os.Exit，统一用 RunE**
  - `healthCmd.Run` → `RunE`，返回 error 替代 `os.Exit(1)`
  - `Execute()` 顶层统一 `os.Exit(1)`，保持唯一退出点
  - 配套 `SilenceErrors: true` + `SilenceUsage: true`

### P1 — 建议优化（质量 & 体验）

- [x] **config.go: 空 error 处理块加输出**
  - `EnsureDefaultConfig` 失败时输出 `⚠️` 警告信息，用户有感知

- [ ] **Result 结构体增加 ErrorCode 字段，区分「不健康」和「错误」**
  - 当前 Status 只有「健康」「不健康」二值，无法区分：
    - 目标服务挂了（如 TCP 拒绝连接）→ UNREACHABLE
    - 工具自身不支持（如 DNS RecordType 不支持）→ UNSUPPORTED
    - 探测超时 → TIMEOUT
    - 系统错误（如权限不足）→ SYSTEM_ERROR
  - 新增 `ErrorCode string` 字段，调用方可基于 ErrorCode 做不同策略
  - 关联修改：所有 Probe() 返回 Result 的地方补充 ErrorCode

- [ ] **tcp.go: Banner 读取超时可配置**
  - 当前写死 `300ms`，对于慢启动协议（SSH/Redis）可能漏读
  - 优化：提取为结构体字段或 config 项

- [ ] **cmd/exec.go: 空壳命令**
  - 当前 `Run: func() { fmt.Println("exec called") }`，无实际功能
  - 要么实现真实功能，要么删除避免迷惑用户

### P2 — 远期规划

- [ ] **Dockerfile: 考虑 distroless/base 替代 scratch**
  - scratch 缺少 /etc/resolv.conf 的兜底，某些 Docker 网络插件下 DNS 解析异常
  - gcr.io/distroless/base 包含基础系统文件，仅比 scratch 大几 MB

- [ ] **全局 http.Client 连接池参数可配置**
  - MaxIdleConns、MaxIdleConnsPerHost、IdleConnTimeout、TLSHandshakeTimeout
  - 支持从 config.yaml 读取或 CLI 参数覆盖

### P0 — 部署适配（v0.2.0 优先）

- [x] **XDG 标准路径适配**
  - 默认配置路径: `~/.config/gopt/config.yaml`
  - 默认日志目录: `~/.local/share/gopt/logs/`
  - `--config` 可选覆盖，不传时自动读取默认路径
  - 首次运行时自动创建目录和默认配置文件

- [x] **Dockerfile**
  - 多阶段构建（builder → scratch）
  - 镜像体积控制在 20MB 以内 (实际: 12.3MB)
  - 默认 `ENTRYPOINT ["/gopt"]`
  - 配置已通过 //go:embed 嵌入二进制，Docker 开箱即用

- [x] **Makefile**
  - `make build` — 编译
  - `make test` — 跑测试（含覆盖率）
  - `make release` — 交叉编译(5平台) + SHA256 校验和
  - `make install` — 安装到 `/usr/local/bin`
  - `make lint` — go vet
  - `make docker` — Docker 构建
  - `make clean` — 清理产物

### 我的补充建议

- [x] **`version` 命令增强** — `-v` 显示编译时间、Go 版本、commit hash（通过 ldflags 注入）
- [ ] **日志自动轮转** — `--save` 写入的 `.jsonl` 文件超过一定大小自动分割
- [ ] **命令补全** — 支持 `gopt completion bash` / `zsh`，方便日常使用

---

## 未来规划

### exec 命令
- 远程命令执行功能待实现

### log 命令
- 日志分析功能待实现