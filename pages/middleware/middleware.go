package middleware

import (
	"net/http"
	"strings"

	"github.com/canpacis/pacis/pages"
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

		h.ServeHTTP(w, r.Clone(pages.Set(r.Context(), "theme", theme)))
	})
}

func Locale(bundle *i18n.Bundle, defaultlang language.Tag) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var locale string
			cookie, err := r.Cookie("pacis_app_locale")
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
			r = r.Clone(pages.Set(r.Context(), "localizer", localizer))
			r = r.Clone(pages.Set(r.Context(), "locale", &tag))
			h.ServeHTTP(w, r)
		})
	}
}
