package api

import "github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
import "github.com/gin-gonic/gin"

// NewRouter 创建并配置一个新的 Gin 引擎
func NewRouter() *gin.Engine {
	// 根据配置设置 Gin 的模式 (debug/release)
	gin.SetMode(config.C.Server.Mode)

	// 创建一个不带任何默认中间件的 Gin 引擎
	// 如果需要 Logger 和 Recovery 中间件，可以使用 gin.Default()
	router := gin.New()

	// 在这里可以添加全局中间件
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	// router.Use(middlewares.Cors()) // 例如：跨域中间件

	// 健康检查路由
	router.GET("/health", HealthCheck)

	// --- Agent 相关的 API 路由组 ---
	agentGroup := router.Group("/api/v1/agent")
	{
		agentGroup.POST("/register", RegisterAgent)
		agentGroup.POST("/heartbeat", Heartbeat)
		agentGroup.GET("/tasks", GetTasks) // 长轮询接口
		agentGroup.POST("/tasks/results", PostTaskResults)
	}

	// --- 内部测试用的 API 路由组 ---
	internalGroup := router.Group("/api/v1/internal")
	{
		internalGroup.POST("/trigger_kb", TriggerKB)
	}

	return router
}
