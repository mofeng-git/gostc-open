package registry

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
)

var (
	root    = registry.LOCAL_MACHINE
	path    = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	appName = "GOSTC_WEBGUI"
)

func Registered() bool {
	key, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer key.Close()
	newAbsPath, _ := os.Executable()
	oldAbsPath, _, err := key.GetStringValue(appName)
	if err != nil {
		return false
	}
	if oldAbsPath != newAbsPath {
		return false
	}
	return true

}

func Register() error {
	key, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %v", err)
	}
	defer key.Close()
	newAbsPath, _ := os.Executable()
	err = key.SetStringValue(appName, newAbsPath)
	if err != nil {
		return fmt.Errorf("修改注册表失败: %v", err)
	}
	return nil
}

func UnRegister() {
	key, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return
	}
	defer key.Close()
	_ = key.DeleteValue(appName)
}
