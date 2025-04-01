package global

import (
	"github.com/lazyfury/pulse/framework/config"
	"github.com/spf13/viper"
)

func GetConfig() *viper.Viper {
	return config.Config
}
