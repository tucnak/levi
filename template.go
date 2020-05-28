package levi

import (
	"fmt"
	"html/template"
	"io"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// Renderer is the endpoint template engine.
type Renderer interface {
	Render(w io.Writer, name string, data interface{}, c echo.Context) error
}

// HtmlRenderer is a regular HTML template
type HtmlRenderer struct {
	// Default: ["public/*.html"]
	Glob []string

	// Default:
	//  {
	//		"htime":       humanize.Time,
	//		"random_uuid": uuid.New,
	//  }
	Funcs map[string]interface{}

	templates *template.Template
	once      sync.Once
}

func (t *HtmlRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var err error
	t.once.Do(func() {
		err = t.load()
	})

	if err != nil {
		return err
	}

	if IsDev() {
		if err := t.load(); err != nil {
			return err
		}
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *HtmlRenderer) load() error {
	defaultFuncs := map[string]interface{}{
		"htime":       humanize.Time,
		"random_uuid": uuid.New,
	}

	t.templates = template.New("")
	t.templates.Funcs(defaultFuncs).Funcs(t.Funcs)

	if t.Glob == nil {
		_, _ = t.templates.ParseGlob("public/*.html")
	}

	for _, glob := range t.Glob {
		_, err := t.templates.ParseGlob(glob)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTmplRepeated, err)
		}
	}

	return nil
}
