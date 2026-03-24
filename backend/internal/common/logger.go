package common

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger

func init() {
	// 默认初始化，防止未调用 InitLogger 时 Logger 为 nil
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	Logger = logger.Sugar()
}

// LogConfig 日志配置接口
type LogConfig interface {
	GetLevel() string
	GetFilename() string
	GetMaxSize() int
	GetMaxBackups() int
	GetMaxAge() int
	GetCompress() bool
}

// InitLogger 根据配置初始化日志
func InitLogger(cfg LogConfig) {
	writeSyncer := getLogWriter(
		cfg.GetFilename(),
		cfg.GetMaxSize(),
		cfg.GetMaxBackups(),
		cfg.GetMaxAge(),
		cfg.GetCompress(),
	)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.GetLevel()))
	if err != nil {
		*l = zapcore.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writeSyncer, l),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), l),
	)

	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func FormatDataForDisplay(results []map[string]interface{}) []map[string]interface{} {
	return results
}
