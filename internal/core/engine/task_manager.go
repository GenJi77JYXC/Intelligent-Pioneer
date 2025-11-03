package engine

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"sync"
	"time"
)

// TaskManager 负责管理和分发所有 Agent 的任务
type TaskManager struct {
	// key: agent_id, value: 一个用于传递任务的 channel
	agentTaskQueues map[string]chan *Task
	mu              sync.RWMutex // 用于保护 agentTaskQueues 的并发访问
}

// TM 是一个全局的任务管理器实例
var TM *TaskManager

// InitTaskManager 初始化全局的任务管理器
func InitTaskManager() {
	TM = &TaskManager{
		agentTaskQueues: make(map[string]chan *Task),
	}
	logger.L.Info("✅ Task Manager initialized successfully!")
}

// getOrCreateAgentQueue 获取或创建一个 Agent 的任务通道
func (tm *TaskManager) getOrCreateAgentQueue(agentID string) chan *Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	queue, exists := tm.agentTaskQueues[agentID]
	if !exists {
		// 创建一个带缓冲的 channel, 防止任务发送阻塞
		queue = make(chan *Task, 10)
		tm.agentTaskQueues[agentID] = queue
		logger.L.Infow("Created new task queue for agent", "agent_id", agentID)
	}
	return queue
}

// SubmitTask 向指定的 Agent 提交一个新任务
func (tm *TaskManager) SubmitTask(task *Task) {
	queue := tm.getOrCreateAgentQueue(task.AgentID)
	logger.L.Infow("Submitting new task to queue", "agent_id", task.AgentID, "task_id", task.ID, "command", task.Command)
	queue <- task
}

// GetTaskForAgent 为指定的 Agent 获取一个任务 (支持长轮询)
func (tm *TaskManager) GetTaskForAgent(agentID string, timeout time.Duration) *Task {
	queue := tm.getOrCreateAgentQueue(agentID)

	select {
	case task := <-queue:
		logger.L.Infow("Dispatched task to agent", "agent_id", agentID, "task_id", task.ID)
		return task
	case <-time.After(timeout):
		// 超时，没有任务
		return nil
	}
}

// TODO: 添加一个清理不活跃 Agent 队列的逻辑 (用于生产环境)
