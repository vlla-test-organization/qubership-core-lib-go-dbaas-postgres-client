package correct

import (
	"embed"
	"github.com/uptrace/bun/migrate"
)

var correctMigrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS // in order to find all *.sql files

func GetCorrectMigrations() (*migrate.Migrations, error) {
	if err := correctMigrations.Discover(sqlMigrations); err != nil {
		return nil, err
	}
	return correctMigrations, nil
}
