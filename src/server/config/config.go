package config

import (
	"fmt"
	"github.com/spf13/viper"
	"thor/src/constanta"
)

func ConfigApps(path string) *DefaultConfig {
	//viper.SetConfigFile(path)
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	//viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(constanta.SysConfigFailedRead)
	}

	var conf DefaultConfig

	if err := viper.Unmarshal(&conf); err != nil {
		panic(constanta.SysConfigUnmarshall)
	}
	fmt.Println(conf)

	return &conf
}
