package bootstrap

import (
	"errors"
	"net/http"
	"sync/atomic"
	"time"
)

//func InitServer(address string) {
//	server := &http.Server{
//		Addr:    address,
//		Handler: engine,
//	}
//	go func() {
//		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
//			Release()
//			global.Logger.Fatal("server listen fail", zap.Error(err))
//		}
//	}()
//	time.Sleep(time.Second)
//	releaseFunc = append(releaseFunc, func() {
//		_ = server.Close()
//	})
//	global.Logger.Info("server listen on " + address)
//}

var server *http.Server
var mu atomic.Int32

func StartServer(address string) (err error) {
	if mu.Load() == 1 {
		return errors.New("is running")
	}
	server = &http.Server{
		Addr:    address,
		Handler: engine,
	}
	go func() {
		err = server.ListenAndServe()
	}()
	time.Sleep(time.Second)
	if err != nil {
		return err
	}
	mu.Store(1)
	return nil
}

func StopServer() {
	if mu.Load() == 0 {
		return
	}
	if server == nil {
		return
	}
	_ = server.Close()
	mu.Store(0)
}

func IsRunningServer() bool {
	return mu.Load() == 1
}
