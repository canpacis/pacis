package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/canpacis/pacis/pages/pcpg/parser"
)

const filetmpl = `package {{ .Package }}

import (
{{ range .Imports }}
	{{ .Ident }} "{{ .Package }}" {{ end }}
)

{{ if ne (len .Embeds) 0 }}
var ( {{ range .Embeds }}
	//go:embed {{ .Source }}
	{{ .Ident }} embed.FS
{{ end }}
)
{{ end }}

func Router(mux *http.ServeMux) (*http.ServeMux, error) {
	if mux == nil {
		mux = http.NewServeMux()
	}

	{{ if ne (len .Assets) 0 }}
	pages.RegisterAssets(map[string]string{ {{ range $key, $value := .Assets }}
		"{{ $key }}": "{{ $value }}", {{ end }}
	})
	{{ end }}

	{{ if ne (len .NotFoundPage) 0 }}
	pages.SetNotFoundPage({{ .NotFoundPage }})
	{{ end }}
	{{ if ne (len .NotFoundPage) 0 }}
	pages.SetErrorPage({{ .ErrorPage }})
	{{ end }}

	{{ if ne (len .I18n.FS) 0 }}
	bundle, err := i18n.Setup({{ .I18n.FS }}, {{ .I18n.DefaultLang }})
	if err != nil {
		return nil, err
	}
	locale := middleware.Locale(bundle, {{ .I18n.DefaultLang }})
	{{ end }}

	{{ range .Middlewares }}
	{{ .Ident }} := {{ .Value }}
	{{ end }}

	head := html.Frag(
		html.Link(html.Rel("stylesheet"), html.Href(pages.Asset("main.css"))), 
		html.Script(html.Src(pages.Asset("app.ts"))), 
		html.Script(html.Src(pages.Asset("main.ts"))),
	)
	body := html.Frag(html.Script(html.Src(pages.Asset("stream.js"))))

	routes := []pages.Route{ {{ range .RouteStrings }}
		{{ . }}, {{ end }}
	}

	staticfs, _ := fs.Sub(static, "static")
	yearcache := middleware.Cache(time.Hour * 24 * 365)
	mux.Handle(
		"GET /static/",
		yearcache(http.StripPrefix("/static/", http.FileServerFS(staticfs))),
	)

	for _, route := range routes {
		var handler http.Handler = route
		mux.Handle(route.Path(), handler)
	}

	return mux, nil
}
`

type FileImport struct {
	Ident   string
	Package string
}

type FileEmbed struct {
	Ident  string
	Source string
}

type FileI18n struct {
	FS          string
	DefaultLang string
}

type FileRouteType string

const (
	HomeRoute     = FileRouteType("pages.NewHomeRoute")
	PageRoute     = FileRouteType("pages.NewPageRoute")
	RedirectRoute = FileRouteType("pages.NewRedirectRoute")
	RawRoute      = FileRouteType("pages.NewRawRoute")
)

type FileLayout struct {
	Name    string
	Wrapper *FileLayout
}

func (fl FileLayout) String() string {
	if fl.Wrapper == nil {
		return fl.Name
	}
	return fmt.Sprintf("pages.WrapLayout(%s, %s)", fl.Name, fl.Wrapper.Name)
}

type FileRoute struct {
	Type FileRouteType

	// Home & Page routes
	Path   string
	Page   string
	Layout FileLayout

	// Redirect routes
	Redirect string
	Code     string

	// Raw routes
	ContentType string
	Content     string

	Middlewares []string
}

type FileMiddleware struct {
	Ident string
	Value string
}

type File struct {
	Package string
	Imports []FileImport
	Embeds  []FileEmbed

	I18n   FileI18n
	Assets map[string]string

	AuthenticationHandler string

	Routes       []*FileRoute
	RouteStrings []string
	Middlewares  []FileMiddleware

	NotFoundPage string
	ErrorPage    string
}

