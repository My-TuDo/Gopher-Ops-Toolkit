# Gopher-Ops-Toolkit

一个用 Go 编写的命令行工具，提供各种运维常用的小功能。

---

## 快速开始

### 构建

```bash
go build -o gopt
```

### 配置

编辑 `configs/config.yaml`，或通过 `--config` 指定自定义路径：

```yaml
output: table        # 输出格式: table | json
save: false          # 是否保存结果到日志文件
log-dir: ./logs      # 日志文件保存目录

health:
    timeout: 5s
    rounds: 1
    targets:
        tcp:  "localhost:8080"
        http: "http://localhost:8080"
        dns:  "baidu.com"
```

---

## 使用指南

### health — 健康检查

#### 全部探测（从配置文件读取）

```bash
gopt health --all
```

#### TCP 端口探测

```bash
gopt health --type tcp --host example.com --port 3306
```

`TCP` 探针会自动进行 DNS 解析耗时 + 握手耗时分解，并尝试读取服务 Banner。

#### HTTP 服务探测

```bash
# 基础用法
gopt health --type http --url http://example.com/health

# 自定义请求方法和请求头
gopt health --type http --url http://example.com/api \
    --method POST -H "Content-Type: application/json"

# 响应体关键字检测
gopt health --type http --url http://example.com/health \
    --keyword "healthy"
```

#### DNS 解析探测

```bash
# 查询 A 记录（默认）
gopt health --type dns --domain baidu.com

# 指定记录类型
gopt health --type dns --domain google.com --record-type MX
gopt health --type dns --domain google.com --record-type NS
gopt health --type dns --domain google.com --record-type TXT

# 使用外置 DNS 服务器
gopt health --type dns --domain baidu.com --dns-server 8.8.8.8
```

#### 输出与保存

```bash
# JSON 格式输出
gopt health --all --output json

# 保存结果到日志文件
gopt health --all --save

# 指定日志目录
gopt health --all --save --log-dir /var/log/ops

# 多轮并发探测（取平均延迟）
gopt health --type tcp --host example.com --port 80 --rounds 5
```

### version — 版本信息

```bash
gopt version
```

---

## 项目结构

```
gopt/
├── cmd/                      # 命令定义
│   ├── health.go             # 健康检查命令（已完成）
│   ├── root.go               # 根命令 & 配置初始化
│   ├── version.go            # 版本信息
│   ├── exec.go               # 远程命令执行（待实现）
│   ├── log.go                # 日志分析（待实现）
│   └── monitor.go            # 批量执行（待实现）
├── configs/
│   └── config.yaml           # 配置文件
├── internal/
│   ├── output/               # 输出格式化（表格 / JSON / 文件保存）
│   ├── prober/               # 健康检查探针
│   │   ├── dns.go            # DNS 探测
│   │   ├── http.go           # HTTP 探测
│   │   ├── result.go         # 结果结构体 & 通用函数
│   │   ├── tcp.go            # TCP 端口探测
│   │   ├── *_test.go         # 单元测试
│   ├── ssh/                  # SSH 连接封装（待实现）
│   └── utils/                # 通用函数
├── go.mod / go.sum
├── main.go                   # CLI 入口
├── README.md
├── TODO.md                   # 开发计划
├── tests/                    # 集成测试
│   └── health_test.go        # health 命令端到端测试
└── reasonix.toml
```

