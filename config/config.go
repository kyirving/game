package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type WorkWxConf struct {
	Host       string
	Corpid     int
	Corpsecret string
}

type Configs struct {
	*WorkWxConf
}

var Config Configs

//初始化配置
func init() {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("ini")      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to look for the config file in
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var workWxConf WorkWxConf
	err = viper.UnmarshalKey("workwx", &workWxConf)
	if err != nil {
		panic(fmt.Errorf("fatal error viper.UnmarshalKey: %w", err))
	}

	Config = Configs{
		&workWxConf,
	}
}
