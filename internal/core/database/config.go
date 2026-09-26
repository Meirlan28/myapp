package database

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Database string `envconfig:"DB" default:"myapp"`
	SSLMode  string `envconfig:"SSLMODE" default:"disable"`

	MaxConns int32 `envconfig:"MAX_CONNS" default:"10"`
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

func LoadConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}
