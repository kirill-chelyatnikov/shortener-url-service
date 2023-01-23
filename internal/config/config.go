package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

const configURL = "internal/config/config.yml"

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"server"`
	App struct {
		ShortedURLLen uint8 `yaml:"shortedURLLen"`
	} `yaml:"app"`
}

// GetConfig - функция получения конфига приложения
func GetConfig(log *logrus.Logger) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configURL, &cfg)
	if err != nil {
		log.Fatalf("can't get config! %s", err)
	}
	log.Info("config received successfully")

	return &cfg
}
