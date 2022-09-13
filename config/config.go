package config

import (
	"github.com/kelseyhightower/envconfig"
	"user_service/pkg/db"
)

var Cfg Config

// Defaults
const (
	DEBUG          = true
	JWT_KEY        = "super secret key"
	TOKEN_TTL_KEYS = 90

	PG_USER_NAME     = "user_service"
	PG_USER_PASSWORD = "user_service"
	PG_HOST          = "localhost"
	PG_PORT          = 55437
	PG_DB_NAME       = "user_service"
)

type Config struct {
	ServerPort int `envconfig:"SERVER_PORT" default:"8080"`

	Debug        bool   `envconfig:"DEBUG" default:"true"`
	JwtKey       []byte `envconfig:"JWT_KEY" default:"super secret"`
	TokenTTLDays int    `envconfig:"TokenTTL" default:"90"`
	// Postgres config
	DbConf *db.PgCfg
	// Redis config
	//RedisHost string `envconfig:"REDIS_HOST" default:"localhost"`
	//RedisPort string `envconfig:"REDIS_PORT" default:"63790"`
	//RedisDb   int    `envconfig:"REDIS_DB" default:"0"`
	//RedisPwd  string `envconfig:"REDIS_PWD" default:""`
}

func Load() {
	envconfig.MustProcess("", &Cfg)
}
