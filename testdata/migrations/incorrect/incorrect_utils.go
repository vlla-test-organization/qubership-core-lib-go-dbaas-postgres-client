package correct

import (
	"embed"
	"github.com/uptrace/bun/migrate"
)

var incorrectMigrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS // in order to find all *.sql files

func GetIncorrectMigrations() (*migrate.Migrations, error) {
	if err := incorrectMigrations.Discover(sqlMigrations); err != nil {
		return nil, err
	}
	return incorrectMigrations, nil
}
