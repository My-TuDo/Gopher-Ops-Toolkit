// Package config 提供 XDG 标准路径的配置管理，
// 并将默认配置通过 //go:embed 编译进二进制。
package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed default.yaml
var DefaultConfig []byte

// ResolveConfigPath 按优先级查找配置文件路径：
//
//	 1. cfgFile（--config 显式指定）
//	 2. GOPT_CONFIG 环境变量
//	 3. XDG 默认路径（~/.config/gopt/config.yaml）
//	 4. 本地开发路径（configs/config.yaml）
//		 → 都不存在时返回 XDG 路径（后续由 EnsureDefaultConfig 自动创建）
func ResolveConfigPath(cfgFile string) (path string, exists bool) {
	// 优先级 1: --config 显式指定
	if cfgFile != "" {
		if exists := pathExists(cfgFile); exists {
			return cfgFile, true
		}
		// 用户显式指定的文件不存在才返回（权限不足时仍尝试加载，让 viper 报错）
		return cfgFile, false
	}

	// 优先级 2: GOPT_CONFIG 环境变量
	if envPath := os.Getenv("GOPT_CONFIG"); envPath != "" {
		if exists := pathExists(envPath); exists {
			return envPath, true
		}
		// 环境变量指向的文件不存在，继续往下查找
	}

	// 优先级 3: XDG 默认路径
	xdgPath := xdgConfigPath()
	if pathExists(xdgPath) {
		return xdgPath, true
	}

	// 优先级 4: 本地开发路径
	localPath := filepath.Join("configs", "config.yaml")
	if pathExists(localPath) {
		return localPath, true
	}

	// 都不存在 → 返回 XDG 路径，由调用方决定是否自动创建
	return xdgPath, false
}

// pathExists 检查文件是否存在，仅当明确「不存在」时返回 false。
// 权限不足等错误会 stderr 告警并返回 false（让调用方 fallback）。
func pathExists(p string) bool {
	_, err := os.Stat(p)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// 权限不足或其他错误：stderr 告警，按不存在处理以触发 fallback
	fmt.Fprintf(os.Stderr, "⚠️ 警告: 无法访问 %s: %v\n", p, err)
	return false
}

// EnsureDefaultConfig 在指定路径写入嵌入的默认配置。
// 如果文件已存在则不操作。返回是否为新创建。
func EnsureDefaultConfig(path string) (created bool, err error) {
	if pathExists(path) {
		return false, nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false, fmt.Errorf("create config dir: %w", err)
	}

	if err := os.WriteFile(path, DefaultConfig, 0644); err != nil {
		return false, fmt.Errorf("write default config: %w", err)
	}

	return true, nil
}

// DefaultLogDir 返回符合 XDG 标准的日志目录。
// 优先使用 $XDG_DATA_HOME，回退 ~/.local/share/gopt/logs 。
func DefaultLogDir() string {
	if envDir := os.Getenv("XDG_DATA_HOME"); envDir != "" {
		return filepath.Join(envDir, "gopt", "logs")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "gopt", "logs")
}

// xdgConfigPath 返回符合 XDG 标准的配置文件路径。
// 优先使用 $XDG_CONFIG_HOME，回退 ~/.config/gopt/config.yaml 。
func xdgConfigPath() string {
	if envDir := os.Getenv("XDG_CONFIG_HOME"); envDir != "" {
		return filepath.Join(envDir, "gopt", "config.yaml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gopt", "config.yaml")
}
