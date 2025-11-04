package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = "agent_config.json"

// Config 定义了 Agent 的配置结构
type Config struct {
	BackendURL string `json:"backend_url"`
	AgentID    string `json:"agent_id"`
	// 未来可以添加更多配置, 如日志级别等
}

var Cfg *Config

// LoadConfig 从文件加载配置，如果文件不存在则使用默认值
func LoadConfig(configDir string) error {
	Cfg = &Config{
		// 设置一个默认的后端地址
		BackendURL: "http://localhost:8080",
	}

	configFile := filepath.Join(configDir, ConfigFileName)
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在是正常情况，直接使用默认配置
			return nil
		}
		return err // 其他读取错误
	}

	return json.Unmarshal(data, Cfg)
}

// SaveConfig 将当前配置保存到文件
func SaveConfig(configDir string) error {
	if Cfg == nil {
		return nil // 没有配置可保存
	}

	configFile := filepath.Join(configDir, ConfigFileName)
	data, err := json.MarshalIndent(Cfg, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}
