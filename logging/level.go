package logging

import (
	"strings"

	"github.com/hellcats88/abstracte/logging"
)

// Atol converts a string into a logger level. Utility function
func Atol(level string) logging.Level {
	var logLevel logging.Level
	switch strings.ToLower(level) {
	case "error":
		logLevel = logging.Error
	case "info":
		logLevel = logging.Info
	case "warn":
		logLevel = logging.Warn
	case "debug":
		logLevel = logging.Debug
	case "trace":
		logLevel = logging.Trace
	default:
		logLevel = logging.Info
	}

	return logLevel
}
