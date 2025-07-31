package main

import (
	"backend/handlers"
	"backend/middleware"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// StaticFileHandler 处理静态文件服务，支持 SPA 路由
func StaticFileHandler(staticPath string) gin.HandlerFunc {
	fileServer := http.FileServer(http.Dir(staticPath))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是 API 请求，跳过静态文件处理
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/health") {
			c.Next()
			return
		}

		// 检查文件是否存在
		fullPath := filepath.Join(staticPath, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// 如果文件不存在，返回 index.html（支持 SPA 路由）
			c.Request.URL.Path = "/"
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func setupLogging() {
	// 确保日志目录存在
	logDir := "/opt/logs/100056529"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		return
	}

	// 创建日志文件，按日期命名
	logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}

	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	gin.DefaultWriter = multiWriter

	log.Printf("Logging initialized, writing to: %s", logFileName)
}

func main() {
	// 设置日志
	setupLogging()

	// 初始化数据库
// 	config.InitDatabase()

	// 创建 Gin 路由器
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(middleware.CORS())

	// 健康检查端点
	r.GET("/health", handlers.HealthCheck)

	// API 路由组
	api := r.Group("/api/v1")
	{
		api.GET("/ping", handlers.Ping)
	}

	// 静态文件服务（放在最后，作为 fallback）
	r.Use(StaticFileHandler("/static"))

	// 启动服务器，默认端口 8080
	r.Run(":8080")
}
