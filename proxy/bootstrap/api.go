package bootstrap

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"proxy/configs"
	"proxy/global"
	"proxy/pkg/middleware"
	"proxy/pkg/proxy"
	"time"
)

type DomainReq struct {
	Domain string `json:"domain"`
	Target string `json:"target"`
	Cert   string `json:"cert"`
	Key    string `json:"key"`
}

func verifyCertificateAndKey(cert, key string) error {
	_, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return fmt.Errorf("证书对验证失败: %v", err)
	}
	return nil
}

func InitApi() {
	if global.Config.ApiAddr == "" {
		return
	}

	if global.MODE == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.Use(middleware.Logger(global.Logger, true, func(c *gin.Context) bool {
		return true
	}))

	engine.POST("/domain", func(c *gin.Context) {
		var req DomainReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.String(500, err.Error())
			return
		}
		if req.Domain == "" {
			return
		}

		if err := verifyCertificateAndKey(req.Cert, req.Key); err != nil {
			global.Logger.Warn("cert valid fail", zap.Error(err))
			return
		}

		var certFile = fmt.Sprintf("%s/data/certs/%s.pem", global.BASE_PATH, req.Domain)
		var keyFile = fmt.Sprintf("%s/data/certs/%s.key", global.BASE_PATH, req.Domain)
		_ = os.WriteFile(certFile, []byte(req.Cert), 0644)
		_ = os.WriteFile(keyFile, []byte(req.Key), 0644)

		global.Config.Domains[req.Domain] = configs.DomainConfig{
			Target: req.Target,
			Cert:   certFile,
			Key:    keyFile,
		}
		marshal, _ := yaml.Marshal(global.Config)
		_ = os.WriteFile(global.BASE_PATH+"/data/config.yaml", marshal, 0644)
		server.UpdateDomain(req.Domain, proxy.DomainConfig{
			Target: req.Target,
			Cert:   req.Cert,
			Key:    req.Key,
		})
	})

	svc := &http.Server{
		Addr:    global.Config.ApiAddr,
		Handler: engine,
	}

	var err error
	go func() {
		err = svc.ListenAndServe()
	}()
	time.Sleep(time.Second)
	if err != nil {
		global.Logger.Warn("api server listen on address: "+global.Config.ApiAddr, zap.Error(err))
		Release()
		os.Exit(1)
	}
	global.Logger.Info("api server listen on address: " + global.Config.ApiAddr)
}
