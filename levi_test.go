package levi

import (
	"github.com/go-pg/pg"
)

type User struct {
	Table

	Login     string
	Password  []byte
	Email     string
	FirstName string
	LastName  string
}

func (*User) Up(mi *Migration) error {
	v := mi.Up()

	switch {
	case v(1):
		// migrations for v1
		fallthrough
	case v(2):
		// migrations for v2
		fallthrough
	case v(3):
		// migrations for v3
		fallthrough
	case v(4):
		// migrations for v4
	}

	return nil
}

func (*User) Down(mi *Migration) error {
	v := mi.Down()

	switch {
	case v(4):
		// de-migration for v4
		fallthrough
	case v(3):
		// de-migration for v3
		fallthrough
	case v(2):
		// de-migration for v2
		fallthrough
	case v(1):
		// de-migration for v1
	}

	return nil
}

func ExampleMigration() {
	const oldVersion = 2

	db.RunInTransaction(func(tx *pg.Tx) error {
		mi := &Migration{Tx: tx, From: oldVersion}

		var model User
		if err := model.Up(mi); err != nil {
			return err
		}

		return nil
	})
}
