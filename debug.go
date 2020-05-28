package levi

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/go-pg/pg/orm"
)

type route struct {
	method, path, name string
}

func debugRoutes() {
	var routes []route

	for _, r := range router.Routes() {
		if r.Name == "github.com/labstack/echo.glob..func1" {
			continue
		}

		routes = append(routes,
			route{r.Method, r.Path, r.Name})
	}

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].path < routes[j].path
	})

	fmt.Printf("ROUTES\n\n")

	for _, r := range routes {
		fmt.Printf("%8s %s --> %s\n", r.method, r.path, r.name)
	}

	fmt.Println()
}

func debugModels() {
	fmt.Printf("MODELS\n\n")

	var tt []*orm.Table

	for _, model := range tables {
		tt = append(tt, orm.GetTable(reflect.TypeOf(model).Elem()))
	}

	sort.Slice(tt, func(i, j int) bool {
		return tt[i].Name < tt[j].Name
	})

	for _, t := range tt {
		info := ""
		if t.SoftDeleteField == nil {
			info += "#"
		}

		if _, ok := t.FieldsMap["created_at"]; !ok {
			info += "$"
		}

		fmt.Printf("%8s %s <%s%s>\n", "TABLE", t.FullName, t.TypeName, info)
	}

	fmt.Println()
}
