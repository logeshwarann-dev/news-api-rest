package migration

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

func New() *migrate.Migrations {
	return Migrations
}

//go:embed *.sql
var sqlMigration embed.FS

func init() {
	Migrations.DiscoverCaller()
	if err := Migrations.Discover(sqlMigration); err != nil {
		panic(err)
	}
}
