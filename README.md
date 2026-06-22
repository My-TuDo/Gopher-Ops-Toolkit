# Gopher-Ops-Toolkit

一个用`GO`编写的命令行工具，提供各种运维常用的小功能。

---

## 快速开始

### 构建

```bash 
go build -o ops-toolkit
```

### 配置
编辑 `configs/config.yaml`，或通过 `--config` 指定自定义路径。

### 使用指南

- **`health`——健康检查**

```bash
# 执行全部探测（从配置文件读取）
ops-toolkit health --all

# TCP 端口探测
ops-toolkit health --type tcp --host example.com --port 3306

# HTTP 服务探测
ops-toolkit health --type http --url http://example.com/health

# DNS 解析探测
ops-toolkit health --type dns --domain baidu.com --reacord-type A
```

- **`version`——版本信息**
```bash
ops-toolkit version
```

---

## 项目结构

```
ops-toolkit/
├── cmd/                      # 命令定义
│   ├── exec.go               # 远程命令执行
│   ├── health.go             # 健康检查命令
│   ├── log.go                # 日志分析命令
│   ├── monitor.go            # 批量执行命令
│   ├── root.go               # 根命令 & 配置初始化
│   └── version.go            # 版本信息
├── configs/
│   └── config.yaml           # 配置文件
├── internal/
│   ├── output/               # 统一输出格式
│   ├── prober/               # 健康检查探针
│   │   ├── dns.go            # DNS 探测
│   │   ├── http.go           # HTTP 探测
│   │   ├── result.go         # 探测结果结构体 & 通用函数
│   │   └── tcp.go            # TCP 端口探测
│   ├── ssh/                  # SSH 连接封装
│   └── utils/                # 通用函数
├── go.mod
├── go.sum
├── LICENSE
├── main.go                   # CLI 入口
├── README.md
├── TODO.md                   # 开发计划
└── reasonix.toml
```

