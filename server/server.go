// Package server provides the core HTTP server implementation for the application,
// including environment configuration, middleware management, static asset handling,
// development server proxying, and graceful shutdown capabilities.
//
// The Server type embeds http.ServeMux and extends it with middleware support,
// asset manifest management, and environment-specific behaviors. It supports both
// development and production environments, allowing for flexible asset resolution
// and request handling.
//
// Usage:
//   - Create a new server instance using New() with appropriate Options.
//   - Register middleware and route handlers as needed.
//   - Serve static assets and pages using provided methods.
//   - Start the server with Serve(), which handles graceful shutdown.
//
// Example:
//
//	options := &server.Options{Env: server.Dev, Port: ":8080"}
//	srv := server.New(options)
//	srv.Serve()
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/canpacis/pacis/internal"
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

// Options holds the configuration settings for the server, including environment,
// port, development server URL, logger, and HTTP request multiplexer.
type Options struct {
	Env       Environment
	Port      string
	DevServer *url.URL
	Logger    *slog.Logger
	Mux       *http.ServeMux
}

type entry struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Names   []string `json:"names"`
	Src     string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
}

type manifest map[string]entry

type Server struct {
	*http.ServeMux
	middlewares []middleware.Middleware
	manifest    manifest
	options     *Options
	notfound    http.Handler
}

// Adds middleware(s) to the application's middleware stack.
func (s *Server) Use(middlewares ...middleware.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func clean(pattern string, prefix string) string {
	replacer := strings.NewReplacer("GET", "", "POST", "", "PATCH", "", "PUT", "", "DELETE", "", "OPTIONS", "")
	pattern = replacer.Replace(pattern)
	pattern = strings.TrimSpace(pattern)
	pattern = strings.TrimSpace(pattern)

	if len(pattern) == 0 {
		if len(prefix) > 0 {
			return prefix + " /"
		}
		return "/"
	}
	if len(prefix) > 0 {
		return prefix + " " + pattern
	}
	return pattern
}

func (s *Server) HandlePage(pattern string, page Page, layout Layout, middlewares ...middleware.Middleware) {
	s.Handle(clean(pattern, "GET"), PageHandler(s, page, layout, middlewares...))

	actioner, ok := page.(interface{ Actions() map[string]ActionFunc })
	if !ok {
		return
	}
	actions := actioner.Actions()
	s.Handle(clean(pattern, "POST"), ActionsHandler(s, actions, middlewares...))
}

func (s *Server) SetBuildDir(name string, dir fs.FS, vite fs.FS) error {
	static, err := fs.Sub(dir, name)
	if err != nil {
		return err
	}

	var handler http.Handler = internal.NewFileServer(static, s.notfound)
	for _, middleware := range s.middlewares {
		handler = middleware.Apply(handler)
	}
	s.Handle("GET /{path}", middleware.NewCache(time.Hour*24*365).Apply(handler))

	file, err := vite.Open(path.Join(name, ".vite/manifest.json"))
	if err != nil {
		return fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&s.manifest); err != nil {
		return fmt.Errorf("failed to decode manifest file: %w", err)
	}
	return nil
}

func (s *Server) RegisterDevHandlers() {
	proxy := httputil.NewSingleHostReverseProxy(s.options.DevServer)
	proxy.ModifyResponse = func(r *http.Response) error {
		if r.StatusCode == http.StatusNotFound {
			return internal.ErrNotFound
		}
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, internal.ErrNotFound) {
			s.notfound.ServeHTTP(w, r)
		}
	}
	var handler http.Handler = proxy
	// for _, middleware := range s.middlewares {
	// 	handler = middleware.Apply(handler)
	// }

	s.Handle("GET /{path}", handler)
	s.Handle("GET /src/", handler)
	s.Handle("GET /@vite/", handler)
	s.Handle("GET /node_modules/", handler)
}

func (s *Server) SetNotFoundPage(page Page, layout Layout) {
	s.notfound = handler(s, page, layout, true)
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.notfound = s.middlewares[i].Apply(s.notfound)
	}
}

// Asset returns the URL or path for a given asset name based on the current application context.
// In development mode, it constructs the asset URL using the development server and webfiles path.
// In production mode, it retrieves the asset entry from the application's entries map.
// If called outside of a server rendering context or if the asset is not found, the function logs a fatal error.
//
// Parameters:
//
//	ctx  - The context, expected to be of type *Context.
//	name - The name of the asset to retrieve.
//
// Returns:
//
//	The URL or path to the requested asset as a string.
func (s *Server) Asset(name string) string {
	if s.options.Env == Dev {
		return name
	}
	entry, ok := s.manifest[strings.TrimPrefix(name, "/")]
	if !ok {
		s.options.Logger.Error("Failed to find static entry", "name", name)
		return "/"
	}
	return "/" + entry.File
}

func (s *Server) Serve() {
	server := &http.Server{
		Addr:              s.options.Port,
		Handler:           s,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	logger := s.options.Logger
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("App setup failed", "error", err)
			os.Exit(1)
		}
		logger.Info("Stopped serving new connections")
	}()

	logger.Info("Server is running", "address", server.Addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	ctx, release := context.WithTimeout(context.Background(), 10*time.Second)
	defer release()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("App setup failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Server shutdowwn")
}

// New creates and returns a new Server instance using the provided Options.
// It sets default values for any missing options, including the HTTP mux, development server URL,
// environment, and logger. The function also attaches default middleware for logging and recovery,
// and sets the default "not found" page handler.
//
// Parameters:
//
//	options - a pointer to an Options struct containing configuration for the server.
//
// Returns:
//
//	A pointer to the initialized Server.
func New(options *Options) *Server {
	if options.Mux == nil {
		options.Mux = http.DefaultServeMux
	}
	if options.DevServer == nil {
		options.DevServer, _ = url.Parse("http://localhost:5173")
	}
	if len(options.Env) == 0 {
		options.Env = Dev
	}
	if options.Logger == nil {
		options.Logger = slog.Default()
	}

	s := &Server{
		ServeMux: options.Mux,
		options:  options,
	}

	s.Use(middleware.NewLogger(s.options.Logger), middleware.NewRecover(s.options.Logger, nil))
	s.SetNotFoundPage(internal.NotFoundPage, DefaultLayout)
	return s
}
