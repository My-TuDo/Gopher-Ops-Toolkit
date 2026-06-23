package output

import (
	"fmt"
	"os"
	"strconv"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/prober"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"golang.org/x/term"
)

func RenderResultTable(results []prober.Result) {
	if len(results) == 0 {
		fmt.Println("没有可显示的结果")
		return
	}
	// 获取终端宽度
	termWidth := 120
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termWidth = w
	}

	// 详情列宽度 = 总宽 - 其他列固定宽 - 边框
	detailWidth := termWidth - 8 - 26 - 6 - 8 - 14
	if detailWidth > 60 {
		detailWidth = 60 // 限制详情列最大宽度为 60
	}

	// 设置表头
	table := tablewriter.NewTable(
		os.Stdout,
		// 设置表格渲染器，使用自定义的蓝图
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{
				Separators: tw.Separators{
					BetweenRows: tw.On, // 行间分隔线
				},
			},
		})),
		// 设置表格配置，包括列宽、对齐方式等
		tablewriter.WithConfig(tablewriter.Config{
			MaxWidth: termWidth,
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap: tw.WrapNormal, // 单元格内自动换行
				},
				Alignment: tw.CellAlignment{
					Global: tw.AlignLeft, // 全局左对齐
				},
				ColMaxWidths: tw.CellWidth{
					PerColumn: tw.Mapper[int, int]{0: 8, 1: 26, 2: 6, 3: 8, 4: detailWidth}, // 每列最大宽度
				},
			},
			// 设置表头配置，包括对齐方式
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignCenter}, // 表头居中对齐
			},
		}),
	)
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

		// 超长文本截断处理
		if len(detailOrError) > 80 {
			detailOrError = string([]rune(detailOrError)[:80]) + "..."
		}

		row := []string{res.Name, res.Target, res.Status, latencyStr, detailOrError}

		table.Append(row)
	}
	// 渲染输出
	table.Render()
}
