package server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `envconfig:"ADDR" default:":8080"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT" default:"10s"`
	AllowedOrigins  []string      `envconfig:"ALLOWED_ORIGINS" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

func LoadConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}
