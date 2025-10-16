package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
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

	"github.com/canpacis/pacis/html"
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

type Options struct {
	Env       Environment
	Host      string
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
}

// Adds middleware(s) to the application's middleware stack.
func (s *Server) Use(middlewares ...middleware.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) HandlePage(pattern string, page Page, layout Layout, middlewares ...middleware.Middleware) {
	pattern = strings.TrimSuffix(pattern, "/") + "/"
	s.Handle(pattern, HandlerOf(s, page, layout, middlewares...))
}

func (s *Server) SetBuildDir(name string, dir fs.FS, vite fs.FS) error {
	static, err := fs.Sub(dir, name)
	if err != nil {
		return err
	}

	s.Handle("GET /{path}", middleware.NewCache(time.Hour*24*365).Apply(http.FileServerFS(static)))

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
	s.Handle("GET /{path}", proxy)
	s.Handle("GET /src/", proxy)
	s.Handle("GET /@vite/", proxy)
	s.Handle("GET /node_modules/", proxy)
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
		log.Fatalf("Failed to find static entry: %s", name)
	}
	return "/" + entry.File
}

func (s *Server) HMR() html.Node {
	if s.options.Env == Prod {
		return html.Fragment()
	}
	return html.Script(html.Type("module"), html.Src("/@vite/client"))
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
	return s
}
