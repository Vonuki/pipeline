package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfg *Config

type Config struct {
	Concurrency int
	URLsCount   int
}

func LoadConfig(filePath string) *Config {
	log.Println("Config loading")
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		processError(err)
	}

	//Read from environment vars
	viper.AutomaticEnv()
	newCfg := Config{}
	err = viper.Unmarshal(&newCfg)
	if err != nil {
		processError(err)
	}
	cfg = &newCfg
	log.Println(fmt.Sprintf("Config loaded: %+v", cfg))
	return cfg
}

func processError(err error) {
	log.Println(fmt.Sprintf("Error with config load: %+v", err))
	os.Exit(2)
}

func GetConfig() *Config {
	if cfg == nil {
		log.Println("Config is not loaded")
		return nil
	}
	return cfg
}
