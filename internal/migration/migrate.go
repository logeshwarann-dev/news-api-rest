package migration

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

func New() *migrate.Migrations {
	return migrations
}

//go:embed *.sql
var sqlMigration embed.FS

func init() {
	migrations.DiscoverCaller()
	if err := migrations.Discover(sqlMigration); err != nil {
		panic(err)
	}
}
