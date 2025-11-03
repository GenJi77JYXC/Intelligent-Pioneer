package api

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/core/engine"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/model"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// RegisterAgent 处理 Agent 注册请求
func RegisterAgent(c *gin.Context) {
	// 1. 绑定并校验请求参数
	var req RegisterAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warnw("Invalid agent registration request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	logger.L.Infow("Processing agent registration request",
		"hostname", req.Hostname,
		"ip_address", req.IPAddress,
		"os", req.OS,
	)

	// 2. 检查该 Agent 是否已经注册过 (基于某些唯一标识，例如 Hostname + IP)
	// 这是一个可选的健壮性设计，MVP 阶段可以先简化
	var existingAgent model.Agent
	result := store.DB.Where("hostname = ? AND ip_address = ?", req.Hostname, req.IPAddress).First(&existingAgent)

	if result.Error == nil {
		// 如果找到了记录，说明已经注册过，直接返回已有的 ID
		logger.L.Infow("Agent already registered, returning existing ID", "agent_id", existingAgent.UUID)
		c.JSON(http.StatusOK, RegisterAgentResponse{
			AgentID: existingAgent.UUID,
			Message: "Agent already registered.",
		})
		return
	}
	// 如果错误不是 "record not found"，说明是其他数据库错误
	if result.Error != gorm.ErrRecordNotFound {
		logger.L.Errorw("Database error while checking for existing agent", "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 3. 创建新的 Agent 记录
	newAgent := model.Agent{
		UUID:      uuid.NewString(), // 生成一个新的唯一ID
		Hostname:  req.Hostname,
		IPAddress: req.IPAddress,
		OS:        req.OS,
		Status:    "offline", // 初始状态为离线，等待心跳
	}
	newAgent.CreatedAt = time.Now() // 手动设置时间或让 GORM 自动处理
	newAgent.UpdatedAt = time.Now()

	// 4. 将新记录存入 PostgreSQL
	createResult := store.DB.Create(&newAgent)
	if createResult.Error != nil {
		logger.L.Errorw("Failed to create new agent in database", "error", createResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register agent"})
		return
	}

	logger.L.Infow("New agent registered successfully", "agent_id", newAgent.UUID)

	// 5. 返回成功响应和新的 agent_id
	c.JSON(http.StatusCreated, RegisterAgentResponse{
		AgentID: newAgent.UUID,
		Message: "Agent registered successfully.",
	})
}

// Heartbeat 处理 Agent 心跳请求
func Heartbeat(c *gin.Context) {
	// 1. 绑定并校验请求参数
	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 对于频繁调用的心跳接口，可以简化日志，避免日志泛滥
		// logger.L.Warnw("Invalid heartbeat request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: agent_id is required."})
		return
	}

	// 2. 在数据库中更新 Agent 状态和时间戳
	// 我们使用 GORM 的 Updates 方法，它只会更新指定的字段，效率更高
	// 并且我们只更新 UUID 匹配的记录
	updateData := map[string]interface{}{
		"status":     "online",
		"updated_at": time.Now(), // GORM 会自动处理 updated_at, 但手动更新更明确
	}

	result := store.DB.Model(&model.Agent{}).Where("uuid = ?", req.AgentID).Updates(updateData)

	// 3. 检查更新操作的结果
	if result.Error != nil {
		logger.L.Errorw("Failed to update agent heartbeat in database", "agent_id", req.AgentID, "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// RowsAffected 返回受影响的行数。如果为 0，说明没有找到对应的 Agent
	if result.RowsAffected == 0 {
		logger.L.Warnw("Heartbeat received from an unknown or unregistered agent", "agent_id", req.AgentID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found. Please register first."})
		return
	}

	// 对于心跳这种高频请求，成功时可以不打印日志，以保持日志清爽
	// logger.L.Debugw("Heartbeat updated for agent", "agent_id", req.AgentID)

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	logger.L.Info("Received agent heartbeat")
}

// GetTasks 处理 Agent 获取任务的长轮询请求
func GetTasks(c *gin.Context) {
	// 1. 从查询参数中获取 agent_id
	agentID := c.Query("agent_id")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'agent_id' is required."})
		return
	}

	// 对于长轮询，不建议打印太多开始日志，可以在返回时打印
	// logger.L.Debugw("Agent polling for tasks...", "agent_id", agentID)

	// 2. 调用任务管理器的长轮询方法
	// 我们在这里设置一个 30 秒的超时时间
	timeout := 30 * time.Second
	task := engine.TM.GetTaskForAgent(agentID, timeout)

	// 3. 根据结果返回响应
	if task == nil {
		// 如果返回 nil，说明是超时，没有任务
		// 返回 204 No Content 是一个很好的实践，表示请求成功，但没有内容返回
		logger.L.Debugw("Polling timeout, no tasks for agent", "agent_id", agentID)
		c.Status(http.StatusNoContent)
	} else {
		// 如果获取到了任务，将其序列化为 JSON 返回
		logger.L.Infow("Dispatched task to agent via polling", "agent_id", agentID, "task_id", task.ID)
		c.JSON(http.StatusOK, task)
	}
	logger.L.Info("Agent polling for tasks...")
	c.JSON(http.StatusOK, gin.H{"tasks": []string{}}) // 暂时返回空任务列表
}

// PostTaskResults 处理 Agent 上报任务结果的请求
func PostTaskResults(c *gin.Context) {
	// 1. 绑定并校验请求参数
	var result engine.TaskResult
	if err := c.ShouldBindJSON(&result); err != nil {
		logger.L.Warnw("Invalid task result submission", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 简单的校验，确保核心字段不为空
	if result.TaskID == "" || result.AgentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id and agent_id are required."})
		return
	}

	logger.L.Infow("Received task result from agent",
		"agent_id", result.AgentID,
		"task_id", result.TaskID,
		"success", result.Success,
	)

	// 2. 将结果异步地交给核心引擎处理
	// 为什么是异步？因为结果的处理（分析、下发新任务）可能需要时间，
	// 我们不应该让 Agent 在这里长时间等待。Agent 只需要知道我们收到了结果即可。
	// 所以我们把它放进一个新的 goroutine 中执行。
	go engine.HandleTaskResult(&result)

	// 3. 立即返回成功响应给 Agent
	// 这告诉 Agent：“答卷已收到，你可以去领下一份卷子了（再次调用 GetTasks）”。
	c.JSON(http.StatusOK, gin.H{"status": "result received and is being processed"})
	logger.L.Info("Received task results from agent")
}
