package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/canpacis/pacis/pages/internal"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func Theme(h http.Handler) http.Handler {
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

		ctx := r.Context()
		ctx = internal.Set(ctx, "theme", theme)

		h.ServeHTTP(w, r.Clone(ctx))
	})
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
			ctx := r.Context()
			ctx = internal.Set(ctx, "localizer", localizer)
			ctx = internal.Set(ctx, "locale", &tag)
			h.ServeHTTP(w, r.Clone(ctx))
		})
	}
}

func Cache(duration time.Duration) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int64(duration.Seconds())))
			h.ServeHTTP(w, r)
		})
	}
}

var Gzip = gziphandler.GzipHandler

type statusRecorder struct {
	http.ResponseWriter
	http.Flusher
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := internal.Get[time.Time](r.Context(), "start")
		reqid := internal.Get[uuid.UUID](r.Context(), "request_id")

		rec := &statusRecorder{ResponseWriter: w, Flusher: w.(http.Flusher), status: http.StatusOK}
		next.ServeHTTP(rec, r)

		slog.Info(
			"request",
			slog.String("request_id", reqid.String()),
			slog.Duration("duration", time.Since(start)),
			slog.String("method", r.Method),
			slog.Int("status", rec.status),
			slog.String("path", r.URL.Path),
			slog.String("addr", r.RemoteAddr),
			slog.String("agent", r.UserAgent()),
		)
	})
}

func Tracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqid := uuid.New()

		ctx := r.Context()
		ctx = internal.Set(ctx, "request_id", reqid)
		ctx = internal.Set(ctx, "start", start)

		next.ServeHTTP(w, r.Clone(ctx))
	})
}

type User interface {
	ID() string
}

type AuthHandler[T User] func(*http.Request) (T, error)

func Authentication[T User](handler AuthHandler[T]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO
			user, _ := handler(r)
			// if err != nil {
			// 	reqid := internal.Get[uuid.UUID](r.Context(), "request_id")
			// 	logger := slog.With(
			// 		slog.String("error", err.Error()),
			// 	)
			// 	if reqid != nil {
			// 		logger = logger.With(slog.String("request_id", reqid.String()))
			// 	}
			// 	logger.Error("failed to run authentication handler")
			// }
			ctx := r.Context()
			ctx = internal.Set(ctx, "user", user)
			next.ServeHTTP(w, r.Clone(ctx))
		})
	}
}
