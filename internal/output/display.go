package output

import (
	"encoding/json"
	"fmt"

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
		panic(err)
	}
	fmt.Println(string(data))
}
