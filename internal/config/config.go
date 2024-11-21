package config

import (
	"database/sql"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port int
	Addr string
	Env  string
	Db   struct {
		Dsn string
		Db  *sql.DB
	}
}

func Configure(fileName string) (*Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(fileName, &cfg)
	return &cfg, err
}
