package levi

import (
	"context"
	"time"
)

// Table is a normal soft-deletable model.
type Table struct {
	baseTable
	timeTable
	DeletedAt *time.Time `pg:",soft_delete" json:"-"`
}

// DestructibleTable (marked ^) is a model for destructible tables.
type DestructibleTable struct {
	baseTable
	timeTable
}

// LighweightTable (marked ˚) is a primitive single-id model.
type LightweightTable struct {
	baseTable
}

type baseTable struct {
	Id int64 `json:"-"`
}

func (baseTable) Model() Archetype    { return TABLE }
func (baseTable) Version() int        { return 0 } // automatic migration
func (baseTable) Up(*Migration) error { return nil }

type timeTable struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *timeTable) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}

	return ctx, nil
}

func (m *timeTable) BeforeUpdate(ctx context.Context) (context.Context, error) {
	m.UpdatedAt = time.Now()

	return ctx, nil
}
