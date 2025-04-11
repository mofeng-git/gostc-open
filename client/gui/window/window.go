package window

import (
	"errors"
	"fmt"
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gostc-sub/gui/global"
	"gostc-sub/webui/backend/bootstrap"
	"log"
	"net"
	"os/exec"
	"sync/atomic"
	"syscall"
	"time"
)

type Option func(window *Window)

func OptionCloseCallback(f func()) Option {
	return func(window *Window) {
		window.closeCallback = f
	}
}

type Window struct {
	win           *app.Window
	title         string
	width         int
	height        int
	mu            atomic.Int32
	closeCallback func()
}

func NewWindow(title string, opts ...Option) *Window {
	win := &Window{title: title}
	for _, opt := range opts {
		opt(win)
	}
	return win
}

func (win *Window) Create() (err error) {
	if win.mu.Load() != 0 {
		return errors.New("is create")
	}
	win.win = new(app.Window)
	win.win.Option(
		app.Title(win.title),
		app.Decorated(true),
		app.Size(unit.Dp(400), unit.Dp(400)),
	)
	go func() {
		err = run(win.win)
		if win.closeCallback != nil {
			win.closeCallback()
		}
		log.Println("close windows")
	}()
	time.Sleep(time.Second)
	if err != nil {
		return err
	}
	win.mu.Store(1)
	return nil
}
func (win *Window) Perform(action system.Action) {
	win.win.Perform(action)
}

func (win *Window) Close() {
	if win.mu.Load() == 0 {
		return
	}
	win.win.Perform(system.ActionClose)
	win.mu.Store(0)
}

func run(window *app.Window) error {
	theme := material.NewTheme()

	//fontTheme := material.NewTheme()
	//fontTheme.ContrastFg = color.NRGBA{
	//	R: 0,
	//	G: 0,
	//	B: 0,
	//	A: 255,
	//}
	var addressEditor widget.Editor
	var saveConfigClick widget.Clickable
	var openUrlClick widget.Clickable
	addressEditor.SetText(global.Config.Address)
	var errMsg = "无异常消息"

	if !bootstrap.IsRunningServer() {
		if err := bootstrap.StartServer(global.Config.Address); err != nil {
			errMsg = err.Error()
		} else {
			errMsg = "无异常消息"
		}
	}

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if saveConfigClick.Clicked(gtx) {
				global.Config.Address = addressEditor.Text()
				_ = global.SaveConfig(global.BasePath + "/config.yaml")
				if bootstrap.IsRunningServer() {
					bootstrap.StopServer()
				}
				if err := bootstrap.StartServer(global.Config.Address); err != nil {
					errMsg = err.Error()
				} else {
					errMsg = "无异常消息"
				}
			}

			if openUrlClick.Clicked(gtx) {
				_, port, _ := net.SplitHostPort(global.Config.Address)
				cmd := exec.Command("cmd", "/c", "start", fmt.Sprintf("http://127.0.0.1:%s", port))
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true,
				}
				_ = cmd.Start()
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return padding(gtx, 10, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return widget.Label{}.Layout(gtx, theme.Shaper, font.Font{}, theme.TextSize, "服务地址：", op.CallOp{})
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return material.Editor(theme, &addressEditor, "******").Layout(gtx)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return padding(gtx, 10, func(gtx layout.Context) layout.Dimensions {
						return material.Button(theme, &saveConfigClick, "保存配置并重启").Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return padding(gtx, 10, func(gtx layout.Context) layout.Dimensions {
						return material.Button(theme, &openUrlClick, "在浏览器打开").Layout(gtx)
					})
				}),
				//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				//	return padding(gtx, 10, func(gtx layout.Context) layout.Dimensions {
				//		_, port, _ := net.SplitHostPort(global.Config.Address)
				//		text := fmt.Sprintf("在浏览器访问http://127.0.0.1:%s管理隧道配置", port)
				//		return material.Label(theme, theme.TextSize, text).Layout(gtx)
				//	})
				//}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return padding(gtx, 10, func(gtx layout.Context) layout.Dimensions {
						text := fmt.Sprintf("服务异常消息：%s", errMsg)
						return material.Label(theme, theme.TextSize, text).Layout(gtx)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}

func padding(gtx layout.Context, padding int, w layout.Widget) layout.Dimensions {
	return layout.UniformInset(unit.Dp(padding)).Layout(gtx, w)
}
