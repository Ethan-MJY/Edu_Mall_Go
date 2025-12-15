// Package logger 日志工具模块
// 职责: 封装zap日志库,提供统一的日志输出接口
// 特性: 支持动态调整日志级别,JSON格式输出
package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var atom = zap.NewAtomicLevelAt(zap.DebugLevel) // 原子日志级别,支持运行时动态修改

// init 初始化zap logger
// 配置: JSON编码,输出到stdout和stderr
// 级别: 默认Debug,可通过SetLevel动态调整
func init() {
	config := zap.Config{
		Level:       atom, // 原子级别,支持动态修改
		Development: false,
		Encoding:    "json", // JSON编码,便于日志收集和分析
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "level",
			TimeKey:    "time",
			CallerKey:  "caller",
			EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	tempLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// 添加调用者信息,跳过1层(logger本身),错误级别自动记录堆栈
	logger = tempLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
}

// SetLevel 动态设置日志级别
// 参数: level 日志级别字符串,如"debug"/"info"/"warn"/"error"
// 调用: main.main -> SetLevel
func SetLevel(level string) {
	tLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		fmt.Printf("invalid level, input: %s", level)
		return
	}
	atom.SetLevel(tLevel)
}

// Debug 输出Debug级别日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info 输出Info级别日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 输出Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 输出Error级别日志
// 自动附加堆栈信息
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
