package logger

import "github.com/phuslu/log"

// PluginLogger 插件日志记录器
var PluginLogger = log.Logger{
	Level:  log.InfoLevel, // 信息日志级别
	Caller: 1,
	Writer: &log.ConsoleWriter{
		ColorOutput:    true,
		EndWithMessage: true,
	},
}
