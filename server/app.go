// Package server provides the core application structure and utilities for configuring and running
// a web server with environment-specific behavior, middleware support, and static asset management.
//
// Types:
//
//   - Environment: Represents the application environment ("DEVELOPMENT" or "PRODUCTION").
//   - AppOptions: Configuration options for the application, such as environment, dev server URL, web files path, and assets directory.
//   - App: The main application struct, holding asset handlers, entry points, options, and middleware stack.
//   - Route: Interface for defining HTTP routes with a path and handler.
//
// Functions:
//
//   - (*App) Register: Registers a Route with a given ServeMux using the GET method.
//   - (*App) Use: Adds middleware(s) to the application's middleware stack.
//   - (*App) SetBuildDir: Configures the asset handler and entry points from a build directory and manifest file.
//   - (*App) ServeAssets: Returns the HTTP handler for serving static assets.
//   - WithEnv: Returns an option function to set the application environment.
//   - WithDevServer: Returns an option function to set the development server URL.
//   - WithWebFiles: Returns an option function to set the web files path.
//   - WithAssetsDir: Returns an option function to set the assets directory name.
//   - NewApp: Constructs a new App instance with the provided options and default middleware.
//
// Usage:
//
//	The server package is designed to be flexible and extensible, allowing developers to configure
//	environment-specific settings, register routes, apply middleware, and serve static assets efficiently.
package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/canpacis/pacis/server/middleware"
)

// Environment describes the app environment and lets the other parts of the app choose environment specific behaviour.
// Valid values are "DEVELOPMENT" and "PRODUCTION".
type Environment string

const (
	// Reprensents a development environment of the app.
	Dev = Environment("DEVELOPMENT")
	// Reprensents a production environment of the app.
	Prod = Environment("PRODUCTION")
)

// Configuration options for the application, such as environment, dev server URL, web files path, and assets directory.
type AppOptions struct {
	env       Environment
	devserver string
	webfiles  string
	assetsdir string
}

// The main application struct, holding asset handlers, entry points, options, and middleware stack.
type App struct {
	assets      http.Handler
	entries     map[string]string
	options     *AppOptions
	middlewares []func(http.Handler) http.Handler
}

// Interface for defining HTTP routes with a path and handler.
type Route interface {
	Path() string
	Handler(*App) http.Handler
}

// Registers a Route with a given ServeMux using the GET method.
func (a *App) Register(mux *http.ServeMux, route Route) {
	mux.Handle("GET "+route.Path(), route.Handler(a))
}

// Adds middleware(s) to the application's middleware stack.
func (a *App) Use(middlewares ...func(http.Handler) http.Handler) {
	a.middlewares = append(a.middlewares, middlewares...)
}

// Configures the asset handler and entry points from a build directory and manifest file.
func (a *App) SetBuildDir(name string, dir fs.FS) error {
	assets, err := fs.Sub(dir, path.Join(name, a.options.assetsdir))
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
			a.entries[strings.TrimPrefix(name, a.options.webfiles+"/")] = "/" + item.File
		}
	}

	return nil
}

// Returns the HTTP handler for serving static assets.
func (a *App) ServeAssets() http.Handler {
	return a.assets
}

// Returns an option function to set the application environment.
func WithEnv(env Environment) func(*AppOptions) {
	return func(ao *AppOptions) {
		ao.env = env
	}
}

// Returns an option function to set the development server URL.
func WithDevServer(url string) func(*AppOptions) {
	return func(ao *AppOptions) {
		ao.devserver = url
	}
}

// Returns an option function to set the web files path.
func WithWebFiles(path string) func(*AppOptions) {
	return func(ao *AppOptions) {
		ao.webfiles = path
	}
}

// Returns an option function to set the assets directory name.
func WithAssetsDir(name string) func(*AppOptions) {
	return func(ao *AppOptions) {
		ao.assetsdir = name
	}
}

// Constructs a new App instance with the provided options and default middleware.
func NewApp(options ...func(*AppOptions)) *App {
	opts := &AppOptions{
		env:       Dev,
		devserver: "http://localhost:5173",
		webfiles:  "src/web",
		assetsdir: "assets",
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
