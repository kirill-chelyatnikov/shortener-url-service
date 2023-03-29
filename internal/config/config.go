package config

import (
	"flag"
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
		FileStorage   string `yaml:"fileStorage"`
	} `yaml:"app"`
}

type Environment struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
}

type flags struct {
	ServerAddress string
	BaseURL       string
	FileStorage   string
}

// GetConfig - функция получения конфига приложения
func GetConfig(log *logrus.Logger, path string) *Config {
	// объявление структур конфига, переменных окружения, флагов
	var cfg Config
	var environment Environment
	var fl flags

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatalf("can't get config! %s", err)
	}

	err = env.Parse(&environment)
	if err != nil {
		log.Fatalf("can't get environments! %s", err)
	}

	flag.StringVar(&fl.ServerAddress, "a", "localhost", "server address")
	flag.StringVar(&fl.BaseURL, "b", "localhost", "base url")
	flag.StringVar(&fl.FileStorage, "f", "", "file storage")

	// замена значений конфига при условии передачи переменных окружения или флагов (переменные окружения в приоритете)
	if environment.ServerAddress != "" {
		cfg.Server.Address = environment.ServerAddress
	} else {
		cfg.Server.Address = fl.ServerAddress
	}

	if environment.BaseURL != "" {
		cfg.App.BaseURL = environment.BaseURL
	} else {
		cfg.App.BaseURL = fl.BaseURL
	}

	if environment.FileStorage != "" {
		cfg.App.FileStorage = environment.FileStorage
	} else {
		cfg.App.FileStorage = fl.FileStorage
	}

	log.Info("config received successfully")

	return &cfg
}
