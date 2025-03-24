package correct

import (
	"context"
	"errors"
	"github.com/uptrace/bun"
)

type Book struct {
	Code int
}

func init() {
	incorrectMigrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		return errors.New("some error during migration")
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().Model((*Book)(nil)).Exec(ctx)
		return err
	})
}
