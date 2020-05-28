package levi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Migrant is a model that is capable of manual migration.
type Migrant interface {
	Version() int
	Up(*Migration) error
}

// Migration is a special transaction type that gets called by
// the migrant model whenever it needs to migrate.
type Migration struct {
	*pg.Tx
	From, To int
}

// Up is called in the Migration function to simplify the outline
// of the up-migration list.
//
//		v := mi.Up()
//		switch {
//		case v(1):
//			mi.Exec("...")
//			fallthrough
//		case v(2):
//			// ...
//		}
//
func (mi *Migration) Up() func(int) bool {
	return func(v int) bool {
		return v > mi.From
	}
}

// Down is called in the Migration function to simplify the outline
// of the down-migration list.
func (mi *Migration) Down() func(int) bool {
	return func(v int) bool {
		return v > mi.To
	}
}

type tableVersion struct {
	tableName struct{} `pg:"migrations"`

	Table   string
	Version int
}

func migrateUp() error {
	var versions []tableVersion
	if err := db.Select(&versions); err != nil {
		return &MigrationError{nil, err, "SELECT * FROM migrations"}
	}
	version := map[string]int{}
	for _, t := range versions {
		version[t.Table] = t.Version
	}

	return db.RunInTransaction(func(tx *pg.Tx) error {
		for _, model := range tables {
			migrant := model.(Migrant)
			tableName := tableOf(model).Name

			mi := &Migration{Tx: tx, From: version[tableName]}
			if err := migrant.Up(mi); err != nil {
				return &MigrationError{model, err, ""}
			}
		}

		return nil
	})
}

func autoMigrate(model Model) error {
	if model.Type() != TABLE {
		return ErrBadArchetype
	}

	return db.RunInTransaction(func(tx *pg.Tx) error {
		err := tx.CreateTable(model, &orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return &MigrationError{model, err, "CREATE TABLE IF NOT EXISTS"}
		}

		table := tableOf(model)
		columns := make([]string, 0, len(table.Fields))
		for _, field := range table.Fields {
			columns = append(columns, `"`+field.SQLName+`" `+field.SQLType)
		}

		for i, column := range columns {
			query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s",
				table.Name, column)

			defaultExpr := table.Fields[i].Default
			if defaultExpr != "" {
				query += " DEFAULT (" + string(defaultExpr) + ")"
			}

			if pgTag, ok := table.Fields[i].Field.Tag.Lookup("pg"); ok {
				if strings.Contains(pgTag, "notnull") {
					query += " NOT NULL"
				}
			}

			_, err := tx.Exec(query)
			if err != nil {
				return &MigrationError{model, err, query}
			}
		}

		return nil
	})
}

func tableOf(model Model) *orm.Table {
	return orm.GetTable(reflect.TypeOf(model).Elem())
}
