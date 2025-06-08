package bootstrap

import (
	"go.uber.org/zap"
	"gostc-sub/internal/service"
	"gostc-sub/pkg/fs"
	"gostc-sub/webui/backend/global"
	"os"
)

func InitFS(basePath string) {
	if basePath == "" {
		basePath = "."
	} else {
		_ = os.MkdirAll(basePath, 0755)
	}
	var err error
	global.ClientFS, err = fs.NewFileStorage(basePath + "/client.json")
	if err != nil {
		Release()
		global.Logger.Fatal("load client file fail", zap.Error(err))
	}
	global.ClientFS.SetAutoPersistInterval(0)
	releaseFunc = append(releaseFunc, func() {
		global.ClientFS.Close()
	})
	global.TunnelFS, err = fs.NewFileStorage(basePath + "/tunnel.json")
	if err != nil {
		Release()
		global.Logger.Fatal("load tunnel file fail", zap.Error(err))
	}
	global.TunnelFS.SetAutoPersistInterval(0)
	releaseFunc = append(releaseFunc, func() {
		global.TunnelFS.Close()
	})
	global.P2PFS, err = fs.NewFileStorage(basePath + "/p2p.json")
	if err != nil {
		Release()
		global.Logger.Fatal("load p2p file fail", zap.Error(err))
	}
	global.P2PFS.SetAutoPersistInterval(0)
	releaseFunc = append(releaseFunc, func() {
		global.P2PFS.Close()
	})
	releaseFunc = append(releaseFunc, func() {
		global.ClientMap.Range(func(key, value any) bool {
			value.(service.Service).Stop()
			return true
		})
		global.TunnelMap.Range(func(key, value any) bool {
			value.(service.Service).Stop()
			return true
		})
		global.P2PMap.Range(func(key, value any) bool {
			value.(service.Service).Stop()
			return true
		})
	})
}
