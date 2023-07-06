package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

const (
	PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
 MIIBOQIBAAJAZSOF9iuJYCxwbUSwFoNteC+Z0rifXvvhJK5NghtimuJmD5xfwySL
 CwXhraKfXEUtz+T6XXA2Rp1tY+pVq+FHwQIDAQABAkAPwNi830smj8VzP5+t4grL
 DZ8IE3m/cbw/2mZ4PYu+VBJ1YjzKecM/HSqq4mvQH8KgGQ02x/3f2MgJ+5Eadw+B
 AiEAvuQGRduC8BeuIfiIHvcOw9rUJaN2DyHmHKYC+4q8zO8CIQCHoqQXLa+0zzg0
 q/KcBJ8SfMxIsWuXTEX2yqdGURaWTwIgaZj+l1ptPp/65jP0KR0Gf/Xn8cJRJuHb
 x/FWKQyAkOUCIQCDy9Rq+Wfc5+aTt+mdFRiFXGMc19nWQLVTZAQ63Zx3HQIgCQzr
 TRtfs3Ax22LNLz2blixObO7FwoV1oC/5ovr6FJE=
 -----END RSA PRIVATE KEY-----`

	PublicKey = `-----BEGIN PUBLIC KEY-----
MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAZSOF9iuJYCxwbUSwFoNteC+Z0rifXvvh
JK5NghtimuJmD5xfwySLCwXhraKfXEUtz+T6XXA2Rp1tY+pVq+FHwQIDAQAB
-----END PUBLIC KEY-----`
)

type Config struct {
	Connection struct {
		Host           string `yaml:"host" env-default:"localhost"`
		HTTPPort       string `yaml:"port" env-default:"80"`
		MasterPassword string `yaml:"master_password" env-default:"135274"`
	} `yaml:"connection"`
	FirstDatabaseConfig struct {
		Host             string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
		Port             string `yaml:"port" env:"DB_PORT" env-default:"5432"`
		User             string `yaml:"user" env:"DB_USER" env-default:"postgres"`
		Password         string `yaml:"password" env:"DB_PASSWORD" env-default:"135274"`
		DbName           string `yaml:"db_name" env:"DB_NAME" env-default:"tasks_db"`
		DBConnectionType string `yaml:"db_connection_type" env-default:"postgres"`
	} `yaml:"first_database_config"`
}

var cfg *Config

func GetConfigs() *Config {
	once := sync.Once{}
	once.Do(func() {
		err := cleanenv.ReadConfig("config.yml", &cfg)
		if err != nil {
			log.Fatalln("Configs wasn't setup")
		}
	})
	return cfg
}
