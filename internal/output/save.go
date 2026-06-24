package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/prober"
)

// SaveToFile 将探测结果以 JSON Lines 格式追加写入日志文件
func SaveToFile(results []prober.Result, logDir string) error {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败：%s", err)
	}

	// 生成文件名:{logDir}/health-2026-06-24.jsonl
	today := time.Now().Format("2006-01-02")
	filePath := filepath.Join(logDir, fmt.Sprintf("health-%s.jsonl", today))

	// 以追加模式打开文件，如果不存在则创建
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败：%s", err)
	}
	defer file.Close()

	// 逐条写入 JSON Lines
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false) // 禁止转义 HTML 字符
	for _, res := range results {
		if err := encoder.Encode(res); err != nil {
			return fmt.Errorf("写入日志失败：%s", err)
		}
	}

	return nil
}
