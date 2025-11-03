package api

// TriggerKBRequest 定义了手动触发知识库工作流的请求体结构
type TriggerKBRequest struct {
	AgentID string `json:"agent_id" binding:"required"` // agent_id 是必需的
	KBID    string `json:"kb_id" binding:"required"`    // kb_id 是必需的
}

// RegisterAgentRequest 定义了 Agent 注册的请求体结构
type RegisterAgentRequest struct {
	Hostname  string `json:"hostname" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	OS        string `json:"os" binding:"required"`
}

// RegisterAgentResponse 定义了 Agent 注册的响应体结构
type RegisterAgentResponse struct {
	AgentID string `json:"agent_id"`
	Message string `json:"message"`
}

// HeartbeatRequest 定义了 Agent 心跳的请求体结构
type HeartbeatRequest struct {
	AgentID string `json:"agent_id" binding:"required"`
}