func GenerateFile(assets map[string]string) ([]byte, error) {
	templ, err := template.New("file").Parse(filetmpl)
	if err != nil {
		return nil, err
	}

	var buf = new(bytes.Buffer)
	var layout = FileLayout{"Layout", nil}

	file := File{
		Package: "app",
		Imports: []FileImport{
			{Package: "embed"},
			{Package: "io/fs"},
			{Package: "net/http"},
			{Package: "time"},
			{Package: "github.com/canpacis/pacis/pages"},
			{Package: "github.com/canpacis/pacis/pages/i18n"},
			{Package: "github.com/canpacis/pacis/pages/middleware"},
			{Package: "github.com/canpacis/pacis/ui/html"},
		},
		Embeds: []FileEmbed{
			{"messages", "messages"},
		},
		I18n: FileI18n{
			FS:          "messages",
			DefaultLang: "language.English",
		},
		Assets: assets,
		Routes: []*FileRoute{
			{Type: HomeRoute, Page: "HomePage", Layout: layout, Middlewares: []string{"auth"}},
			{Type: PageRoute, Path: "GET /docs/{slug}", Page: "DocsPage", Layout: FileLayout{"DocLayout", &layout}, Middlewares: []string{"auth"}},
			{Type: PageRoute, Path: "GET /share/{slug}", Page: "SharePage", Layout: layout, Middlewares: []string{"auth"}},
			{Type: PageRoute, Path: "GET /auth/login", Page: "LoginPage", Layout: layout, Middlewares: []string{"auth"}},
			{Type: PageRoute, Path: "GET /auth/callback", Page: "AuthCallbackPage", Layout: layout, Middlewares: []string{"auth"}},
			{Type: RedirectRoute, Path: "GET /docs/", Redirect: "/docs/introduction", Code: "http.StatusFound"},
			{Type: RedirectRoute, Path: "GET /docs/components", Redirect: "/docs/alert", Code: "http.StatusFound"},
			{Type: RawRoute, Path: "GET /robots.txt", ContentType: "text/plain; charset=utf-8", Content: "robots"},
			{Type: RawRoute, Path: "GET /sitemap.xml", ContentType: "application/xml", Content: "sitemap"},
		},
		Middlewares: []FileMiddleware{
			{"auth", "middleware.Authentication(AuthHandler)"},
		},
		NotFoundPage: "NotFoundPage",
		ErrorPage:    "ErrorPage",
	}

	if len(file.I18n.FS) != 0 {
		file.Imports = append(
			file.Imports,
			FileImport{Package: "golang.org/x/text/language"},
		)

		for _, route := range file.Routes {
			route.Middlewares = append(route.Middlewares, "locale")
		}
	}

	for _, route := range file.Routes {
		switch route.Type {
		case PageRoute, HomeRoute:
			route.Middlewares = append(route.Middlewares, "middleware.Theme", "middleware.Gzip", "middleware.Tracer")
		case RedirectRoute:
			route.Middlewares = append(route.Middlewares, "middleware.Logger", "middleware.Tracer")
		case RawRoute:
			route.Middlewares = append(route.Middlewares, "middleware.Logger", "middleware.Gzip", "middleware.Tracer")
		}
	}

	// Always include the static embed
	file.Embeds = append(file.Embeds, FileEmbed{"static", "static"})

	for _, route := range file.Routes {
		str, err := GenerateRoute(route)
		if err != nil {
			return nil, err
		}
		file.RouteStrings = append(file.RouteStrings, str)
	}

	if err := templ.Execute(buf, file); err != nil {
		return nil, err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

const hometempl = `pages.NewHomeRoute(
	{{ .Page }},
	{{ .Layout }},
	head, body,
	{{ range .Middlewares }}
	{{ . }}, {{ end }}
)`

const pagetempl = `pages.NewPageRoute(
	"{{ .Path }}",
	{{ .Page }},
	{{ .Layout }},
	head, body,
	{{ range .Middlewares }}
	{{ . }}, {{ end }}
)`

const redirecttempl = `pages.NewRedirectRoute(
	"{{ .Path }}",
	"{{ .Redirect }}",
	{{ .Code }},
	{{ range .Middlewares }}
	{{ . }}, {{ end }}
)`

const rawtempl = `pages.NewRawRoute(
	"{{ .Path }}",
	"{{ .ContentType }}",
	{{ .Content }},
	{{ range .Middlewares }}
	{{ . }}, {{ end }}
)`

func GenerateRoute(route *FileRoute) (string, error) {

	var text string
	switch route.Type {
	case HomeRoute:
		text = hometempl
	case PageRoute:
		text = pagetempl
	case RedirectRoute:
		text = redirecttempl
	case RawRoute:
		text = rawtempl
	}

	templ, err := template.New("route").Parse(text)
	if err != nil {
		return "", err
	}

	var buf = new(bytes.Buffer)
	if err := templ.Execute(buf, route); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateRoutes(list *parser.DirectiveList) error {
	fmt.Println(list)
	return nil
}
