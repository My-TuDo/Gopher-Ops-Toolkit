package output

import (
	"fmt"
	"os"
	"strconv"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/prober"
	"github.com/olekukonko/tablewriter"
)

func RenderResultTable(results []prober.Result) {
	if len(results) == 0 {
		fmt.Println("没有可显示的结果")
		return
	}

	// 设置表头
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"探测项目", "探测目标", "状态", "耗时(ms)", "详情/错误信息"})

	// 循环注入数据
	for _, res := range results {
		latencyStr := "—"
		if res.Latency > 0 {
			latencyStr = strconv.FormatInt(res.Latency, 10) + "ms"
		}

		detailOrError := res.Detail
		if res.Error != "" {
			detailOrError = res.Error
		}

		row := []string{res.Name, res.Target, res.Status, latencyStr, detailOrError}

		table.Append(row)
	}

	fmt.Println()

	// 渲染输出
	table.Render()
	fmt.Println()
}
