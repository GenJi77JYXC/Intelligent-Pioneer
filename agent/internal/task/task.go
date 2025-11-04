package task

import (
	"context"
	"log"
	"time"

	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/client"
	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/executor"
)

// StartPolling 启动任务拉取循环
func StartPolling(ctx context.Context, apiClient *client.APIClient, agentID string) {
	log.Println("Task polling service started.")
	for {
		select {
		case <-ctx.Done():
			log.Println("Task polling service stopped.")
			return
		default:
			log.Println("Polling for new tasks...")
			task, err := apiClient.FetchTasks(ctx, agentID)
			if err != nil {
				if ctx.Err() == context.Canceled {
					// 如果是主动取消，则正常退出
					return
				}
				log.Printf("ERROR: Failed to fetch tasks: %v. Retrying in 10s...", err)
				time.Sleep(10 * time.Second)
				continue
			}

			if task != nil {
				log.Printf("New task received: ID=%s, Command=%s", task.ID, task.Command)
				// 异步执行任务，避免阻塞任务拉取循环
				go func(t *client.Task) {
					result := executor.Execute(agentID, t)
					if err := apiClient.PostResult(result); err != nil {
						log.Printf("ERROR: Failed to post task result: %v", err)
					} else {
						log.Printf("Task result posted successfully: ID=%s", t.ID)
					}
				}(task)
			}
		}
	}
}
