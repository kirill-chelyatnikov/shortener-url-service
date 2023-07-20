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
		SecretKey     string `yaml:"secretKey"`
	} `yaml:"app"`
	DB struct {
		CDN string `yaml:"cdn"`
	} `yaml:"db"`
}

type Environment struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

type Flags struct {
	ServerAddress string
	BaseURL       string
	FileStorage   string
	DatabaseDSN   string
}

// GetConfig - функция получения конфига приложения
func GetConfig(log *logrus.Logger, path string, fl *Flags) *Config {
	// объявление структур конфига, переменных окружения, флагов
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

	if environment.DatabaseDSN != "" {
		cfg.DB.CDN = environment.DatabaseDSN
	} else {
		cfg.DB.CDN = fl.DatabaseDSN
	}

	log.Info("config received successfully")

	return &cfg
}

func GetFlags() *Flags {
	var fl Flags

	flag.StringVar(&fl.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&fl.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&fl.FileStorage, "f", "", "file storage")
	flag.StringVar(&fl.DatabaseDSN, "d", "", "DatabaseDSN")

	return &fl
}
