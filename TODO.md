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
- [ ] 完善 TCP 探测（握手耗时分解 + Banner 读取）
- [ ] 新增更多探测指标（Ping 等）

## 计划中

### exec 命令
- 远程命令执行功能待实现

### log 命令
- 日志分析功能待实现