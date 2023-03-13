package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type AuthConfig struct {
	BaseDBConfig     DBConfig
	BaseServerConfig ServerConfig
}

func NewAuthConfig() (*AuthConfig, error) {
	var c AuthConfig

	viper.AddConfigPath("../../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	var dbConf DBConfig
	if err := viper.Unmarshal(&dbConf); err != nil {
		return nil, fmt.Errorf("NewConfiguration.UnmarshalDB %s", err)
	}

	var serverConf ServerConfig
	if err := viper.Unmarshal(&serverConf); err != nil {
		return nil, fmt.Errorf("NewConfiguration.UnmarshalServer %s", err)
	}

	c = AuthConfig{
		BaseDBConfig:     dbConf,
		BaseServerConfig: serverConf,
	}
	validate := validator.New()
	if err := validate.Struct(&c); err != nil {
		return nil, fmt.Errorf("NewConfiguration.Validate: Missing required attributes %v", err)
	}

	return &c, nil
}
