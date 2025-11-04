package scheduler

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

// InitScheduler 初始化并启动后台任务调度器
func InitScheduler() {
	logger.L.Info("Initializing scheduler...")

	// 创建一个新的 cron 调度器，支持秒级精度
	c = cron.New(cron.WithSeconds())

	// 注册我们的离线检测任务
	_, err := c.AddFunc(config.C.Agent.OfflineCheckCron, CheckOfflineAgents)
	if err != nil {
		logger.L.Fatalw("Failed to add offline agent check job to scheduler", "error", err)
	}

	// 在一个新的 goroutine 中启动调度器，避免阻塞主线程
	go c.Start()

	logger.L.Info("✅ Scheduler started successfully!")
}

// StopScheduler 优雅地停止调度器 (在程序退出时调用)
func StopScheduler() {
	logger.L.Info("Stopping scheduler...")
	ctx := c.Stop() // Stop 会等待所有正在运行的任务完成
	<-ctx.Done()
	logger.L.Info("Scheduler stopped gracefully.")
}
