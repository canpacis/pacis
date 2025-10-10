package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/server/middleware"
)

type Environment string

const (
	Dev  = Environment("DEVELOPMENT")
	Prod = Environment("PRODUCTION")
)

type AppOptions struct {
	env       Environment
	devserver string
}

type SpeculationRule struct {
	URLs []string `json:"urls,omitempty"`
}

type Speculation struct {
	Prerender []SpeculationRule `json:"prerender,omitempty"`
	Prefetch  []SpeculationRule `json:"prefetch,omitempty"`
}

type App struct {
	assets      http.Handler
	entries     map[string]string
	options     *AppOptions
	speculation Speculation
	middlewares []func(http.Handler) http.Handler
}

type Route interface {
	Path() string
	Handler(*App) http.Handler
}

func (a *App) Register(mux *http.ServeMux, route Route) {
	mux.Handle("GET "+route.Path(), route.Handler(a))
}

func (a *App) Speculations() html.Node {
	return html.Script(html.Type("speculationrules"), html.JSON(a.speculation))
}

func (a *App) Use(middlewares ...func(http.Handler) http.Handler) {
	a.middlewares = append(a.middlewares, middlewares...)
}

func (a *App) SetBuildDir(name string, dir fs.FS) error {
	assets, err := fs.Sub(dir, path.Join(name, "assets"))
	if err != nil {
		return fmt.Errorf("failed to read asset directory: %w", err)
	}
	a.assets = http.FileServerFS(assets)

	file, err := dir.Open(path.Join(name, ".vite/manifest.json"))
	if err != nil {
		return fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer file.Close()

	type Item struct {
		File    string   `json:"file"`
		Name    string   `json:"name"`
		Names   []string `json:"names"`
		Src     string   `json:"src"`
		IsEntry bool     `json:"isEntry"`
	}
	type Manifest map[string]Item

	manifest := new(Manifest)
	if err := json.NewDecoder(file).Decode(manifest); err != nil {
		return fmt.Errorf("failed to decode manifest file: %w", err)
	}

	for name, item := range *manifest {
		if item.IsEntry {
			a.entries[strings.TrimPrefix(name, "web/src/")] = "/" + item.File
		}
	}

	return nil
}

func (a *App) ServeAssets() http.Handler {
	return a.assets
}

func Asset(app *App, name string) string {
	if app.options.env == Dev {
		return app.options.devserver + "/src/web/" + name
	}
	entry, ok := app.entries[name]
	if !ok {
		log.Fatalf("failed to retrieve asset %s", name)
	}
	return entry
}

func WithDevServer(url string) func(*AppOptions) {
	return func(ao *AppOptions) {
		ao.devserver = url
	}
}

func WithEnv(env Environment) func(*AppOptions) {
	return func(ao *AppOptions) {
		ao.env = env
	}
}

func NewApp(options ...func(*AppOptions)) *App {
	opts := &AppOptions{
		env:       Dev,
		devserver: "http://localhost:5173",
	}

	for _, opt := range options {
		opt(opts)
	}

	app := &App{
		assets:  http.NotFoundHandler(),
		entries: make(map[string]string),
		options: opts,
	}
	app.Use(middleware.ColorScheme, middleware.Cache(time.Hour*24*365))
	return app
}
