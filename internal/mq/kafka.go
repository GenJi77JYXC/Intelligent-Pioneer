package mq

import (
	"context"
	"errors"
	"fmt"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/segmentio/kafka-go"
	"time"
)

// MetricProducer 是一个全局的、用于发送指标数据的 Kafka 生产者
var MetricProducer *kafka.Writer

// LogProducer 是一个全局的、用于发送日志数据的 Kafka 生产者
var LogProducer *kafka.Writer

// InitKafka 初始化 Kafka 生产者和消费者
func InitKafka() {
	logger.L.Info("Initializing Kafka...")

	// 1. 确保 Topic 存在
	createTopics()

	// 2. 初始化生产者
	initProducers()

	// 3. 初始化并启动消费者 (在独立的 goroutine 中)
	go initConsumers()

	logger.L.Info("✅ Kafka initialized successfully!")
}

// createTopics 检查并创建所需的 Kafka Topic
func createTopics() {
	conn, err := kafka.Dial("tcp", config.C.Kafka.Brokers[0])
	if err != nil {
		logger.L.Fatalw("Failed to connect to Kafka for topic creation", "error", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		logger.L.Fatalw("Failed to get Kafka controller", "error", err)
	}
	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		logger.L.Fatalw("Failed to connect to Kafka controller", "error", err)
	}
	defer controllerConn.Close()

	topicsToCreate := []kafka.TopicConfig{
		{
			Topic:             config.C.Kafka.Topics.AgentMetrics,
			NumPartitions:     1, // MVP 阶段，1个分区足够
			ReplicationFactor: 1, // 因为我们是单节点 Kafka
		},
		{
			Topic:             config.C.Kafka.Topics.AgentLogs,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicsToCreate...)
	if err != nil {
		// 忽略 "Topic with this name already exists" 错误
		var ke kafka.Error
		if errors.As(err, &ke) && errors.Is(ke, kafka.TopicAlreadyExists) {
			logger.L.Infow("Topics already exist, skipping creation.", "topics", topicsToCreate)
		}
	} else {
		logger.L.Info("Kafka topics created or already exist.")
	}
}

// initProducers 初始化所有生产者
func initProducers() {
	MetricProducer = newProducer(config.C.Kafka.Topics.AgentMetrics)
	LogProducer = newProducer(config.C.Kafka.Topics.AgentLogs)
	logger.L.Info("Kafka producers initialized.")
}

// newProducer 创建一个新的 Kafka 生产者 (Writer)
func newProducer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(config.C.Kafka.Brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{}, // 负载均衡策略
		// 在生产环境中，可以配置更复杂的参数，如重试、超时等
	}
}

// initConsumers 初始化并启动所有消费者
func initConsumers() {
	// 在这里，我们为每个需要消费的 Topic 启动一个消费者 goroutine
	go consumeMetrics()
	// go consumeLogs() // 如果需要消费日志，也在这里启动
	logger.L.Info("Kafka consumers started.")
}

// consumeMetrics 是一个示例消费者，用于消费 Agent 的性能指标
func consumeMetrics() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.C.Kafka.Brokers,
		Topic:    config.C.Kafka.Topics.AgentMetrics,
		GroupID:  "pioneer-metrics-consumer-group", // 消费者组 ID
		MinBytes: 10e3,                             // 10KB
		MaxBytes: 10e6,                             // 10MB
		MaxWait:  1 * time.Second,
	})
	defer reader.Close()

	logger.L.Infof("Consumer for topic '%s' is running...", config.C.Kafka.Topics.AgentMetrics)

	for {
		// 阻塞式读取消息
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.L.Errorw("Error reading message from Kafka", "topic", config.C.Kafka.Topics.AgentMetrics, "error", err)
			continue // or break
		}

		// TODO: 在这里处理消息
		// 例如: 将消息解析后存入 InfluxDB (时序数据库)
		logger.L.Infow("Received message from Kafka",
			"topic", m.Topic,
			"partition", m.Partition,
			"offset", m.Offset,
			"key", string(m.Key),
			"value", string(m.Value),
		)
	}
}
