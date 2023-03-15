package configs

import "time"

type DBConfig struct {
	Dsn           string        `mapstructure:"DSN" validate:"required"`
	TimeoutRead   time.Duration `mapstructure:"TIMEOUT_READ" validate:"required"`
	TimeoutDbconn time.Duration `mapstructure:"TIMEOUT_DBCONN" validate:"required"`
}
type ServerConfig struct {
	Port       int    `mapstructure:"SERVER_PORT" validate:"required"`
	JWT_SECRET string `mapstructure:"JWT_SECRET" validate:"required"`
}
