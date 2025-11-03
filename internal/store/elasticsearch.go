package store

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/logger"
	"github.com/elastic/go-elasticsearch/v8"
)

// ESClient 是一个全局的、可导出的 Elasticsearch 客户端实例
var ESClient *elasticsearch.Client

// InitElasticsearch 根据配置初始化 Elasticsearch 连接
func InitElasticsearch() {
	logger.L.Info("Initializing Elasticsearch connection...")

	// 从全局配置中获取 ES 配置
	esConfig := config.C.Database.Elasticsearch

	cfg := elasticsearch.Config{
		Addresses: esConfig.Addresses,
		// 在生产环境中，这里可能需要配置用户名、密码、证书等
	}

	// 创建客户端
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.L.Fatalw("Failed to create Elasticsearch client", "error", err)
	}

	// 验证连接：发送一个 Info 请求到 ES 服务器
	res, err := client.Info()
	if err != nil {
		logger.L.Fatalw("Failed to get Elasticsearch info", "error", err)
	}
	defer res.Body.Close() // 确保关闭响应体

	if res.IsError() {
		logger.L.Fatalw("Elasticsearch server returned an error", "status", res.Status())
	}

	// 将客户端实例赋值给全局变量
	ESClient = client
	logger.L.Info("✅ Elasticsearch connection established successfully!", "version", elasticsearch.Version)
}
