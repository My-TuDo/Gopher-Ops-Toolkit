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

## 未来规划

### exec 命令
- 远程命令执行功能待实现

### log 命令
- 日志分析功能待实现