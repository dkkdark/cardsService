package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

const (
	PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
-----END RSA PRIVATE KEY-----`

	PublicKey = `-----BEGIN PUBLIC KEY-----
-----END PUBLIC KEY-----`
)

type Config struct {
	Connection struct {
		Host           string `yaml:"host" env-default:"localhost"`
		HTTPPort       string `yaml:"port" env-default:"80"`
		MasterPassword string `yaml:"master_password" env-default:"123456"`
	} `yaml:"connection"`
	FirstDatabaseConfig struct {
		Host             string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
		Port             string `yaml:"port" env:"DB_PORT" env-default:"5432"`
		User             string `yaml:"user" env:"DB_USER" env-default:"postgres"`
		Password         string `yaml:"password" env:"DB_PASSWORD" env-default:"123456"`
		DbName           string `yaml:"db_name" env:"DB_NAME" env-default:"tasks_db"`
		DBConnectionType string `yaml:"db_connection_type" env-default:"postgres"`
	} `yaml:"first_database_config"`
	SecondDatabaseConfig struct {
		Host string `yaml:"host" env:"MONGO_HOST" env-default:"localhost"`
		Port string `yaml:"port" env:"MONGO_PORT" env-default:"27017"`
	} `yaml:"second_database_config"`
}

var cfg *Config

func GetConfigs() *Config {
	once := sync.Once{}
	once.Do(func() {
		cfg = &Config{}
		err := cleanenv.ReadConfig("./config.yml", cfg)
		if err != nil {
			log.Fatalln("Configs wasn't setup")
		}
	})
	return cfg
}
