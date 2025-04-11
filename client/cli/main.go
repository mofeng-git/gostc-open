package main

import (
	"fmt"
	system_service "github.com/kardianos/service"
	_ "gostc-sub/webui/backend/router"
	_ "gostc-sub/webui/backend/todo"
	"os"
	"path/filepath"
)

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
	SvcCfg.WorkingDirectory, _ = os.Executable()
	SvcCfg.WorkingDirectory = filepath.Dir(SvcCfg.WorkingDirectory)
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
