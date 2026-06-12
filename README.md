# Gopher-Ops-Toolkit

一个用`GO`编写的命令行工具，提供各种运维常用的小功能。


ops-toolkit/
├── cmd/                  # 子命令目录
│   ├── health/           # 健康检查命令
│   ├── log/              # 日志分析命令
│   ├── exec/             # 批量执行命令
│   ├── monitor/          # 简易监控命令
│   └── version/          # 版本信息
├── internal/             # 内部共享代码
│   ├── ssh/              # SSH 连接封装
│   ├── output/           # 统一输出格式（JSON/表格）
│   └── utils/            # 通用函数
├── go.mod
├── main.go               # CLI 入口（使用 cobra 库）
└── README.md