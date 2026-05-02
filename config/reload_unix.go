//go:build linux || darwin || freebsd || openbsd || netbsd

package config

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const defaultPidFile = "/tmp/frptables.pid"

// getPidFile 获取 PID 文件路径
func getPidFile() string {
	// 可以根据配置覆盖，暂时用默认值
	return defaultPidFile
}

// SavePid 保存 PID 文件
func SavePid() error {
	pidFile := getPidFile()
	dir := filepath.Dir(pidFile)
	os.MkdirAll(dir, 0755)

	return os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
}

// RemovePid 删除 PID 文件
func RemovePid() {
	pidFile := getPidFile()
	os.Remove(pidFile)
}

// SetupReloadHandler 设置信号处理器，监听 SIGUSR1
func SetupReloadHandler(reloadFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	go func() {
		for range c {
			fmt.Println("Received SIGUSR1, reloading config...")
			reloadFunc()
		}
	}()
}

// SendReload 发送 reload 信号给运行中的进程
func SendReload() error {
	pidFile := getPidFile()

	data, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("failed to read pid file %s: %v", pidFile, err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return fmt.Errorf("invalid pid in file: %v", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process not found: %v", err)
	}

	err = process.Signal(syscall.SIGUSR1)
	if err != nil {
		return fmt.Errorf("failed to send signal: %v", err)
	}

	fmt.Println("Reload signal sent successfully")
	return nil
}
