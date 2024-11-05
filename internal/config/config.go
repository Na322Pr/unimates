package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TG `yaml:"telegram"`
	PG `yaml:"postgres"`
}

type TG struct {
	Token string `yaml:"token" env:"TG_BOT_TOKEN"`
}

type PG struct {
	DB       string `yaml:"db" env:"PG_DB"`
	Host     string `yaml:"host" env:"PG_HOST"`
	Port     int    `yaml:"port" env:"PG_PORT"`
	User     string `yaml:"user" env:"PG_USER"`
	Password string `yaml:"password" env:"PG_PASSWORD"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		log.Fatal("config path is empty")
	}

	return MustLoadPath(path)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config path does not exist: " + err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Print("failed to read config: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("failed to read environment variables" + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
