package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	BaseConfig
}

var cfg *AppConfig

func Load() *AppConfig {
	base := &AppConfig{}
	viper.GetViper().AddConfigPath("../../")
	viper.GetViper().AddConfigPath("../../../")
	viper.GetViper().AddConfigPath("../../../../")
	base.LoadWithOptions(map[string]interface{}{})
	cfg = base
	return base
}

func Get() *AppConfig {
	return cfg
}

func AppName() string {
	return cfg.GetOptionalValue("APP_NAME", "catalog-service")
}
