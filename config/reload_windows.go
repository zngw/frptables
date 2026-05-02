//go:build windows

package config

import (
	"fmt"
)

// SavePid Windows 上不保存 PID 文件
func SavePid() error {
	// Windows 不支持 signal-based reload
	return nil
}

// RemovePid Windows 上不删除 PID 文件
func RemovePid() {
	// Windows 不支持 signal-based reload
}

// SetupReloadHandler Windows 上不设置 signal handler
func SetupReloadHandler(reloadFunc func()) {
	// Windows 不支持 SIGUSR1
	fmt.Println("Warning: reload via signal not supported on Windows")
}

// SendReload Windows 上不支持
func SendReload() error {
	return fmt.Errorf("reload not supported on Windows, please restart the service")
}
