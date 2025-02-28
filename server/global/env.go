package global

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var basePath string
var devMode bool
var logLevel string

func CmdPerInit(cmds ...*cobra.Command) {
	executable, _ := os.Executable()
	for _, cmd := range cmds {
		cmd.Flags().StringVarP(&basePath, "path", "p", executable, "app run dir")
		cmd.Flags().StringVarP(&logLevel, "log-level", "", "warn", "log level debug|info|warn|error|fatal")
		cmd.Flags().BoolVarP(&devMode, "dev", "d", false, "app run dev")
	}
}

func CmdInit() {
	// 模式
	if devMode {
		MODE = "dev"
	} else {
		MODE = "prod"
	}
	LOGGER_LEVEL = logLevel
	BASE_PATH = filepath.Dir(basePath)
	LOGGER_FILE_PATH = BASE_PATH + "/data/gostc.log"

	fmt.Printf(`
========================================
MODE: %s
VERSION: %s
BASE_PATH: %s
LOGGER_FILE_PATH: %s
LOGGER_LEVEL: %s
========================================
`, MODE, VERSION, BASE_PATH, LOGGER_FILE_PATH, LOGGER_LEVEL)

}

var (
	VERSION   = "v1.0.0" // 版本
	BASE_PATH = ""       // 基础目录
	/*
		prod:生产模式，程序运行根路径和程序的所在目录保持一致
		dev:开发模式，程序运行根路径和pwd输出的目录保持一致
	*/
	MODE             = "dev"
	LOGGER_FILE_PATH = ""     // 日志文件路径
	LOGGER_LEVEL     = "info" // 日志等级
)
