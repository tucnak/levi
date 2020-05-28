package levi

import (
	"runtime/debug"
	"strings"

	"github.com/labstack/echo"
)

// Echo provides access to application's main HTTP router.
func Echo() *echo.Echo {
	return router
}

func newRouter() *echo.Echo {
	e := echo.New()
	e.Renderer = renderer
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			lv := &Lv{Context: c}

			defer func(lv *Lv) {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					stack = stack[strings.Index(stack, "panic("):]
					lv.Panicf("%s\n%s", r, stack)
					err = lv.NoContent(403)
					lv.End(nil)
				}
			}(lv)

			lv.Begin()
			err = next(lv)
			lv.End(err)

			return
		}
	})

	return e
}
