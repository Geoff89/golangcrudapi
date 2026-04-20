package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string `mapstructure:"port"`
	ConnectionString string `mapstructure:"connection_string"`
}

var AppConfig Config

func LoadAppConfig() {
	log.Println("Loading server configurations...")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatal(err)
	}

	if AppConfig.Port == "" {
		AppConfig.Port = "8080"
	}
	if AppConfig.ConnectionString == "" {
		log.Fatal(fmt.Errorf("connection_string must be configured"))
	}
}
