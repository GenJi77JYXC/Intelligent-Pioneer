# Gin HTTP 服务器初始化与路由配置指南

本文档详细介绍了如何在 `Intelligent-Pioneer` 项目中集成 [Gin](https://github.com/gin-gonic/gin) Web 框架，创建 API 路由，并实现一个支持优雅关停（Graceful Shutdown）的生产级 HTTP 服务器。

## 目标

完成本项目后端服务的基础架构搭建，使其能够：
1.  接收并处理来自外部（终端 Agent、Web 前端）的 HTTP 请求。
2.  拥有清晰、可扩展的 API 路由结构。
3.  在服务停止时能够优雅地处理完当前请求，确保数据一致性。

## 1. 安装 Gin 框架

首先，确保 Gin 框架已经被添加到项目依赖中。如果尚未安装，请在项目根目录执行：

```bash
go get github.com/gin-gonic/gin
```

## 2. 构建 API 模块 (`internal/api`)

为了保持代码的组织性和可维护性，我们将所有与 API 相关的代码（路由定义、请求处理器）都放在 `internal/api` 包下。

### 2.1. 定义路由 (`router.go`)

创建一个 `router.go` 文件来集中管理所有的 API 路由。使用 Gin 的路由组（`Group`）功能可以有效地组织不同模块的 API。

**文件路径:** `internal/api/router.go`

```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/heyang-code/intelligent-pioneer/internal/config"
)

// NewRouter 创建并配置一个新的 Gin 引擎
func NewRouter() *gin.Engine {
	// 根据配置设置 Gin 的模式 (debug/release)
	gin.SetMode(config.C.Server.Mode)

	// 创建一个不带默认中间件的引擎，以便更精细地控制
	router := gin.New()

	// 可以在此添加自定义的全局中间件
	// 例如: 日志、恢复、跨域等
	// router.Use(gin.Logger(), gin.Recovery(), middlewares.Cors())

	// --- 基础路由 ---
	router.GET("/health", HealthCheck)

	// --- Agent API 路由组 (v1) ---
	agentGroup := router.Group("/api/v1/agent")
	{
		agentGroup.POST("/register", RegisterAgent)
		agentGroup.POST("/heartbeat", Heartbeat)
		agentGroup.GET("/tasks", GetTasks) // 用于 Agent 长轮询任务
		agentGroup.POST("/tasks/results", PostTaskResults)
	}

	// --- 内部/测试 API 路由组 (v1) ---
	internalGroup := router.Group("/api/v1/internal")
	{
		internalGroup.POST("/trigger_kb", TriggerKB) // 手动触发知识库任务
	}

	return router
}
```

### 2.2. 实现请求处理器 (Handlers)

为每个路由创建对应的处理函数。为了保持代码整洁，我们将不同模块的 Handler 放入不同的文件中。

#### Agent Handlers (`handler_agent.go`)

**文件路径:** `internal/api/handler_agent.go`

```go
package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/heyang-code/intelligent-pioneer/internal/logger"
)

// RegisterAgent 处理 Agent 注册请求 (目前为模拟实现)
func RegisterAgent(c *gin.Context) {
	logger.L.Info("Received agent registration request")
	// TODO: 实现真正的注册逻辑
	c.JSON(http.StatusOK, gin.H{
		"message":   "Agent registered successfully (mock)",
		"agent_id": "mock-uuid-12345",
	})
}

// Heartbeat 处理 Agent 心跳请求 (目前为模拟实现)
func Heartbeat(c *gin.Context) {
	logger.L.Info("Received agent heartbeat")
	// TODO: 实现真正的心跳逻辑
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ... 其他 Agent 相关的 Handler ...
```
*(注：为保持文档简洁，此处省略了 `GetTasks` 和 `PostTaskResults` 的模拟代码)*

#### Internal Handlers (`handler_internal.go`)

**文件路径:** `internal/api/handler_internal.go`

```go
package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/heyang-code/intelligent-pioneer/internal/logger"
)

// TriggerKB 手动触发知识库任务 (目前为模拟实现)
func TriggerKB(c *gin.Context) {
	logger.L.Info("Manual KB trigger received")
	// TODO: 实现真正的触发逻辑
	c.JSON(http.StatusOK, gin.H{"message": "KB task triggered successfully (mock)"})
}

// HealthCheck 提供服务健康状态检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
```

## 3. 在主程序中启动服务器 (`main.go`)

最后，我们在 `main.go` 中将所有部分整合起来，初始化路由并启动一个支持优雅关停的 HTTP 服务器。

**文件路径:** `cmd/intelligent-pioneer/main.go`

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heyang-code/intelligent-pioneer/internal/api"
	"github.com/heyang-code/intelligent-pioneer/internal/config"
	"github.com/heyang-code/intelligent-pioneer/internal/logger"
	"github.com/heyang-code/intelligent-pioneer/internal/mq"
	"github.com/heyang-code/intelligent-pioneer/internal/store"
)

func main() {
	// --- 初始化所有依赖服务 ---
	config.LoadConfig()
	logger.InitLogger()
	store.InitPostgres()
	store.InitElasticsearch()
	mq.InitKafka()

	logger.L.Info("🚀 All services initialized. Starting HTTP server...")

	// 1. 初始化 Gin 路由
	router := api.NewRouter()

	// 2. 配置并启动 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.C.Server.Port),
		Handler: router,
	}

	go func() {
		// 在一个 goroutine 中启动服务器，避免阻塞主线程
		logger.L.Infof("Server is listening on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatalf("Listen error: %s\n", err)
		}
	}()

	// 3. 实现优雅关停 (Graceful Shutdown)
	// 创建一个 channel 来等待系统中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞，直到接收到信号

	logger.L.Info("Shutting down server...")

	// 创建一个5秒超时的 context，用于通知服务器有5秒时间来处理剩余请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用 Shutdown() 方法来优雅地关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Fatalw("Server forced to shutdown:", "error", err)
	}

	logger.L.Info("Server exiting.")
}
```

## 4. 验证

1.  **启动所有依赖服务：**
    ```bash
    make docker-up
    ```
2.  **运行后端应用：**
    ```bash
    make run
    ```
    你应该会在日志中看到服务器已在 `8080` 端口上监听。

3.  **测试 API 端点：**
    *   **健康检查:**
        ```bash
        curl http://localhost:8080/health
        # 预期输出: {"status":"UP"}
        ```
    *   **模拟 Agent 注册:**
        ```bash
        curl -X POST http://localhost:8080/api/v1/agent/register
        # 预期输出: {"agent_id":"mock-uuid-12345","message":"Agent registered successfully (mock)"}
        ```
4.  **测试优雅关停：**
    在运行应用的终端按下 `Ctrl+C`。应用会打印关停日志并等待几秒后退出，而不是立即崩溃。

---

至此，`Intelligent-Pioneer` 项目的后端基础架构已全部搭建完成。我们拥有了一个结构清晰、功能完备、具备生产级特性的应用框架，为后续的业务功能开发奠定了坚实的基础。