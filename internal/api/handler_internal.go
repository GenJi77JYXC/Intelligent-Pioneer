package api

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/core/engine"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/model"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TriggerKB 手动触发一个知识库任务
func TriggerKB(c *gin.Context) {
	// 1. 绑定并校验请求参数
	var req TriggerKBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warnw("Invalid request to trigger KB", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	logger.L.Infow("Manual KB trigger received", "agent_id", req.AgentID, "kb_id", req.KBID)

	// 2.验证 Agent 是否存在且在线
	var agent model.Agent
	// First 方法会在找不到记录时返回 gorm.ErrRecordNotFound 错误
	result := store.DB.Where("uuid = ? AND status = ?", req.AgentID, "online").First(&agent)
	if result.Error != nil {
		logger.L.Warnw("Agent not found or not online for KB trigger", "agent_id", req.AgentID, "error", result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found or is offline."})
		return
	}

	// 3. 调用引擎，启动工作流
	// 注意：StartKBWorkflow 目前返回的是 error，未来可以修改它返回 (workflowID, error)
	workflowID, err := engine.StartKBWorkflow(req.AgentID, req.KBID)
	if err != nil {
		logger.L.Errorw("Failed to start KB workflow", "agent_id:", req.AgentID, "kb_id:", req.KBID, "workflowID:", workflowID, "error", err)
		// 根据错误类型返回不同的 HTTP 状态码
		// 例如，如果是因为 KB 不存在，可以返回 404
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start workflow: " + err.Error()})
		return
	}

	// 4. 返回成功响应
	// 当 StartKBWorkflow 返回 workflowID 时，在这里一并返回
	c.JSON(http.StatusOK, gin.H{
		"message":     "KB workflow triggered successfully.",
		"workflow_id": workflowID,
	})
	logger.L.Info("Manual KB trigger received")
}

// HealthCheck 是一个简单的健康检查端点
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
