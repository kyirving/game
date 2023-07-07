package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type WorkWxConf struct {
	Host       string
	WebhookKey string
	Corpid     int
	Corpsecret string
	Touser     string
}

type GameConf struct {
	ServerList   string
	ServerStatus string
	Host         string
	PoolNum      int
	GameId       int
	PtId         int
}

type Configs struct {
	*WorkWxConf
	*GameConf
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
	var gameConf GameConf

	err = viper.UnmarshalKey("workwx", &workWxConf)
	if err != nil {
		panic(fmt.Errorf("fatal error viper.UnmarshalKey['workwx']: %w", err))
	}

	err = viper.UnmarshalKey("game", &gameConf)
	if err != nil {
		panic(fmt.Errorf("fatal error viper.UnmarshalKey['game']: %w", err))
	}

	Config = Configs{
		&workWxConf,
		&gameConf,
	}
}
