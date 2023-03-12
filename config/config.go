package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type ConfDB struct {
	Dsn           string        `mapstructure:"DSN" validate:"required"`
	TimeoutRead   time.Duration `mapstructure:"TIMEOUT_READ" validate:"required"`
	TimeoutDbconn time.Duration `mapstructure:"TIMEOUT_DBCONN" validate:"required"`
}

type ConfServer struct {
	Port int `mapstructure:"SERVER_PORT" validate:"required"`
}

type Conf struct {
	DBConf     ConfDB
	ServerConf ConfServer
}

func New() (*Conf, error) {
	var c Conf

	viper.AddConfigPath("../../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	var dbConf ConfDB
	if err := viper.Unmarshal(&dbConf); err != nil {
		return nil, fmt.Errorf("NewConfiguration.UnmarshalDB %s", err)
	}

	var serverConf ConfServer
	if err := viper.Unmarshal(&serverConf); err != nil {
		return nil, fmt.Errorf("NewConfiguration.UnmarshalServer %s", err)
	}

	c = Conf{
		DBConf:     dbConf,
		ServerConf: serverConf,
	}
	validate := validator.New()
	if err := validate.Struct(&c); err != nil {
		return nil, fmt.Errorf("NewConfiguration.Validate: Missing required attributes %v", err)
	}

	return &c, nil
}
