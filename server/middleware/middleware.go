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

func GetColorScheme(ctx context.Context) string {
	return ctx.Value(keytyp("theme")).(string)
}

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

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	return ctx.Value(keytyp("localizer")).(*i18n.Localizer)
}

func GetLocale(ctx context.Context) *language.Tag {
	return ctx.Value(keytyp("locale")).(*language.Tag)
}

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

var Gzip = gziphandler.GzipHandler
