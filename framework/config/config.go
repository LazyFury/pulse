package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func NewConfig() *viper.Viper {
	return NewConfigFromFile()
}

func NewConfigFromBytes(_bytes []byte) *viper.Viper {
	if len(_bytes) == 0 {
		panic("config file bytes is empty")
	}

	config := viper.New()
	config.SetConfigType("yaml")

	if err := config.ReadConfig(bytes.NewBuffer(_bytes)); err != nil {
		panic("failed to read config file, err: " + err.Error())
	}
	setup(config)
	return config
}

func NewConfigFromFile() *viper.Viper {

	config := viper.New()
	config.SetConfigType("yaml")
	config.AddConfigPath("./")              // current dir
	config.AddConfigPath("./config/")       // config dir
	config.AddConfigPath("../../")          //relative root
	config.AddConfigPath("$HOME/.yoomall/") // home
	config.SetConfigName("config")

	config.SetEnvPrefix("yoomall")
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		fmt.Println("Failed to read config file, err: ", err)
		// gen default config.yaml
		os.Exit(1)
	}
	setup(config)
	return config
}

var Config *viper.Viper

func setup(config *viper.Viper) {
	log.Info(fmt.Sprintf("config: %+v", config.AllSettings()))
	Config = config
}
