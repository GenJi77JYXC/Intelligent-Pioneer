package model

import "time"

type Workflow struct {
	ID            string `gorm:"primaryKey"`
	KBID          string
	AgentID       string `gorm:"index"` // 为 agent_id 添加索引以加快查询
	Status        string
	CurrentTaskID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
