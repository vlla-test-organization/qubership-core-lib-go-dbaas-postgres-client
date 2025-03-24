package correct

import (
	"context"
	"github.com/uptrace/bun"
)

type Book struct {
	Code int
}

func init() {
	correctMigrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		bookToInsert := Book{Code: 111}
		_, errInsert := db.NewInsert().Model(&bookToInsert).Exec(ctx)
		return errInsert
	}, func(ctx context.Context, db *bun.DB) error {
		return nil
	})
}
