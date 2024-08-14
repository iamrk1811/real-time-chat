package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Site struct {
		Title        string `mapstructure:"TITLE"`
		Port         string `mapstructure:"PORT"`
		ReadTimeout  int    `mapstructure:"HTTP_READ_TIMEOUT"`
		WriteTimeout int    `mapstructure:"HTTP_WRITE_TIMEOUT"`
	} `mapstructure:"SITE"`

	DB struct {
		CONN_URL string `mapstructure:"CONN_URL"`
	} `mapstructure:"DB"`
	ProtectedPaths ProtectedPaths
}

func NewConfig() *Config {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("can't read config file")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("error while reading config")
	}
	config.ProtectedPaths = *NewProtectedPaths()
	return &config
}
