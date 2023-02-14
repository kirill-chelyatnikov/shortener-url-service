package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	} `yaml:"server"`
	App struct {
		ShortedURLLen uint8 `yaml:"shortedURLLen"`
	} `yaml:"app"`
}

// GetConfig - функция получения конфига приложения
func GetConfig(log *logrus.Logger, path string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatalf("can't get config! %s", err)
	}

	log.Info("config received successfully")

	return &cfg
}
