package prober

import (
	"fmt"
	"time"
)

type Result struct {
	Name    string `json:"name"`              // 项目名称
	Target  string `json:"target"`            // 目标地址
	Status  string `json:"status"`            // 健康状态
	Error   string `json:"error,omitempty"`   // 错误信息（如果有）
	Detail  string `json:"detail,omitempty"`  // 详细信息（可选）
	Latency int64  `json:"latency,omitempty"` // 响应时间（如果健康）
}

// String 格式化输出结果
func (r Result) String() string {
	if r.Error != "" {
		return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s, 错误: %s", r.Name, r.Target, r.Status, r.Error)
	}
	return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s", r.Name, r.Target, r.Status)
}

// 定义一个接口，所有探测器都实现这个接口
type Prober interface {
	Probe(target string, timeout time.Duration) Result
}

// 全局注册表
var Probers = map[string]Prober{
	"tcp":  &TCPProber{},
	"http": &HTTPProber{},
}
