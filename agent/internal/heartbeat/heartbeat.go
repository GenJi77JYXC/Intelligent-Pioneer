package heartbeat

import (
	"context"
	"log"
	"time"

	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/client"
)

// Start 启动心跳循环
func Start(ctx context.Context, apiClient *client.APIClient, agentID string) {
	log.Println("Heartbeat service started.")
	ticker := time.NewTicker(1 * time.Minute) // 每分钟发送一次心跳
	defer ticker.Stop()

	// 立即发送一次心跳，而不是等一分钟
	if err := apiClient.SendHeartbeat(agentID); err != nil {
		log.Printf("ERROR: Failed to send initial heartbeat: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := apiClient.SendHeartbeat(agentID); err != nil {
				log.Printf("ERROR: Failed to send heartbeat: %v", err)
			} else {
				log.Println("Heartbeat sent successfully.")
			}
		case <-ctx.Done():
			log.Println("Heartbeat service stopped.")
			return
		}
	}
}
