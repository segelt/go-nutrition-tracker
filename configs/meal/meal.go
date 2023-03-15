package meal

import (
	"fmt"
	"nutritiontracker/configs"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type MealDBConfig struct{}

type MealServerConfig struct{}

type MealConfig struct {
	BaseDBConfig     configs.DBConfig
	BaseServerConfig configs.ServerConfig
}

func NewMealConfig() (*MealConfig, error) {
	var c MealConfig

	viper.AddConfigPath("../../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	var dbConf configs.DBConfig
	if err := viper.Unmarshal(&dbConf); err != nil {
		return nil, fmt.Errorf("NewConfiguration.UnmarshalDB %s", err)
	}

	var serverConf configs.ServerConfig
	if err := viper.Unmarshal(&serverConf); err != nil {
		return nil, fmt.Errorf("NewConfiguration.UnmarshalServer %s", err)
	}

	c = MealConfig{
		BaseDBConfig:     dbConf,
		BaseServerConfig: serverConf,
	}
	validate := validator.New()
	if err := validate.Struct(&c); err != nil {
		return nil, fmt.Errorf("NewConfiguration.Validate: Missing required attributes %v", err)
	}

	return &c, nil
}
