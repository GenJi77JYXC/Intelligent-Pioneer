package store

import (
	"fmt"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 是一个全局的、可导出的数据库连接实例
var DB *gorm.DB

// InitPostgres 根据配置初始化 PostgreSQL 连接
func InitPostgres() {
	logger.L.Info("Initializing PostgreSQL connection...")

	// 从全局配置中获取数据库配置
	pgConfig := config.C.Database.Postgres

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		pgConfig.Host,
		pgConfig.User,
		pgConfig.Password,
		pgConfig.DBName,
		pgConfig.Port,
		pgConfig.SSLMode,
	)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.L.Fatalw("Failed to connect to PostgreSQL database", "error", err)
	}

	// 将连接实例赋值给全局变量
	DB = db
	logger.L.Info("✅ PostgreSQL connection established successfully!")

	// 自动迁移模式：自动创建或更新表结构
	migrateDatabase()
}

func migrateDatabase() {
	logger.L.Info("Running database migrations...")
	// GORM 会自动检查 `Agent` 结构体对应的表是否存在，不存在则创建
	err := DB.AutoMigrate(
		&model.Agent{},
		&model.Workflow{},
	)
	if err != nil {
		logger.L.Fatalw("Failed to migrate database", "error", err)
	}
	logger.L.Info("✅ Database migration completed.")
}
