package levi

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

const figurine = `
  _       _____  __     __  ___      _      _____   _   _      _      _   _
 | |     | ____| \ \   / / |_ _|    / \    |_   _| | | | |    / \    | \ | |
 | |     |  _|    \ \ / /   | |    / _ \     | |   | |_| |   / _ \   |  \| |
 | |___  | |___    \ V /    | |   / ___ \    | |   |  _  |  / ___ \  | |\  |
 |_____| |_____|    \_/    |___| /_/   \_\   |_|   |_| |_| /_/   \_\ |_| \_|
`

type Config struct {
	Port        string `os:"PORT"`         // no default
	Production  string `os:"PRODUCTION"`   // no default
	Domain      string `os:"DOMAIN"`       // ex: veritas.icu
	DatabaseURL string `os:"DATABASE_URL"` // ex: postgresql://user@localhost/db

	// Inb4 is guaranteed to execute before the request.
	Inb4 func(*Lv)

	// The interval, in microseconds, for which the same-level adjacent
	// request logs will be automatically grouped.
	//
	// This has an advantage of showing intricate deltas for granular
	// sequential operations within a request.
	//
	// Default: 100 μs.
	LogGroupWindow int `os:"LOG_GROUP_WINDOW"`

	Logger   Logger
	Renderer Renderer
}

var (
	// listen port
	port string
	// postgres
	db *pg.DB
	// web server
	router *echo.Echo
	// serverDomain (required) is used for cookies.
	serverDomain string
	// inb4 is executed for every request before the handler kicks in.
	inb4 func(*Lv)
	// prod indicates whether the app is in prod
	prod bool
	// renderer manages endpoint templates
	renderer Renderer
	// the logger backlog
	logger Logger
	// models
	tables, cues, graphs []Model
)

// Here it all begins.
func Wake(cfg *Config) {
	if err := parseConfig(cfg); err != nil {
		panic(err)
	}

	// 0. Routine checks
	if serverDomain == "" {
		panic("levi: serverDomain global variable must be set before wake")
	}

	// 1. Say hello.
	fmt.Print(figurine)

	// 2. Do migrations.
	if err := migrateUp(); err != nil {
		panic(err)
	}

	debugModels()

	// 3. Set up routing.
	http.Handle("/", router)
	debugRoutes()

	// 4. Listen.
	fmt.Println("WOKE", ":"+port)
	fmt.Println()
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func parseConfig(cfg *Config) error {
	if cfg.Port != "" {
		port = cfg.Port
	}

	if cfg.Production != "" {
		prod = true
	}

	if cfg.Domain != "" {
		serverDomain = cfg.Domain
	}

	if cfg.Inb4 != nil {
		inb4 = cfg.Inb4
	}

	if cfg.Logger != nil {
		logger = cfg.Logger
	} else {
		logger = &StdLogger{}
	}

	if cfg.Renderer != nil {
		renderer = cfg.Renderer
	} else {
		renderer = &HtmlRenderer{}
	}

	if cfg.LogGroupWindow == 0 {

	}

	return nil
}

// IsProd is true when the app is running in production.
func IsProd() bool {
	return prod
}

// IsDev is true when the app is run in dev environment.
func IsDev() bool {
	return !prod
}
