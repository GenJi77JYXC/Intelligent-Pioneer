package executor

import (
	"context"
	"log"
	"os/exec"
	"time"

	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/client"
)

// Execute 执行一个任务并返回结果
func Execute(agentID string, task *client.Task) client.TaskResult {
	log.Printf("Executing command: %s", task.Command)

	// 设置命令执行超时，例如 1 分钟
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// 使用 sh -c 来执行命令，以便支持管道等 shell 特性
	cmd := exec.CommandContext(ctx, "sh", "-c", task.Command)

	output, err := cmd.CombinedOutput() // 合并 stdout 和 stderr

	result := client.TaskResult{
		TaskID:  task.ID,
		AgentID: agentID,
		Output:  string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1 // 表示不是正常的退出错误（如超时）
		}
		log.Printf("Command execution failed: %v", err)
	} else {
		result.Success = true
		result.ExitCode = 0
		log.Println("Command executed successfully.")
	}

	return result
}
