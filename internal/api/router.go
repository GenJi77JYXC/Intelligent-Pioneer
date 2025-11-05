package api

import (
	"time"

	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"github.com/gin-contrib/cors"
)
import "github.com/gin-gonic/gin"

// NewRouter 创建并配置一个新的 Gin 引擎
func NewRouter() *gin.Engine {
	// 根据配置设置 Gin 的模式 (debug/release)
	gin.SetMode(config.C.Server.Mode)

	// 创建一个不带任何默认中间件的 Gin 引擎
	// 如果需要 Logger 和 Recovery 中间件，可以使用 gin.Default()
	router := gin.New()

	// 配置并使用 CORS 中间件
	// router.Use() 用于注册全局中间件，它会对所有请求生效
	router.Use(cors.New(cors.Config{
		// AllowOrigins 是允许跨域请求的源地址列表。
		// 在开发环境中，我们通常允许来自 Vite 开发服务器的地址。
		// 为了方便，可以使用 "*" 允许所有源，但在生产环境中应该指定明确的域名。
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},

		// AllowMethods 是允许的 HTTP 方法列表。
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},

		// AllowHeaders 是允许的请求头列表。
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},

		// ExposeHeaders 是允许客户端访问的响应头列表。
		ExposeHeaders: []string{"Content-Length"},

		// AllowCredentials 指示是否允许发送 Cookie。
		// 如果你的前端需要发送凭证（如 session cookie 或 JWT in cookie），请设置为 true。
		AllowCredentials: true,

		// MaxAge 指定了预检请求 (OPTIONS) 的结果可以被缓存多久。
		MaxAge: 12 * time.Hour,
	}))
	// --- CORS 配置结束 ---

	// 在此之后可以添加其他全局中间件，如 Logger 和 Recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 在这里可以添加全局中间件
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	// router.Use(middlewares.Cors()) // 例如：跨域中间件

	// 健康检查路由
	router.GET("/health", HealthCheck)

	// --- Agent 相关的 API 路由组 ---
	agentGroup := router.Group("/api/v1/agent")
	{
		agentGroup.GET("", GetAllAgents)
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
