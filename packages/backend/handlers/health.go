package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查处理器
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Backend service is running123",
	})
}

// Ping 测试端点处理器
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong1235",
	})
}
