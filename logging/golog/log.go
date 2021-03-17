package golog

import (
	"log"

	"github.com/hellcats88/abstracte/logging"
	rem "github.com/hellcats88/rem/logging"
)

func writeWithStd(level logging.Level, msg string) {
	log.Printf("%s\n", msg)
}

// NewLoggerConsole creates new logger based on Console implementation
func New(config logging.Config) logging.Logger {
	return rem.New(config, writeWithStd)
}

// LoggerConsoleFactory is a factory based on Console logger
type factory struct {
	config logging.Config
}

// NewLoggerConsoleFactory creates a new factory
func NewFactory(config logging.Config) logging.Factory {
	return factory{config: config}
}

func (f factory) Create() logging.Logger {
	return New(f.config)
}

func (f factory) CreateFromConfig(config logging.Config) logging.Logger {
	return New(config)
}
