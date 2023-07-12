package config

import (
	"fmt"
	"os"

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

var (
	configFilePath string
	configFileName string
)

//初始化配置
func init() {
	env := os.Getenv("APP_ENV")
	//生产环境
	if env == "prod" {
		//服务器配置目录
		configFilePath = "/data/go/game/config"
		configFileName = "config.prod"
	} else {
		configFilePath = "./config"
		configFileName = "config.dev"
	}
	fmt.Println("env:", env)
	fmt.Printf("configName:%s\n configPath:%s\n", configFileName, configFilePath)

	viper.SetConfigName(configFileName) // name of config file (without extension)
	viper.SetConfigType("ini")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configFilePath) // path to look for the config file in
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
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
