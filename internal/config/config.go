package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

const (
	configPath = "config"
)

type Config struct {
	Version string
	Logger  struct {
		Level string
	}
	Server struct {
		Port         string
		ReadTimeout  int
		WriteTimeout int
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Cache struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
	JwtSecret string
}

var (
	once     sync.Once
	instance *Config
)

func New() *Config {
	once.Do(func() {
		viper.SetConfigType("yaml")
		instance = &Config{}
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&instance); err != nil {
			panic(err)
		}

		log.Println("config initialized")
	})
	return instance
}
