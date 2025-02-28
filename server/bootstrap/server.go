package bootstrap

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"server/global"
	"time"
)

func InitServer() {
	server := &http.Server{
		Addr:    global.Config.Address,
		Handler: engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			Release()
			global.Logger.Fatal("server listen fail", zap.Error(err))
		}
	}()
	time.Sleep(time.Second)
	releaseFunc = append(releaseFunc, func() {
		_ = server.Close()
	})
	global.Logger.Info("server listen on address: " + global.Config.Address)
}
