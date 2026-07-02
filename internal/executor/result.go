package executor

type ExecResult struct {
	Line     string // 一行输出
	IsStderr bool   // 是否是 stderr
	Done     bool   // 是否执行完毕
	ExitCode int    // 退出码（仅在 Done == true 时有意义）
	Err      error  // 错误信息（仅在Done == true 时有意义）
}
