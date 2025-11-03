package model

import "gorm.io/gorm"

// Agent 代表一个终端智能体
type Agent struct {
	gorm.Model        // gorm.Model 内嵌了 ID, CreatedAt, UpdatedAt, DeletedAt 字段
	UUID       string `gorm:"uniqueIndex;not null"` // Agent 的唯一标识符
	Hostname   string
	IPAddress  string
	OS         string
	Status     string // 例如: "online", "offline"
}
