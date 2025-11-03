package main

import (
	"context"
	"fmt"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/api"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/core/engine"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/mq"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/store"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("🚀 Starting Intelligent-Pioneer... The journey begins!")

	// 1. Load configuration (Viper)
	config.LoadConfig()

	// 为了验证配置是否加载成功，我们可以打印一些值
	fmt.Println("Server mode:", config.C.Server.Mode)
	fmt.Printf("PostgreSQL Host: %s, Port: %s\n", config.C.Database.Postgres.Host, config.C.Database.Postgres.Port)
	fmt.Println("Kafka Brokers:", config.C.Kafka.Brokers)
	// 2. Initialize logger (Zap/Logrus)
	logger.InitLogger()

	// 使用全局Logger打印日志
	logger.L.Info("🚀 Starting Intelligent-Pioneer... The journey begins!")
	logger.L.Debugw("Configuration loaded successfully.",
		"server_mode", config.C.Server.Mode,
		"postgres_host", config.C.Database.Postgres.Host,
	)
	// 3. Initialize database connections (PostgreSQL, Elasticsearch)
	store.InitPostgres()
	store.InitElasticsearch()
	// 4. Initialize message queue producer/consumer (Kafka)
	mq.InitKafka()
	// 示例：启动后发送一条测试消息
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	logger.L.Info("Sending a test message to Kafka...")
	//	err := mq.MetricProducer.WriteMessages(context.Background(),
	//		kafka.Message{
	//			Key:   []byte("test-key"),
	//			Value: []byte("{\"cpu_usage\": 10.5}"),
	//		},
	//	)
	//	if err != nil {
	//		logger.L.Errorw("Failed to send test message", "error", err)
	//	}
	//}()
	//
	//time.Sleep(20 * time.Second)

	// 初始化任务管理器/引擎
	engine.InitTaskManager()

	// 5. Initialize HTTP server (Gin) and register routes
	router := api.NewRouter()
	// 6. Start the server and wait for shutdown signal
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.C.Server.Port),
		Handler: router,
	}
	go func() {
		// 启动服务器
		logger.L.Infof("Server is listening on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatalf("Listen error: %s\n", err)
		}
	}()

	// ---- 优雅关停逻辑 ----
	// 创建一个 channel 来接收系统信号
	quit := make(chan os.Signal, 1)
	// 我们只关心 SIGINT 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞在此，直到接收到信号
	<-quit

	logger.L.Info("Shutting down server...")

	// 创建一个有超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用服务器的 Shutdown 方法
	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Fatalw("Server forced to shutdown:", "error", err)
	}

	logger.L.Info("Server exiting.")
}
