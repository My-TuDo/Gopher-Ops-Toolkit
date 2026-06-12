# Gopher-Ops-Toolkit

一个用`GO`编写的命令行工具，提供各种运维常用的小功能。


ops-toolkit/
├── cmd
│   ├── exec.go         # 子命令目录
│   ├── health.go       # 健康检查命令
│   ├── log.go          # 日志分析命令
│   ├── monitor.go      # 批量执行命令
│   ├── root.go         # 简易监控命令
│   └── version.go      # 版本信息
├── configs
│   └── config.yaml
├── go.mod
├── go.sum
├── internal            # 内部共享代码
│   ├── output          # 统一输出格式（JSON/表格）
│   ├── ssh             # SSH 连接封装
│   └── utils           # 通用函数
├── LICENSE
├── main.go             # CLI 入口（使用 cobra 库）
├── README.md
└── reasonix.toml        
        

        
        
       
       
        
        

        
