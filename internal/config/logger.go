package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level      string `yaml:"level" env:"LOG_LEVEL" default:"info"`
	Format     string `yaml:"format" env:"LOG_FORMAT" default:"json"`
	OutputPath string `yaml:"output_path" env:"LOG_OUTPUT_PATH" default:"stdout"`
}

// NewLogger creates a structured logger based on configuration
func NewLogger(cfg LoggerConfig) (*zap.Logger, error) {
	var zapConfig zap.Config

	if cfg.Format == "json" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// Set log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Set output path
	if cfg.OutputPath != "" && cfg.OutputPath != "stdout" {
		zapConfig.OutputPaths = []string{cfg.OutputPath}
	}

	// Build logger
	logger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return logger, nil
}
