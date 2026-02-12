package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/logeshwarann-dev/news-api-rest/internal/migration"
	"github.com/logeshwarann-dev/news-api-rest/internal/postgres"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func main() {
	db, err := postgres.NewDB(&postgres.Config{
		DbHost:   os.Getenv("DATABASE_HOST"),
		DbPort:   os.Getenv("DATABASE_PORT"),
		DbName:   os.Getenv("DATABASE_NAME"),
		UserName: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal(err)
	}

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(),
	))

	app := &cli.App{
		Name: "migrate",
		Commands: []*cli.Command{
			newMigrationCmd(migrate.NewMigrator(db, migration.New(),
				migrate.WithMarkAppliedOnSuccess(true))),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newMigrationCmd(m *migrate.Migrator) *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "migrate is used for performing db schema setup",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "initialize the migration",
				Action: func(ctx *cli.Context) error {
					return m.Init(ctx.Context)
				},
			},
			{
				Name:  "up",
				Usage: "run up migration",
				Action: func(ctx *cli.Context) error {
					if err := m.Lock(ctx.Context); err != nil {
						return err
					}
					defer m.Unlock(ctx.Context)

					group, err := m.Migrate(ctx.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Println("no new migrations to perform")
						return nil
					}
					fmt.Printf("migrated to %s", group)
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "run down migration",
				Action: func(ctx *cli.Context) error {
					if err := m.Lock(ctx.Context); err != nil {
						return err
					}
					defer m.Unlock(ctx.Context)
					group, err := m.Rollback(ctx.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Println("no new migrations to run")
					}
					fmt.Printf("rolled back to %s\n", group)
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create new migration",
				Action: func(ctx *cli.Context) error {
					name := strings.Join(ctx.Args().Slice(), "_")
					files, err := m.CreateTxSQLMigrations(ctx.Context, name)
					if err != nil {
						return err
					}
					for _, f := range files {
						fmt.Printf("created migration %s (%s)\n", f.Name, f.Path)
					}
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "status of migration",
				Action: func(ctx *cli.Context) error {
					mg, err := m.MigrationsWithStatus(ctx.Context)
					if err != nil {
						return err
					}
					fmt.Printf("migration status: %s\n", mg)
					fmt.Printf("unapplied migration: %s\n", mg.Unapplied())
					fmt.Printf("last group migration: %s\n", mg.LastGroup())
					return nil
				},
			},
		},
	}
}
