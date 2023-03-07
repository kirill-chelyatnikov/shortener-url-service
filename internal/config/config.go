package config

import (
	"github.com/caarlos0/env"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"server"`
	App struct {
		ShortedURLLen uint8  `yaml:"shortedURLLen"`
		BaseURL       string `yaml:"baseURL"`
	} `yaml:"app"`
}

type Environment struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

// GetConfig - функция получения конфига приложения
func GetConfig(log *logrus.Logger, path string) *Config {
	//структуры конфига и переменных окружения
	var cfg Config
	var environment Environment

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatalf("can't get config! %s", err)
	}

	err = env.Parse(&environment)
	if err != nil {
		log.Fatalf("can't get environments! %s", err)
	}

	if environment.ServerAddress != "" {
		cfg.Server.Address = environment.ServerAddress
	}

	if environment.BaseURL != "" {
		cfg.App.BaseURL = environment.BaseURL
	}

	log.Info("config received successfully")

	return &cfg
}
