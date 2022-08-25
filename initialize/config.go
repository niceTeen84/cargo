package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	FILE_NAME     = "config"
	FILE_TYPE     = "yaml"
	FILE_LOCATION = "."
)

var Conf *viper.Viper

func init() {
	v := viper.New()
	v.SetConfigName(FILE_NAME)
	v.SetConfigType(FILE_TYPE)
	v.AddConfigPath(FILE_LOCATION)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file not found; ignore error if desired")
		} else {
			panic("Config file was found but another error was produced")
		}
	}
	v.OnConfigChange(func(event fsnotify.Event) {
		fmt.Print("config file changes")
	})
	v.WatchConfig()
	Conf = v
}
