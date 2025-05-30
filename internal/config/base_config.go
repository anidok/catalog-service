package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config interface {
	GetValue(string) string
	GetIntValue(string) int
}

type configuration map[string]interface{}

var config configuration

type BaseConfig struct {
}

func (b BaseConfig) Load() {
	b.LoadWithOptions(map[string]interface{}{})
}

func (self BaseConfig) LoadWithOptions(options map[string]interface{}) {
	viper.AutomaticEnv()
	viper.SetConfigName("application")
	if options["configPath"] != nil {
		viper.AddConfigPath(options["configPath"].(string))
	} else {
		viper.AddConfigPath("./")
		viper.AddConfigPath("../")
	}
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	config = configuration{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("error unmarshalling config: %w", err))
	}
}

func (b BaseConfig) GetValue(key string) string {
	if _, ok := config[key]; !ok {
		config[key] = getStringOrPanic(key)
	}
	return config[key].(string)
}

func (b BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	if _, ok := config[key]; !ok {
		var value string
		if value = viper.GetString(key); !viper.IsSet(key) {
			value = defaultValue
		}
		return value
	}
	return config[key].(string)
}

func (b BaseConfig) GetIntValue(key string) int {
	if _, ok := config[key]; !ok {
		config[key] = getIntOrPanic(key)
	}
	return config[key].(int)
}

func (b BaseConfig) GetOptionalIntValue(key string, defaultValue int) int {
	if _, ok := config[key]; !ok {
		var value int
		if value = viper.GetInt(key); !viper.IsSet(key) {
			value = defaultValue
		}
		return value
	}
	return config[key].(int)
}

func checkKey(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Errorf("%s key is not set", key))
	}
}

func getStringOrPanic(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func getIntOrPanic(key string) int {
	checkKey(key)
	v, err := strconv.Atoi(viper.GetString(key))
	panicIfErrorForKey(err, key)
	return v
}

func panicIfErrorForKey(err error, key string) {
	if err != nil {
		panic(fmt.Errorf("could not parse key: %s. error: %v", key, err))
	}
}
