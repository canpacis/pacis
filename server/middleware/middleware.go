// Package middleware provides HTTP middleware utilities for handling
// color scheme preferences, locale selection, caching, logging, and gzip compression.
//
// The package includes the following middleware:
//
//   - ColorScheme: Detects and sets the user's color scheme preference (light or dark) via cookies.
//   - Locale: Determines the user's locale from cookies, query parameters, or Accept-Language headers,
//     and injects an i18n.Localizer into the request context.
//   - Cache: Sets Cache-Control headers for HTTP responses to enable client-side caching.
//   - Logger: Logs HTTP requests with method, status, path, remote address, user agent, and duration.
//   - Gzip: Provides gzip compression for HTTP responses.
//
// Helper functions are provided to retrieve the color scheme, localizer, and locale from the request context.
package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type KeyType string

type Middleware interface {
	Name() string
	Apply(http.Handler) http.Handler
}

// ColorScheme is a middleware that manages the user's color scheme preference via a cookie.
// It checks for the "pacis_color_scheme" cookie in the incoming request. If the cookie exists
// and its value is "light" or "dark", it uses that value as the theme. Otherwise, it defaults
// to "light" and sets the cookie accordingly. The selected theme is stored in the request's
// context under the key "theme" for downstream handlers to access.
type ColorScheme struct {
	key string
}

func (*ColorScheme) Name() string {
	return "ColorScheme"
}

func (m *ColorScheme) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		set := func(theme string) {
			http.SetCookie(w, &http.Cookie{
				Name:     m.key,
				Value:    theme,
				Path:     "/",
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
		}

		cookie, err := r.Cookie(m.key)
		var scheme string

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				header := r.Header.Get("Sec-CH-Prefers-Color-Scheme")
				if len(header) > 0 {
					switch header {
					case "light", "dark":
						scheme = header
					default:
						scheme = "light"
					}
				} else {
					w.Header().Set("Accept-CH", "Sec-CH-Prefers-Color-Scheme")
					w.Header().Set("Critical-CH", "Sec-CH-Prefers-Color-Scheme")
					w.Header().Set("Vary", "Sec-CH-Prefers-Color-Scheme")
				}
			} else {
				scheme = "light"
				set(scheme)
			}
		} else {
			switch cookie.Value {
			case "light", "dark":
				scheme = cookie.Value
			default:
				scheme = "light"
				set(scheme)
			}
		}

		ctx := context.WithValue(r.Context(), KeyType("color-scheme"), scheme)
		h.ServeHTTP(w, r.Clone(ctx))
	})
}

var DefaultColorScheme = &ColorScheme{key: "pacis_color_scheme"}

func NewColorScheme(key string) *ColorScheme {
	return &ColorScheme{key: key}
}

// GetColorScheme retrieves the color scheme (theme) from the provided context.
// It expects the context to have a value associated with the key "theme" of type string.
// If the value is not present or not a string, this function will panic.
func GetColorScheme(ctx context.Context) string {
	return ctx.Value(KeyType("color-scheme")).(string)
}

// Locale is a middleware that determines the user's preferred language from a cookie ("pacis_locale"),
// a "lang" form value, or the "Accept-Language" HTTP header, falling back to the provided default language
// if none are set or valid. It parses the locale, creates an i18n.Localizer, and injects both the localizer
// and the language tag into the request context for downstream handlers to use.
type Locale struct {
	key         string
	bundle      *i18n.Bundle
	defaultlang language.Tag
}

func (*Locale) Name() string {
	return "Locale"
}

func (l *Locale) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var locale string
		cookie, err := r.Cookie(l.key)
		if err == nil {
			locale = cookie.Value
		} else {
			header := r.Header.Get("Accept-Language")
			switch {
			case len(r.FormValue("lang")) > 0:
				locale = r.FormValue("lang")
			case len(header) > 0:
				locale = strings.Split(header, ",")[0]
			default:
				locale = l.defaultlang.String()
			}
		}
		tag, err := language.Parse(locale)
		if err != nil {
			tag = l.defaultlang
		}

		localizer := i18n.NewLocalizer(l.bundle, tag.String())
		ctx := context.WithValue(r.Context(), KeyType("localizer"), localizer)
		ctx = context.WithValue(ctx, KeyType("locale"), &tag)

		h.ServeHTTP(w, r.Clone(ctx))
	})
}

func NewLocale(key string, bundle *i18n.Bundle, defaultlang language.Tag) *Locale {
	return &Locale{key: key, bundle: bundle, defaultlang: defaultlang}
}

// GetLocalizer retrieves the localizer struct from the provided context.
// It expects the context to have a value associated with the key "localizer" of type *i18n.Localizer
// If the value is not present or not *i18n.Localizer, this function will panic.
func GetLocalizer(ctx context.Context) *i18n.Localizer {
	return ctx.Value(KeyType("localizer")).(*i18n.Localizer)
}

// GetLocale retrieves the locale value from the provided context.
// It expects the context to have a value associated with the key "locale" of type *language.Tag
// If the value is not present or not *language.Tag, this function will panic.
func GetLocale(ctx context.Context) *language.Tag {
	return ctx.Value(KeyType("locale")).(*language.Tag)
}

// Cache returns a middleware that sets the "Cache-Control" header on HTTP responses,
// specifying that the response can be cached by any cache and defining the maximum age
// (in seconds) that the response is considered fresh. The duration parameter determines
// the max-age value. This middleware should be used to control client and proxy caching
// behavior for HTTP handlers.
type Cache struct {
	duration time.Duration
}

func (*Cache) Name() string {
	return "Cache"
}

func (c *Cache) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int64(c.duration.Seconds())))
		h.ServeHTTP(w, r)
	})
}

func NewCache(dur time.Duration) *Cache {
	return &Cache{duration: dur}
}

type logwriter struct {
	http.ResponseWriter
	http.Flusher
	status int
}

func (lw *logwriter) WriteHeader(code int) {
	lw.status = code
	lw.ResponseWriter.WriteHeader(code)
}

// Logger returns a middleware that logs HTTP requests using the provided slog.Logger.
// It records the request method, status code, path, remote address, user agent, and duration.
// The middleware wraps the next http.Handler and logs the request details after it is served.
//
// Example usage:
//
//	mux.Use(Logger(logger))
//
// Parameters:
//
//	logger - the slog.Logger instance used for logging request details.
//
// Returns:
//
//	A middleware function compatible with http.Handler.
type Logger struct {
	logger *slog.Logger
}

func (*Logger) Name() string {
	return "Logger"
}

func (l *Logger) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &logwriter{ResponseWriter: w, Flusher: w.(http.Flusher), status: http.StatusOK}
		h.ServeHTTP(lw, r)

		agent := strings.Split(r.UserAgent(), " ")[0]
		l.logger.Info(
			"request",
			slog.Duration("duration", time.Since(start)),
			slog.String("method", r.Method),
			slog.Int("status", lw.status),
			slog.String("path", r.URL.Path),
			slog.String("addr", r.RemoteAddr),
			slog.String("agent", agent),
		)
	})
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger: logger}
}

// GzipHandler wraps an HTTP handler, to transparently gzip the response body if the client supports
// it (via the Accept-Encoding header). This will compress at the default compression level.
type Gzip struct{}

func (*Gzip) Name() string {
	return "Gzip"
}

func (*Gzip) Apply(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

var DefaultGzip = &Gzip{}
