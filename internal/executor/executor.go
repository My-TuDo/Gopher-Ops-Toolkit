package executor

import "context"

// Executor 接口定义了一个执行器应该有的能力
type Executor interface {
	// Exec 执行命令，返回一个 Result channel,实时流式传输
	// ctx 控制超时和取消
	Exec(ctx context.Context, command string) (<-chan ExecResult, error)
}
