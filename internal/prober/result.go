package prober

type Result struct {
	Name      string `json:"name"`              // 项目名称
	Target    string `json:"target"`            // 目标地址
	Status    string `json:"status"`            // 健康状态
	Error     string `json:"error,omitempty"`   // 错误信息（如果有）
	Detail    string `json:"detail,omitempty"`  // 详细信息（可选）
	Latency   int64  `json:"latency,omitempty"` // 响应时间（如果健康）
	ErrorCode string `json:"error_code,omitempty"`
}

const (
	ErrCodeNone        = ""             // 健康，无错误
	ErrCodeUnreachable = "UNREACHABLE"  // 目标不可达
	ErrCodeTimeout     = "TIMEOUT"      // 超时
	ErrCodeUnsupported = "UNSUPPORTED"  // 功能不支持
	ErrCodeDNSFailure  = "DNS_FAILURE"  // DNS 解析失败
	ErrCodeSystemError = "SYSTEM_ERROR" // 系统错误
)

// // String 格式化输出结果
// func (r Result) String() string {
// 	if r.Error != "" {
// 		return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s, 错误: %s", r.Name, r.Target, r.Status, r.Error)
// 	}
// 	return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s, 详情：%s", r.Name, r.Target, r.Status, r.Detail)
// }
