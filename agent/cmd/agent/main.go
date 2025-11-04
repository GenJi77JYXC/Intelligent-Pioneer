package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/client"
	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/config"
	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/heartbeat"
	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/sysinfo"
	"github.com/GenJi77JYXC/intelligent-pioneer/agent/internal/task"
)

func main() {
	log.Println("Starting Intelligent-Pioneer Agent...")

	// 1. 加载配置
	// 我们将配置文件放在用户主目录下的一个隐藏文件夹里
	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/.intelligent-pioneer"
	if err := config.LoadConfig(configDir); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化 API 客户端
	apiClient := client.NewAPIClient(config.Cfg.BackendURL)

	// 3. 注册 (如果需要)
	if config.Cfg.AgentID == "" {
		log.Println("Agent not registered. Attempting to register...")
		// 采集本机信息
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("WARN: Could not get hostname: %v", err)
			hostname = "unknown-host"
		}

		ip, err := sysinfo.GetPrimaryIP()
		if err != nil {
			log.Printf("WARN: Could not get primary IP: %v", err)
			ip = "0.0.0.0" // 使用一个默认值
		}

		osInfo, err := sysinfo.GetOSInfo()
		if err != nil {
			log.Printf("WARN: Could not get OS info: %v", err)
			osInfo = "unknown-os"
		}

		agentID, err := apiClient.Register(hostname, ip, osInfo)
		if err != nil {
			log.Fatalf("Failed to register agent: %v", err)
		}
		config.Cfg.AgentID = agentID
		if err := config.SaveConfig(configDir); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
		log.Printf("Agent registered successfully with ID: %s", agentID)
	} else {
		log.Printf("Agent already registered with ID: %s", config.Cfg.AgentID)
	}

	// 4. 创建一个可以被取消的 context，用于优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5. 在 goroutine 中启动心跳和任务拉取
	go heartbeat.Start(ctx, apiClient, config.Cfg.AgentID)
	go task.StartPolling(ctx, apiClient, config.Cfg.AgentID)

	log.Println("Agent is running. Press Ctrl+C to exit.")

	// 6. 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received. Exiting gracefully...")
	// cancel() 会通知所有使用 ctx 的 goroutine 退出
}
