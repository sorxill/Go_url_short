package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Address     string
	Storage     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func MustLoad() *Config {
	// По соглашению приставка Must - функция будет вызывать панику, вместо возврата ошибки. Такая логика не лучший способ.

	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("reading the config is not exist")
	}

	cfg.Env = os.Getenv("ENV")
	cfg.Storage = os.Getenv("STORAGE")
	cfg.Address = os.Getenv("ADDRESS")

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatal("bad timeout")
	}
	cfg.Timeout = timeout

	idle_timeout, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal("bad idle_timeout")
	}
	cfg.IdleTimeout = idle_timeout

	return &cfg
}
