package logging

import (
	"go.uber.org/zap"
)

var Debug bool

// DebugLog prints debug messages if Debug is enabled.
func DebugLog(format string, v ...interface{}) {
	if Debug {
		logger, _ := zap.NewProduction()
		logger.Debug(format, zap.Any("args", v))
	}
}
