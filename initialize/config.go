package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
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
	keys := v.AllKeys()
	fmt.Println(keys)
	return v
}
