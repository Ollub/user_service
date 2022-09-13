package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PgCfg struct {
	User     string `envconfig:"PG_USER" default:"user_service"`
	DbName   string `envconfig:"PG_NAME" default:"user_service"`
	Password string `envconfig:"PG_PASSWORD" default:"user_service"`
	Host     string `envconfig:"PG_HOST" default:"localhost"`
	Port     string `envconfig:"PG_PORT" default:"55437"`
	SSL      bool   `envconfig:"PG_SSL" default:"false"`
	MaxConn  int    `envconfig:"PG_MAX_CONN" default:"10"`
}

func (c *PgCfg) Dsn() string {
	var ssl_mode string
	if c.SSL {
		ssl_mode = "enable"
	} else {
		ssl_mode = "disable"
	}
	return fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		c.User,
		c.DbName,
		c.Password,
		c.Host,
		c.Port,
		ssl_mode,
	)
}

func GetPostgres(cfg *PgCfg) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Dsn())
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}
	err = db.Ping() // first connection here
	if err != nil {
		return nil, fmt.Errorf("connection to postgres failed: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxConn)
	return db, nil
}
