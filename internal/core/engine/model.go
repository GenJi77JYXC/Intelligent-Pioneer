package engine

import "time"

// Task 代表一个需要 Agent 执行的具体指令
type Task struct {
	ID         string    `json:"ID"`      // 任务的唯一ID
	AgentID    string    `json:"AgentID"` // 目标 Agent
	WorkflowID string    `json:"WorkflowID"`
	Type       string    `json:"Type"`      // 任务类型, e.g., "diagnostic", "remediation"
	Command    string    `json:"Command"`   // 要执行的命令
	CreatedAt  time.Time `json:"CreatedAt"` // 创建时间
}

// TaskResult 代表 Agent 执行任务后返回的结果
type TaskResult struct {
	TaskID   string `json:"task_id"`
	AgentID  string `json:"agent_id"`
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	Error    string `json:"error"`
	ExitCode int    `json:"exit_code"`
}

// Workflow 代表一个完整的自动化工作流实例
// 我们可以把它存到数据库里，用于追踪状态
type Workflow struct {
	ID            string `gorm:"primaryKey"` // 使用 UUID 作为主键
	KBID          string
	AgentID       string
	Status        string // "pending(待处理)", "diagnosing(诊断中)", "remediating(修复中)", "completed(已完成)", "failed(已失败)"
	CurrentTaskID string // 当前正在执行的任务ID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type KnowledgeBaseItem struct {
	Diagnostics   []map[string]string `json:"diagnostics"`
	AnalysisLogic string              `json:"analysis_logic"`
	Remediation   map[string]string   `json:"remediation"`
}
