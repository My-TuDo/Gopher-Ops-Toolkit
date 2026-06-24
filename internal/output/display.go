package output

import "github.com/My-TuDo/Gopher-Ops-Toolkit/internal/prober"

func RenderResult(results []prober.Result, format string) {
	switch format {
	case "json":
		RenderResultJSON(results) // 终端答应 JSON
	default:
		RenderResultTable(results)
	}
}

func RenderResultJSON(results []prober.Result) {

}
