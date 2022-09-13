package config

import (
	"github.com/Ollub/user_service/pkg/db"
	"github.com/kelseyhightower/envconfig"
)

var Cfg Config

type Config struct {
	ServerPort int `envconfig:"SERVER_PORT" default:"8080"`

	Debug        bool   `envconfig:"DEBUG" default:"true"`
	JwtKey       []byte `envconfig:"JWT_KEY" default:"super secret"`
	TokenTTLDays int    `envconfig:"TokenTTL" default:"90"`
	// Postgres config
	DbConf *db.PgCfg
}

func Load() {
	envconfig.MustProcess("", &Cfg)
}
