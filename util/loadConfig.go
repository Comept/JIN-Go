package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	db       DB     `mapstructure"db"`
	servAddr string `mapstructure"servAddr"`
}

type DB struct {
	Addr     string `mapstructure"address"`
	User     string `mapstructure"user"`
	Password string `mapstructure"password"`
	Database string `mapstructure"dataBase"`
}

func LoadConfig() (Config, error) {
	conf := viper.New()
	conf.SetConfigName("conf")
	conf.SetConfigType("env")
	conf.AddConfigPath("./util")
	err := conf.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = conf.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	fmt.Println(config)
	return config, nil
}
