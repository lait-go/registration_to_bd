package conf

import (
	"os"

	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	Env       string
	Host      string
	LogPath   string
	ErrorPath string
	DateConf  string
}

func Config_work() *Config {
	_ = godotenv.Load(".env") // загружаем .env по пути /env

	conf := &Config{
		Env:       os.Getenv("ENV"),
		Host:      os.Getenv("HOST"),
		LogPath:   os.Getenv("LOG_PATH"),
		DateConf:  os.Getenv("DATE_CONF"),
		ErrorPath: os.Getenv("ERROR_PATH"),
	}

	return conf
}