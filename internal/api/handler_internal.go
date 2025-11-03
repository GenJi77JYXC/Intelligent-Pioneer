package api

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TriggerKB 手动触发一个知识库任务
func TriggerKB(c *gin.Context) {
	// TODO: 实现手动触发知识库的逻辑
	// 1. 从请求体中获取 agent_id 和 kb_id
	// 2. 调用知识库引擎的核心工作流
	logger.L.Info("Manual KB trigger received")
	c.JSON(http.StatusOK, gin.H{"message": "KB task triggered successfully (mock)"})
}

// HealthCheck 是一个简单的健康检查端点
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
