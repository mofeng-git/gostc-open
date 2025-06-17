package main

import (
	_ "embed"
	"fmt"
	"gioui.org/io/system"
	service2 "github.com/SianHH/frp-package/package"
	"github.com/energye/systray"
	"gostc-sub/gui/global"
	"gostc-sub/gui/registry"
	"gostc-sub/gui/window"
	"gostc-sub/internal/common"
	"gostc-sub/webui/backend/bootstrap"
	_ "gostc-sub/webui/backend/router"
	_ "gostc-sub/webui/backend/todo"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
)

//go:embed static/logo.ico
var logo []byte

var win *window.Window

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	moveWebuiCfgDir(global.BasePath, global.BasePath+"/data") // 移动旧配置文件内容
	global.BasePath = global.BasePath + "/data"
	if err := global.LoadConfig(global.BasePath + "/config.yaml"); err != nil {
		log.Fatalln(err)
	}

	bootstrap.InitLogger()
	bootstrap.InitFS(global.BasePath)
	bootstrap.InitTodo()
	bootstrap.InitTask()
	bootstrap.InitRouter()
	_ = bootstrap.StartServer(global.Config.Address)

	systray.SetIcon(logo)
	systray.SetTitle("GOSTC")
	systray.SetTooltip("GOSTC GUI " + common.VERSION)
	systray.SetOnClick(func(menu systray.IMenu) {
		if WindowShow {
			win.Perform(system.ActionRaise)
			win.Perform(system.ActionUnmaximize)
			return
		} else {
			win = window.NewWindow("GOSTC GUI "+common.VERSION, window.OptionCloseCallback(func() {
				WindowShow = false
			}))
			_ = win.Create()
			WindowShow = true
		}
	})
	systray.SetOnRClick(func(menu systray.IMenu) {
		_ = menu.ShowMenu()
	})

	mRegistry := systray.AddMenuItem("开机自启", "修改注册表实现开机自启动")
	if registry.Registered() {
		mRegistry.Check()
	}
	mRegistry.Click(func() {
		if mRegistry.Checked() {
			registry.UnRegister()
			mRegistry.Uncheck()
		} else {
			if err := registry.Register(); err == nil {
				mRegistry.Check()
			} else {
				fmt.Println(err)
			}
		}
	})

	mOpen := systray.AddMenuItem("在浏览器打开", "用浏览器打开Web服务地址")
	mOpen.Click(func() {
		_, port, _ := net.SplitHostPort(global.Config.Address)
		cmd := exec.Command("cmd", "/c", "start", fmt.Sprintf("http://127.0.0.1:%s", port))
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
		_ = cmd.Start()
	})

	mQuit := systray.AddMenuItem("退出", "彻底退出程序")
	mQuit.Click(func() {
		systray.Quit()
	})
}

func onExit() {
	if WindowShow {
		win.Close()
	}
	bootstrap.Release()
	service2.Range(func(key string, value service2.Service) {
		value.Stop()
	})
	fmt.Println("onExit")
}

func moveWebuiCfgDir(basePath string, targetPath string) {
	_ = os.MkdirAll(targetPath, 0755)
	if cfgData, err := os.ReadFile(basePath + "/config.yaml"); err == nil {
		if err = os.WriteFile(targetPath+"/config.yaml", cfgData, 0644); err == nil {
			_ = os.Remove(basePath + "/config.yaml")
		}
	}
	if cfgData, err := os.ReadFile(basePath + "/client.json"); err == nil {
		if err = os.WriteFile(targetPath+"/client.json", cfgData, 0644); err == nil {
			_ = os.Remove(basePath + "/client.json")
		}
	}
	if cfgData, err := os.ReadFile(basePath + "/tunnel.json"); err == nil {
		if err = os.WriteFile(targetPath+"/tunnel.json", cfgData, 0644); err == nil {
			_ = os.Remove(basePath + "/tunnel.json")
		}
	}
	if cfgData, err := os.ReadFile(basePath + "/p2p.json"); err == nil {
		if err = os.WriteFile(targetPath+"/p2p.json", cfgData, 0644); err == nil {
			_ = os.Remove(basePath + "/p2p.json")
		}
	}
}
