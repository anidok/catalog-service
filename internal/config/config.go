package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	BaseConfig
	openSearchConfig *OpenSearchConfig
}

var cfg *AppConfig

func Load() *AppConfig {
	base := &AppConfig{}
	viper.GetViper().AddConfigPath("../../")
	viper.GetViper().AddConfigPath("../../../")
	viper.GetViper().AddConfigPath("../../../../")
	base.LoadWithOptions(map[string]interface{}{})
	cfg = base
	cfg.openSearchConfig = NewOpenSearchConfig(cfg)
	return base
}

func Get() *AppConfig {
	return cfg
}

func AppName() string {
	return cfg.GetOptionalValue("APP_NAME", "catalog-service")
}

func LogLevel() string {
	return cfg.GetOptionalValue("LOG_LEVEL", "debug")
}

func LogFormat() string {
	return cfg.GetOptionalValue("LOG_FORMAT", "json")
}

func OpenSearch() *OpenSearchConfig {
	return cfg.openSearchConfig
}
