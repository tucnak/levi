package levi

import (
	"errors"
	"fmt"
)

var (
	ø = func(err string) error { return errors.New("levi: " + err) }

	ErrNilModel      = ø("nil model interface")
	ErrMigrationFail = ø("migration failed")
	ErrBadArchetype  = ø("archetype not supported")
	ErrBadPaperwork  = ø("paperwork bind fail")
	ErrTmplRepeated  = ø("template loaded repeatedly")
)

// ValidationError should commonly be used in forms.
type ValidationError struct {
	OK      bool   `json:"ok"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("levi: validation error: %s (field %s)", err.Message, err.Field)
}

// MigrationError occurs whenever the table migration fails.
type MigrationError struct {
	Model Model
	Err   error
	Query string
}

func (err *MigrationError) Error() string {
	main := fmt.Sprintf("levi: %T failed to migrate: %v", err.Model, err.Err)
	if err.Query == "" {
		return main
	}
	return main + "\n>>> " + err.Query
}

func (err *MigrationError) Unwrap() error {
	return err.Err
}
