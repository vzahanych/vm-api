package core

import (
	"time"
)

type (
	Config struct {
		Server        ServerConfig        `yaml:"server"`
		Opentelemetry OpentelemetryConfig `yaml:"opentelemetry"`
		Healthcheck   HealthcheckConfig   `yaml:"healthcheck"`
		Logging       LoggingConfig       `yaml:"logging"` // Added logging section
	}

	ServerConfig struct {
		Address               string        `yaml:"address"`
		ReadTimeout           time.Duration `yaml:"read_timeout"`
		WriteTimeout          time.Duration `yaml:"write_timeout"`
		GracefulShutdownDelay time.Duration `yaml:"graceful_shutdown_delay"`
	}

	OpentelemetryConfig struct {
		Enabled     bool   `yaml:"enabled"`
		Endpoint    string `yaml:"endpoint"`
		ServiceName string `yaml:"service_name"`
	}

	HealthcheckConfig struct {
		Enabled bool   `yaml:"enabled"`
		URL     string `yaml:"url"`
	}

	LoggingConfig struct {
		LogLevel string `yaml:"log_level"` // Log level (debug, info, warn, error)
	}
)
