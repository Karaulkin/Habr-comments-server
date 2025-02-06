package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Storage    DB     `yaml:"db" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string        `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Database string `yaml:"database" env-default:"postgres"`
}

func MustLoad() *Config {
	if err := loadEnv(); err != nil {
		log.Printf("error loading environment variables: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH") //загрузка из переменной окружения
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func loadEnv() error {
	var err error

	if err = godotenv.Load("local.env"); err == nil {
		return nil
	}

	if err = godotenv.Load("dev.env"); err == nil {
		return nil
	}

	if err = godotenv.Load("prod.env"); err == nil {
		return nil
	}

	return err
}
