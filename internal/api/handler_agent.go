package api

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterAgent 处理 Agent 注册请求
func RegisterAgent(c *gin.Context) {
	// TODO: 实现 Agent 注册逻辑
	// 1. 从请求体中解析 Agent 信息 (JSON)
	// 2. 生成一个唯一的 UUID
	// 3. 将 Agent 信息存入 PostgreSQL
	// 4. 返回 UUID 给 Agent
	logger.L.Info("Received agent registration request")
	c.JSON(http.StatusOK, gin.H{
		"message":  "Agent registered successfully (mock)",
		"agent_id": "mock-uuid-12345",
	})
}

// Heartbeat 处理 Agent 心跳请求
func Heartbeat(c *gin.Context) {
	// TODO: 实现 Agent 心跳逻辑
	// 1. 从请求体中获取 agent_id
	// 2. 更新数据库中该 Agent 的状态为 "online" 和最后心跳时间
	logger.L.Info("Received agent heartbeat")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetTasks 处理 Agent 获取任务的长轮询请求
func GetTasks(c *gin.Context) {
	// TODO: 实现长轮询获取任务的逻辑
	// 1. 从查询参数中获取 agent_id
	// 2. 检查该 Agent 是否有待执行的任务
	// 3. 如果有，立即返回任务信息
	// 4. 如果没有，阻塞请求一段时间 (e.g., 30秒)，期间如果收到新任务则返回，否则超时后返回空
	logger.L.Info("Agent polling for tasks...")
	c.JSON(http.StatusOK, gin.H{"tasks": []string{}}) // 暂时返回空任务列表
}

// PostTaskResults 处理 Agent 上报任务结果的请求
func PostTaskResults(c *gin.Context) {
	// TODO: 实现任务结果上报逻辑
	// 1. 从请求体中解析任务ID和结果
	// 2. 更新数据库中任务的状态
	// 3. 根据结果，可能触发知识库工作流的下一步
	logger.L.Info("Received task results from agent")
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
