package news

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Config struct {
	DbHost      string
	DbPort      string
	DbName      string
	UserName    string
	Password    string
	SSLMode     string
	MaxIdleConn int
	MaxOpenConn int
	Debug       bool
}

func (c *Config) getDsn() string {
	return fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=%s",
		c.DbName,
		c.DbHost,
		c.DbPort,
		c.UserName,
		c.Password,
		c.SSLMode,
	)
}

func NewDB(c *Config) (*bun.DB, error) {
	pgConfig, err := pgx.ParseConfig(c.getDsn())
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDB(*pgConfig)
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)

	db := bun.NewDB(sqlDB, pgdialect.New())
	if c.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return db, nil
}
