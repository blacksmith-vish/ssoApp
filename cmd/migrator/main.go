package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "storage-path")
	flag.StringVar(&migrationsPath, "migrations-path", "", "migrations-path")
	flag.StringVar(&migrationsTable, "migrations-table", "", "migrations-table")

	flag.Parse()

	if validate.Var(storagePath, "required") != nil {
		panic("storagePath required")
	}

	if validate.Var(migrationsPath, "required") != nil {
		panic("migrationsPath required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}

	fmt.Println("migrated successfully")

}
