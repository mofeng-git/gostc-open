package main

import (
	"fmt"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config/parsing"
	xlogger "github.com/go-gost/x/logger"
	system_service "github.com/kardianos/service"
	"gostc-sub/common"
	log2 "gostc-sub/p2p/pkg/util/log"
	"log"
	"os"
)

var Writer *common.LogWriter

func Init(logLevel string, console bool) {
	fixDNSResolver()
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	level := logger.WarnLevel
	switch logLevel {
	case "trace":
		level = logger.TraceLevel
	case "debug":
		level = logger.DebugLevel
	case "info":
		level = logger.InfoLevel
	case "warn":
		level = logger.WarnLevel
	case "error":
		level = logger.ErrorLevel
	case "fatal":
		level = logger.FatalLevel
	}
	Writer = common.NewLogWriter(console)
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(level), xlogger.OutputOption(Writer)))
	tlsConfig, _ := parsing.BuildDefaultTLSConfig(nil)
	parsing.SetDefaultTLSConfig(tlsConfig)
	log2.RefreshDefault()
}

func main() {
	var controlMap = map[string]bool{
		"install":   true,
		"uninstall": true,
		"start":     true,
		"stop":      true,
		"restart":   true,
	}
	var svrArgs []string
	for _, item := range os.Args[1:] {
		if controlMap[item] {
			continue
		}
		svrArgs = append(svrArgs, item)
	}
	SvcCfg.Arguments = append(SvcCfg.Arguments, svrArgs...)
	svr, err := system_service.New(Program, SvcCfg)
	if err != nil {
		fmt.Println("build service fail", err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			_ = svr.Stop()
			_ = svr.Uninstall()
			if err := svr.Install(); err != nil {
				fmt.Println("install service fail", err)
				os.Exit(1)
			}
			fmt.Println("install service success")
			return
		case "uninstall":
			_ = svr.Stop()
			if err := svr.Uninstall(); err != nil {
				fmt.Println("uninstall service fail", err)
				os.Exit(1)
			}
			fmt.Println("uninstall service success")
			return
		case "start", "restart", "stop":
			if err := system_service.Control(svr, os.Args[1]); err != nil {
				fmt.Println(os.Args[1]+" service fail", err)
				os.Exit(1)
			}
			fmt.Println(os.Args[1] + " service success")
			return
		}
	}
	_ = svr.Run()
}
