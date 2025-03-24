package model

import (
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
)

// PgOptions is a struct which should be used for datasource configuration
type PgOptions struct {
	// Options for bun.DB object creation
	BunOptions []bun.DBOption
	// Options for sql.DB object creation
	SqlOptions []stdlib.OptionOpenDB
}
