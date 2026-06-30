package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/prober"
)

func RenderResult(results []prober.Result, format string) {
	switch format {
	case "json":
		RenderResultJSON(results) // 终端打印 JSON
	default:
		RenderResultTable(results)
	}
}

// 返回 JSON 格式结果
func RenderResultJSON(results []prober.Result) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("JSON 序列化失败: %v\n", err))
		return
	}
	fmt.Println(string(data))
}
