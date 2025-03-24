# How to migrate from pg/go-pg to uptrace/bun

This article contains some information about migration process from versions below 1.0.0-beta.1,
ehich were based on pg/go-pg library, to versions equal or upper than 1.0.0-beta.1, which are now
based on uptrace/bun library.

* [Guide from uptrace/bun authors](#guide-from-uptracebun-authors)
* [Migrations](#migrations)
* [Schemas](#schemas)
* [Queries](#queries)

## Guide from uptrace/bun authors

Authors of uptrace/bun library provide some basic instructions for convenient migration from go-pg. You may
find these instructions here: https://bun.uptrace.dev/guide/pg-migration.html .

## Migrations 

uptrace/bun has it's own implementation of migration procedure. In go-pg table with information about
migrations was named as `gopg_migrations`. In uptrace/bun table has name `bun_migrations`. If you already 
has some executed migrations you should manually mark them as 'applied'. 

Also note that now files with migrations should have new name format: 
* timestamp_name.up.sql for sql migrations
* timestamp_name.go for go migrations

We recommend using order number instead of timestamp for convenience, like "00000000000001_init.go"

More information about migrations may be found here: https://bun.uptrace.dev/guide/migrations.html 

```go
    var migrations = migrate.NewMigrations()
    var lastVersion int64 // select info about last applied migration
	err = db.QueryRow("select version from " + schemaVersionTable + " where version <> 0 order by version::int desc limit 1;").
		Scan(&lastVersion)
	if err != nil {
		log.Panic("can not get last version from %s: %v", schemaVersionTable, err.Error())
	}
    sorted := migrations.Sorted() // al discovered migrations
	
	// get already applied migrations in order not to mark
	// one migration as applied several times
    applied, _ := migrator.AppliedMigrations(ctx)
    for i := 0; i < len(sorted); i++ {
        migrationNum, err := strconv.ParseInt(sorted[i].Name, 10, 0)
        if err != nil {
            log.Errorf("Wrong migration name %+v", err.Error())
            return err
        }
        if migrationNum > version { // migration is new, it should be executed
            return nil
        }
        if !contains(applied, sorted[i]) { // if current migration was not yet marked as applied
            err := migrator.MarkApplied(ctx, &sorted[i])
            if err != nil {
                log.Errorf("Error with marking existing migrations as applied %+v", err.Error())
                return err
            }
        }
    }
```

## Schemas

Let's go with example here. 

```go
type Entity struct {
    tableName struct{} `pg:"entities"`
	
	Id          int     `pg:,pk`
	Name        string  `pg:name`
	Age         string  `pg:age`
	IsEmployee  bool    `pg:is_employee,default:false`
}
```

To make this entity work with bun user should:
1. Replace tableName struct{} `pg:"entities"` with bun.BaseModel `bun:"entities"`
2. replace all `pg` keys to `bun`
3. Unlike go-pg, Bun does not marshal Go zero values as SQL NULLs by default. To get the old behavior, use `nullzero` tag option
4. Information about ORM relations may be found here: https://bun.uptrace.dev/guide/relations.html 
5. Although bun declares support for `default` tag, sometimes it doesn't work as expected with update queries (check https://github.com/uptrace/bun/issues/664, https://github.com/uptrace/bun/discussions/659, https://github.com/uptrace/bun/issues/523)

More info about defining model may be found here: https://bun.uptrace.dev/guide/models.html . 
After all transformations user should have something like:

```go
type Entity struct {
    bun.BaseModel `bun:"entities"`
	
	Id          int     `bun:,pk`
	Name        string  `bun:name`
	Age         string  `bun:age`
	IsEmployee  bool    `bun:is_employee`
}
```

## Queries 

`*pg.Query `is split into smaller structs, for example, `bun.SelectQuery`, `bun.InsertQuery`,
`bun.UpdateQuery`, `bun.DeleteQuery` and so on. 

```go
// go-pg API
err := db.Model(&users).Select()
res, err := db.Model(&users).Insert()
res, err := db.Model(&user).WherePK().Update()
res, err := db.Model(&users).WherePK().Delete()
```

Should be transformed to:
```go
// uptrace/bun API
err := db.NewSelect().Model(&users).Scan(ctx) // use Scan(ctx) terminate function
res, err := db.NewInsert().Model(&users).Exec(ctx) // use Exec(ctx) to execute request
res, err := db.NewUpdate().Model(&users).WherePK().Exec(ctx) // use Exec(ctx) to execute request
res, err := db.NewDelete().Model(&users).WherePK().Exec(ctx) // use Exec(ctx) to execute request
```