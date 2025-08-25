package service

import (
	"crypto/tls"
	"fmt"
	"github.com/fatedier/golib/log"
	"github.com/lesismal/arpc"
	"github.com/radovskyb/watcher"
	"gostc-sub/gui/global"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func certWatcher(client *arpc.Client) (done func()) {
	global.BasePath = "."
	certpath, _ := filepath.Abs(global.BasePath + "/data/certs")
	_ = os.MkdirAll(certpath, 0755)
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(
		watcher.Create,
		watcher.Write,
		watcher.Rename,
		watcher.Chmod,
		watcher.Move,
	)
	go func() {
		for {
			select {
			case e := <-w.Event:
				if e.Path == certpath {
					continue
				}
				// 跳过不符合的目录
				if len(strings.Split(e.Name(), "/")) > 1 {
					continue
				}
				fmt.Println("域名：", e.Name())
				certFile := e.Path + "/cert.pem"
				keyFile := e.Path + "/cert.key"

				certStat, err := os.Stat(certFile)
				if err != nil {
					continue
				}
				if certStat.IsDir() {
					continue
				}
				keyStat, err := os.Stat(keyFile)
				if err != nil {
					continue
				}
				if keyStat.IsDir() {
					continue
				}
				if certStat.Size() > 1024*128 || keyStat.Size() > 1024*128 {
					log.Warn("the file size exceeds 128KB")
					continue
				}
				certFileBytes, err := os.ReadFile(certFile)
				if err != nil {
					continue
				}
				keyFileBytes, err := os.ReadFile(keyFile)
				if err != nil {
					continue
				}
				if _, err := tls.X509KeyPair(certFileBytes, keyFileBytes); err != nil {
					log.Warn("certificate verification failed.", err)
					continue
				}
				var reply string
				if err := client.Call("rpc/client/host_set_domain_certs", map[string]string{
					"domain": e.Name(),
					"cert":   string(certFileBytes),
					"key":    string(keyFileBytes),
				}, &reply, time.Second*10); err != nil {
					log.Warn("failed to upload certificate.", e.Name(), certFile, keyFile)
					continue
				}
				if reply != "success" {
					log.Warn("failed to upload certificate.", e.Name(), reply)
					continue
				}
			case _ = <-w.Error:
				time.Sleep(time.Second)
				_ = os.MkdirAll(certpath, 0755)
				time.Sleep(time.Second)
				if err := w.AddRecursive(certpath); err != nil {
					log.Warn("watcher certificate directory failed", err)
					return
				}
			case <-w.Closed:
				return
			}
		}
	}()
	if err := w.AddRecursive(certpath); err != nil {
		log.Warn("watcher certificate directory failed", err)
		return
	}
	if err := w.Start(time.Second); err != nil {
		log.Warn("watcher certificate directory failed", err)
		return
	}
	log.Info("watcher certificate directory success")
	done = func() {
		w.Close()
	}
	return done
}
