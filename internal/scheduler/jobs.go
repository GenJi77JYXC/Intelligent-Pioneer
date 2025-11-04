package scheduler

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/model"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/store"
	"time"
)

// CheckOfflineAgents 是一个定时任务，用于检测并标记离线的 Agent
func CheckOfflineAgents() {
	logger.L.Debug("Running job: CheckOfflineAgents")

	// 解析配置中的超时时间
	timeoutDuration, err := time.ParseDuration(config.C.Agent.HeartbeatTimeout)
	if err != nil {
		logger.L.Errorw("Invalid heartbeat_timeout duration in config, using default 5m", "error", err)
		timeoutDuration = 5 * time.Minute
	}

	// 计算超时的临界时间点
	// 任何 'online' 的 Agent，如果其 updated_at 在这个时间点之前，就认为它离线了
	deadline := time.Now().Add(-timeoutDuration)

	// 构造更新条件
	// 我们使用 GORM 的原生 SQL 功能来执行一个更高效的批量更新
	// 这比“查询出来再逐个更新”要快得多
	result := store.DB.Model(&model.Agent{}).
		Where("status = ? AND updated_at < ?", "online", deadline).
		Update("status", "offline")

	if result.Error != nil {
		logger.L.Errorw("Error checking for offline agents", "error", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		logger.L.Infow("Marked agents as offline", "count", result.RowsAffected)
	} else {
		logger.L.Debug("No offline agents found.")
	}
}
