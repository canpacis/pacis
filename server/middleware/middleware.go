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
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type keytyp string

// ColorScheme is a middleware that manages the user's color scheme preference via a cookie.
// It checks for the "pacis_color_scheme" cookie in the incoming request. If the cookie exists
// and its value is "light" or "dark", it uses that value as the theme. Otherwise, it defaults
// to "light" and sets the cookie accordingly. The selected theme is stored in the request's
// context under the key "theme" for downstream handlers to access.
func ColorScheme(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		set := func(theme string) {
			http.SetCookie(w, &http.Cookie{
				Name:     "pacis_color_scheme",
				Value:    theme,
				Path:     "/",
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
		}

		cookie, err := r.Cookie("pacis_color_scheme")
		var theme string
		if err == nil {
			switch cookie.Value {
			case "light", "dark":
				theme = cookie.Value
			default:
				theme = "light"
				set(theme)
			}
		} else {
			theme = "light"
			set(theme)
		}

		ctx := context.WithValue(r.Context(), keytyp("theme"), theme)

		h.ServeHTTP(w, r.Clone(ctx))
	})
}

// GetColorScheme retrieves the color scheme (theme) from the provided context.
// It expects the context to have a value associated with the key "theme" of type string.
// If the value is not present or not a string, this function will panic.
func GetColorScheme(ctx context.Context) string {
	return ctx.Value(keytyp("theme")).(string)
}

// Locale is a middleware that determines the user's preferred language from a cookie ("pacis_locale"),
// a "lang" form value, or the "Accept-Language" HTTP header, falling back to the provided default language
// if none are set or valid. It parses the locale, creates an i18n.Localizer, and injects both the localizer
// and the language tag into the request context for downstream handlers to use.
func Locale(bundle *i18n.Bundle, defaultlang language.Tag) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var locale string
			cookie, err := r.Cookie("pacis_locale")
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
					locale = defaultlang.String()
				}
			}
			tag, err := language.Parse(locale)
			if err != nil {
				tag = defaultlang
			}

			localizer := i18n.NewLocalizer(bundle, tag.String())
			ctx := context.WithValue(r.Context(), keytyp("localizer"), localizer)
			ctx = context.WithValue(ctx, keytyp("locale"), &tag)

			h.ServeHTTP(w, r.Clone(ctx))
		})
	}
}

// GetLocalizer retrieves the localizer struct from the provided context.
// It expects the context to have a value associated with the key "localizer" of type *i18n.Localizer
// If the value is not present or not *i18n.Localizer, this function will panic.
func GetLocalizer(ctx context.Context) *i18n.Localizer {
	return ctx.Value(keytyp("localizer")).(*i18n.Localizer)
}

// GetLocale retrieves the locale value from the provided context.
// It expects the context to have a value associated with the key "locale" of type *language.Tag
// If the value is not present or not *language.Tag, this function will panic.
func GetLocale(ctx context.Context) *language.Tag {
	return ctx.Value(keytyp("locale")).(*language.Tag)
}

// Cache returns a middleware that sets the "Cache-Control" header on HTTP responses,
// specifying that the response can be cached by any cache and defining the maximum age
// (in seconds) that the response is considered fresh. The duration parameter determines
// the max-age value. This middleware should be used to control client and proxy caching
// behavior for HTTP handlers.
func Cache(duration time.Duration) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int64(duration.Seconds())))
			h.ServeHTTP(w, r)
		})
	}
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
func Logger(logger *slog.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &logwriter{ResponseWriter: w, Flusher: w.(http.Flusher), status: http.StatusOK}
			next.ServeHTTP(rec, r)

			agent := strings.Split(r.UserAgent(), " ")[0]
			logger.Info(
				"request",
				slog.Duration("duration", time.Since(start)),
				slog.String("method", r.Method),
				slog.Int("status", rec.status),
				slog.String("path", r.URL.Path),
				slog.String("addr", r.RemoteAddr),
				slog.String("agent", agent),
			)
		})
	}
}

// GzipHandler wraps an HTTP handler, to transparently gzip the response body if the client supports it (via the Accept-Encoding header). This will compress at the default compression level.
var Gzip = gziphandler.GzipHandler
