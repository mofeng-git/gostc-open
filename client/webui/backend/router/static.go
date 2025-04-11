package router

import (
	"archive/zip"
	"bytes"
	"github.com/gin-gonic/gin"
	"gostc-sub/webui/backend/web"
	"io"
	"net/http"
	"strings"
)

var fileContentTypeMap = []struct {
	Ext         string
	ContentType string
}{
	{Ext: ".js", ContentType: "application/javascript"},
	{Ext: ".css", ContentType: "text/css"},
	{Ext: "manifest", ContentType: "application/octet-stream"},
	{Ext: ".png", ContentType: "image/png"},
	{Ext: ".jpg", ContentType: "image/jpeg"},
	{Ext: ".jpeg", ContentType: "image/jpeg"},
	{Ext: "", ContentType: "text/html; charset=utf-8"},
}

// MatchFile 匹配不同文件类型的Content-Type
func MatchFile(fileName string) (result string) {
	result = "text/html; charset=utf-8"
	for _, value := range fileContentTypeMap {
		if strings.Contains(fileName, value.Ext) {
			result = value.ContentType
			break
		}
	}
	return result
}

func StaticFile(zipFile []byte, callback func(fileMap map[string][]byte)) {
	var result = make(map[string][]byte)
	zipReader, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return
	}
	for _, file := range zipReader.File {
		open, err := file.Open()
		if err != nil {
			panic("读取静态资源失败" + err.Error())
		}
		data, err := io.ReadAll(open)
		if err != nil {
			panic("读取静态资源失败" + err.Error())
		}
		_ = open.Close()
		result[file.Name] = data
	}
	callback(result)
}

func InitStatic(engine *gin.Engine) {
	StaticFile(web.Static(), func(fileMap map[string][]byte) {
		for k, data := range fileMap {
			// 规避forRange复用k,data
			fileKey := k
			fileBytes := data
			ginStaticFilePath := strings.Replace(fileKey, "dist/", "", 1)
			if ginStaticFilePath == "" {
				continue
			}
			engine.GET(ginStaticFilePath, func(c *gin.Context) {
				c.Writer.Header().Set("Cache-Control", "public,max-age=86400")
			}, func(c *gin.Context) {
				c.Data(http.StatusOK, MatchFile(fileKey), fileBytes)
			})
		}
		engine.NoRoute(func(c *gin.Context) {
			c.Data(http.StatusOK, "text/html; charset=utf-8", fileMap["dist/index.html"])
		})
	})
}
