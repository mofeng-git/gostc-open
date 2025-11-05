package router

import (
    "archive/zip"
    "bytes"
    "fmt"
    "github.com/gin-gonic/gin"
    "html/template"
    "io"
    "net/http"
    "path/filepath"
    "server/model"
    cache2 "server/repository/cache"
    dashboardSvc "server/service/admin/dashboard"
    "server/web"
    "strings"
    "time"
)

var fileContentTypeMap = map[string]string{
	".js":       "application/javascript",
	".mjs":      "application/javascript",
	".css":      "text/css",
	".manifest": "text/cache-manifest",
	".png":      "image/png",
	".jpg":      "image/jpeg",
	".jpeg":     "image/jpeg",
	".svg":      "image/svg+xml",
	".ico":      "image/x-icon",
	".json":     "application/json",
	".html":     "text/html; charset=utf-8",
	".htm":      "text/html; charset=utf-8",
	".txt":      "text/plain; charset=utf-8",
	".wasm":     "application/wasm",
}

// MatchFile 更安全地匹配文件Content-Type
func MatchFile(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if contentType, ok := fileContentTypeMap[ext]; ok {
		return contentType
	}
	return "application/octet-stream" // 更安全的默认值
}

func StaticFile(zipFile []byte, callback func(fileMap map[string][]byte)) error {
	result := make(map[string][]byte)
	zipReader, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return fmt.Errorf("加载静态资源失败: %w", err)
	}

	for _, file := range zipReader.File {
		open, err := file.Open()
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %w", file.Name, err)
		}

		data, err := io.ReadAll(open)
		_ = open.Close()
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", file.Name, err)
		}

		result[file.Name] = data
	}

	callback(result)
	return nil
}

func InitStatic(engine *gin.Engine) error {
    return StaticFile(web.Static(), func(fileMap map[string][]byte) {
        // 预检查index.html是否存在
        indexHTML, ok := fileMap["dist/index.html"]
        if !ok {
            panic("dist/index.html 文件不存在")
        }

		for fileKey, fileBytes := range fileMap {
			// 创建局部变量避免闭包问题
			fileKey := fileKey
			fileBytes := fileBytes

			ginStaticFilePath := strings.TrimPrefix(fileKey, "dist/")
			if ginStaticFilePath == "" {
				continue
			}

            // 1) 默认根路径下的静态资源，例如 /assets/xxx
            engine.GET(ginStaticFilePath,
                cacheControlMiddleware(),
                serveFileHandler(fileKey, fileBytes),
            )

            // 2) 兼容 Vite 构建使用的 base 前缀（web/vite.config.js: base: '/extras/gostc/'）
            //    让 /extras/gostc/assets/xxx 也能命中到同一份静态资源
            engine.GET("extras/gostc/"+ginStaticFilePath,
                cacheControlMiddleware(),
                serveFileHandler(fileKey, fileBytes),
            )
		}

        engine.NoRoute(func(c *gin.Context) {
            if c.Request.URL.Path == "/" {
                // 首页分流：未开启则302到 /login；开启则渲染首页
                var hCfg model.SystemConfigHome
                cache2.GetSystemConfigHome(&hCfg)
                if hCfg.HomeEnable != "1" {
                    c.Redirect(http.StatusFound, "/login")
                    return
                }
                c.Writer.Header().Set("Cache-Control", "no-store")
                htmlBytes, status := renderHomeHTML()
                c.Data(status, "text/html; charset=utf-8", htmlBytes)
                return
            }
            // 其余路径继续交给前端路由
            c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
        })
    })
}

// 中间件和处理器工厂函数
func cacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "public, max-age=86400")
		c.Next()
	}
}

func serveFileHandler(fileKey string, fileBytes []byte) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Data(http.StatusOK, MatchFile(fileKey), fileBytes)
    }
}

// renderHomeHTML 组装并渲染首页HTML（无软缓存，每次请求实时计算）
func renderHomeHTML() ([]byte, int) {
    // 读取首页配置与基础配置（首页开关已在 NoRoute 中处理，这里只负责渲染）
    var hCfg model.SystemConfigHome
    cache2.GetSystemConfigHome(&hCfg)

    // 站点基础信息
    var b model.SystemConfigBase
    cache2.GetSystemConfigBase(&b)

    // 统计（实时）
    cnt := dashboardSvc.Service.Count()
    todayFlow := formatBytes(cnt.InputBytes + cnt.OutputBytes)
    tunnelTotal := cnt.Host + cnt.Forward + cnt.Tunnel + cnt.Proxy + cnt.P2P

    data := map[string]any{
        "config": map[string]any{
            "title":   b.Title,
            "favicon": b.Favicon,
        },
        "stats": map[string]any{
            "today_flow":    todayFlow,
            "user_total":    cnt.User,
            "checkin_today": cnt.CheckInTotal,
            "tunnel_total":  tunnelTotal,
            "node_online":   cnt.NodeOnline,
            "client_online": cnt.ClientOnline,
            "updated_at":    time.Now().Format(time.DateTime),
        },
        "login_url": "/login",
    }

    // 模板选择：优先使用自定义模板，否则使用内置默认模板
    tplText := string(web.HomeTpl())
    if strings.TrimSpace(hCfg.HomeTpl) != "" {
        tplText = hCfg.HomeTpl
    }

    // 使用 html/template 渲染，保证数据自动转义
    t, err := template.New("home").Parse(tplText)
    if err != nil {
        return []byte("template parse error"), http.StatusInternalServerError
    }
    var buf bytes.Buffer
    if err := t.Execute(&buf, data); err != nil {
        return []byte("template execute error"), http.StatusInternalServerError
    }
    return buf.Bytes(), http.StatusOK
}

// formatBytes 将字节数转成人类可读字符串（B/KB/MB/GB/TB，保留两位小数）
func formatBytes(n int64) string {
    units := []string{"B", "KB", "MB", "GB", "TB"}
    f := float64(n)
    i := 0
    for f >= 1024 && i < len(units)-1 {
        f /= 1024
        i++
    }
    if i == 0 {
        return fmt.Sprintf("%d%s", n, units[i])
    }
    return fmt.Sprintf("%.2f%s", f, units[i])
}
