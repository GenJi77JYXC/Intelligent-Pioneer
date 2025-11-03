package logger

import (
	"github.com/GenJi77JYXC/intelligent-pioneer/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

// L 是一个全局的、可导出的 SugaredLogger
// 项目的其他部分可以通过 logger.L 来访问
var L *zap.SugaredLogger

// InitLogger 根据配置初始化全局的Logger
func InitLogger() {
	// 根据配置决定日志编码器 (console or json)
	encoder := getEncoder()

	// 根据配置决定日志写入位置 (stdout or file)
	writeSyncer := getWriteSyncer()

	// 根据配置决定日志级别
	logLevel := zapcore.DebugLevel
	if err := logLevel.UnmarshalText([]byte(config.C.Logger.Level)); err != nil {
		// 如果配置中的级别无效，则默认为 debug
		logLevel = zapcore.DebugLevel
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	// 创建 Logger
	// zap.AddCaller() 会在日志中记录调用者的文件名和行号
	logger := zap.New(core, zap.AddCaller())

	// 使用 SugaredLogger，它提供了更方便的 API (如 Printf, Infof)
	L = logger.Sugar()

	L.Info("✅ Logger initialized successfully!")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// 人类可读的时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 小写字母的日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if strings.ToLower(config.C.Logger.Format) == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	// 默认使用 console 格式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	output := strings.ToLower(config.C.Logger.Output)
	if output == "stderr" {
		return zapcore.AddSync(os.Stderr)
	}
	if strings.HasSuffix(output, ".log") {
		// TODO: 在生产环境中，这里应该使用日志轮转库，如 lumberjack
		// 为了MVP，我们先用简单的文件写入
		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// 如果文件打开失败，则回退到标准输出
			return zapcore.AddSync(os.Stdout)
		}
		return zapcore.AddSync(file)
	}
	// 默认输出到标准输出
	return zapcore.AddSync(os.Stdout)
}
