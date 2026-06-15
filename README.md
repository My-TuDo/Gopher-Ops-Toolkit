# Gopher-Ops-Toolkit

一个用`GO`编写的命令行工具，提供各种运维常用的小功能。


ops-toolkit/
├── cmd/
│   ├── exec.go         # 远程命令执行
│   ├── health.go       # 健康检查命令
│   ├── log.go          # 日志分析命令
│   ├── monitor.go      # 批量执行命令
│   ├── root.go         # 根命令 & 配置初始化
│   └── version.go      # 版本信息
├── configs/
│   └── config.yaml     # 配置文件
├── internal/
│   ├── output/         # 统一输出格式（JSON/表格）
│   ├── prober/         # 健康检查探针
│   │   ├── http.go     # HTTP 探测
│   │   ├── result.go   # 探测结果结构体
│   │   └── tcp.go      # TCP 端口探测
│   ├── ssh/            # SSH 连接封装
│   └── utils/          # 通用函数
├── go.mod
├── go.sum
├── LICENSE
├── main.go             # CLI 入口（使用 cobra 库）
├── README.md
└── reasonix.toml
