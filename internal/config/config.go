package config

import "fmt"
import "github.com/spf13/viper"

// Config 是整个应用程序的配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Agent    AgentConfig    `mapstructure:"agent"`
}

// ServerConfig 对应 server 部分的配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 对应 database 部分的配置
type DatabaseConfig struct {
	Postgres      PostgresConfig      `mapstructure:"postgres"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
}

// PostgresConfig 对应 postgres 数据库的配置
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// ElasticsearchConfig 对应 elasticsearch 的配置
type ElasticsearchConfig struct {
	Addresses []string `mapstructure:"addresses"`
}

// KafkaConfig 对应 kafka 的配置
type KafkaConfig struct {
	Brokers []string          `mapstructure:"brokers"`
	Topics  KafkaTopicsConfig `mapstructure:"topics"`
}

// KafkaTopicsConfig 对应 kafka.topics 的配置
type KafkaTopicsConfig struct {
	AgentMetrics string `mapstructure:"agent_metrics"`
	AgentLogs    string `mapstructure:"agent_logs"`
}

// LoggerConfig 对应 logger 部分的配置
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// AgentConfig 对应 agent 部分的配置
type AgentConfig struct {
	HeartbeatTimeout string `mapstructure:"heartbeat_timeout"`
	OfflineCheckCron string `mapstructure:"offline_check_cron"`
}

// C 是一个全局变量，用于存储加载后的配置
// 外部包可以通过 config.C 访问配置
var C *Config

// LoadConfig 从 configs/config.yaml 文件加载配置
func LoadConfig() {
	v := viper.New()

	// 1. 设置配置文件路径
	v.AddConfigPath("./configs") // . 表示从执行程序的当前目录开始查找
	// 2. 设置配置文件名 (不带扩展名)
	v.SetConfigName("config")
	// 3. 设置配置文件类型
	v.SetConfigType("yaml")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 将读取的配置反序列化到全局变量 C 中
	if err := v.Unmarshal(&C); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}

	fmt.Println("✅ Configuration loaded successfully!")
}
