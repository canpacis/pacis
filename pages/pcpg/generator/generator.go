package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/canpacis/pacis/pages/pcpg/parser"
	"golang.org/x/text/language"
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
	{{ if ne (len .ErrorPage) 0 }}
	pages.SetErrorPage({{ .ErrorPage }})
	{{ end }}

	{{ if ne (len .I18n.FS) 0 }}
	lang := language.MustParse("{{ .I18n.DefaultLang }}")
	bundle, err := i18n.Setup({{ .I18n.FS }}, lang)
	if err != nil {
		return nil, err
	}
	locale := middleware.Locale(bundle, lang)
	{{ end }}

	{{ range .Middlewares }}
	{{ .Ident }} := {{ .Value }}
	{{ end }}

	{{ if .HasConstituents }}
	head := html.Frag(
		html.Link(html.Rel("stylesheet"), html.Href(pages.Asset("main.css"))), 
		html.Script(html.Src(pages.Asset("app.ts"))), 
		html.Script(html.Src(pages.Asset("main.ts"))),
	)
	body := html.Frag(html.Script(html.Src(pages.Asset("stream.js"))))
	{{ end }}

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
	ActionRoute   = FileRouteType("pages.NewActionRoute")
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
	Layout *FileLayout

	// Redirect routes
	Redirect string
	Code     string

	// Action routes
	Action string

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
	HasConstituents       bool

	Routes       []*FileRoute
	RouteStrings []string
	Middlewares  []FileMiddleware

	NotFoundPage string
	ErrorPage    string
}

func GenerateFile(file *File) ([]byte, error) {
	templ, err := template.New("file").Parse(filetmpl)
	if err != nil {
		return nil, err
	}

	var buf = new(bytes.Buffer)

	if len(file.I18n.FS) != 0 {
		file.Imports = append(
			file.Imports,
			FileImport{Package: "golang.org/x/text/language"},
			FileImport{Package: "github.com/canpacis/pacis/pages/i18n"},
		)

		for _, route := range file.Routes {
			route.Middlewares = append(route.Middlewares, "locale")
		}
	}

	if file.HasConstituents {
		file.Imports = append(file.Imports, FileImport{Package: "github.com/canpacis/pacis/ui/html"})
	}

	unchecked := file.Middlewares
	file.Middlewares = []FileMiddleware{}
	for _, middleware := range unchecked {
		var used bool
		for _, route := range file.Routes {
			if slices.Contains(route.Middlewares, middleware.Ident) {
				used = true
				break
			}
		}
		if used {
			file.Middlewares = append(file.Middlewares, middleware)
		}
	}

	for _, route := range file.Routes {
		switch route.Type {
		case PageRoute, HomeRoute, ActionRoute:
			route.Middlewares = append(route.Middlewares, "middleware.Logger", "middleware.Theme", "middleware.Gzip", "middleware.Tracer")
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
	{{ if eq .Layout nil }}
	pages.EmptyLayout, {{ else }}
	{{ .Layout }}, {{ end }}
	head, body,
	{{ range .Middlewares }}
	{{ . }}, {{ end }}
)`

const pagetempl = `pages.NewPageRoute(
	"{{ .Path }}",
	{{ .Page }},
	{{ if eq .Layout nil }}
	pages.EmptyLayout, {{ else }}
	{{ .Layout }}, {{ end }}
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

const actiontempl = `pages.NewActionRoute(
	"{{ .Path }}",
	{{ .Action }},
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
	case ActionRoute:
		text = actiontempl
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

func cleanpath(p string) string {
	return strings.TrimSuffix(strings.TrimPrefix(p, "/"), "/")
}

func param(name string, dir *parser.Directive) (string, error) {
	value, ok := dir.Params[name]
	if !ok {
		return "", fmt.Errorf("%w: %s %s %s %s", parser.ErrMissingDirectiveParam, name, "parameter is required for", dir.Type.String(), "directives")
	}
	return value, nil
}

func CreateFile(list *parser.DirectiveList, assets map[string]string) (*File, error) {
	file := &File{
		Package: "app",
		Assets:  assets,
	}

	file.Imports = []FileImport{
		{Package: "embed"},
		{Package: "io/fs"},
		{Package: "net/http"},
		{Package: "time"},
		{Package: "github.com/canpacis/pacis/pages"},
		{Package: "github.com/canpacis/pacis/pages/middleware"},
	}

	for _, middleware := range list.Middleware {
		label, err := param("label", middleware)
		lclname, lclerr := param("name", middleware)

		var name string
		fn, ok := middleware.Node.(*ast.FuncDecl)
		if ok {
			name = fn.Name.String()
		} else {
			gen, ok := middleware.Node.(*ast.GenDecl)
			if !ok {
				return nil, fmt.Errorf("middleware %s is incorrectly placed, place it before a variable definition or a function", lclname)
			}
			spec, ok := gen.Specs[0].(*ast.ValueSpec)
			if !ok {
				return nil, fmt.Errorf("middleware %s is incorrectly placed, place it before a variable definition or a function", lclname)
			}
			name = spec.Names[0].Name
		}

		if err == nil {
			switch label {
			case "authentication":
				file.Middlewares = append(file.Middlewares, FileMiddleware{
					Ident: "auth",
					Value: fmt.Sprintf("middleware.Authentication(%s)", name),
				})
			default:
				return nil, fmt.Errorf("unknown middleware label: %s", label)
			}
		} else {
			if lclerr != nil {
				return nil, err
			}
			file.Middlewares = append(file.Middlewares, FileMiddleware{
				Ident: lclname,
				Value: name,
			})
		}
	}

	if len(list.Language) > 0 {
		lang := list.Language[len(list.Language)-1]
		defaultlang, err := param("default", lang)
		if err != nil {
			return nil, err
		}
		tag, err := language.Parse(defaultlang)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", parser.ErrInvalidDireciveParams, err)
		}

		file.I18n = FileI18n{
			FS:          "messages",
			DefaultLang: tag.String(),
		}
		file.Embeds = append(file.Embeds, FileEmbed{Ident: "messages", Source: "messages"})
	}

	layoutpaths := []string{}
	layouts := map[string]FileLayout{}
	for _, layout := range list.Layout {
		path, err := param("path", layout)
		if err != nil {
			return nil, err
		}
		path = "/" + cleanpath(path)

		fn, ok := layout.Node.(*ast.FuncDecl)
		if !ok {
			return nil, fmt.Errorf("layout %s is incorrectly placed, place it before a layout function", path)
		}

		layoutpaths = append(layoutpaths, path)
		layouts[path] = FileLayout{Name: fn.Name.String()}
	}
	// sort them by length, largest one is potentially the most enclosed
	sort.Slice(layoutpaths, func(i, j int) bool {
		return len(layoutpaths[i]) < len(layoutpaths[j])
	})
	// from tail to head, check every layout if the ones before it enclose it
	for i := len(layoutpaths) - 1; i >= 0; i-- {
		path := layoutpaths[i]
		// layouts that come before this one
		for j := len(layoutpaths[:i]) - 1; j >= 0; j-- {
			before := layoutpaths[j]

			if strings.HasPrefix(path, before) {
				parent := layouts[before]
				if layouts[path].Wrapper == nil {
					layouts[path] = FileLayout{
						Name:    layouts[path].Name,
						Wrapper: &parent,
					}
				} else {
					layouts[path] = FileLayout{
						Name: layouts[path].Name,
						Wrapper: &FileLayout{
							Name:    parent.Name,
							Wrapper: layouts[path].Wrapper,
						},
					}
				}
			}
		}
	}

	for _, redirect := range list.Redirect {
		to, err := param("to", redirect)
		if err != nil {
			return nil, err
		}
		from, err := param("from", redirect)
		if err != nil {
			return nil, err
		}
		file.Routes = append(file.Routes, &FileRoute{
			Type: RedirectRoute,
			// TODO: This might change if the base path changes
			Path:     "GET /" + cleanpath(from),
			Redirect: "/" + cleanpath(to),
			Code:     "http.StatusFound",
		})
	}

	for _, action := range list.Action {
		path, err := param("path", action)
		if err != nil {
			return nil, err
		}
		method, err := param("method", action)
		if err != nil {
			return nil, err
		}

		var mtd string

		switch strings.TrimSpace(strings.ToLower(method)) {
		case "post":
			mtd = "POST"
		case "patch":
			mtd = "PATCH"
		case "put":
			mtd = "PUT"
		case "delete":
			mtd = "DELETE"
		default:
			return nil, fmt.Errorf("invalid action method %s, valid methods are post, patch, put and delete", method)
		}

		middlewarelist, _ := param("middlewares", action)
		middlewares := strings.Split(middlewarelist, ",")
		middlewares = slices.DeleteFunc(middlewares, func(m string) bool {
			return len(m) == 0
		})

		fn, ok := action.Node.(*ast.FuncDecl)
		if !ok {
			return nil, fmt.Errorf("action is incorrectly placed, place it before a action function")
		}

		file.Routes = append(file.Routes, &FileRoute{
			Type:        ActionRoute,
			Path:        mtd + " /" + cleanpath(path),
			Action:      fn.Name.String(),
			Middlewares: middlewares,
		})
	}

	for _, page := range list.Page {
		label, ok := page.Params["label"]
		if ok {
			switch label {
			case "not-found":
				fn, ok := page.Node.(*ast.FuncDecl)
				if !ok {
					return nil, fmt.Errorf("not-found is incorrectly placed, place it before a page function")
				}
				file.NotFoundPage = fn.Name.String()
			case "error":
				fn, ok := page.Node.(*ast.FuncDecl)
				if !ok {
					return nil, fmt.Errorf("error is incorrectly placed, place it before a page function")
				}
				file.ErrorPage = fn.Name.String()
			case "robots":
				decl := page.Node.(*ast.GenDecl)
				spec, ok := decl.Specs[0].(*ast.ValueSpec)
				if !ok {
					return nil, fmt.Errorf("robots directive is incorrectly placed, place it before a variable declaration")
				}
				if len(spec.Names) == 0 {
					return nil, fmt.Errorf("robots directive needs to be put before a named variable")
				}

				file.Routes = append(file.Routes, &FileRoute{
					Type:        RawRoute,
					Path:        "GET /robots.txt",
					ContentType: "text/plain; charset=utf-8",
					Content:     spec.Names[0].String(),
				})
			case "sitemap":
				decl := page.Node.(*ast.GenDecl)
				spec, ok := decl.Specs[0].(*ast.ValueSpec)
				if !ok {
					return nil, fmt.Errorf("sitemap directive is incorrectly placed, place it before a variable declaration")
				}
				if len(spec.Names) == 0 {
					return nil, fmt.Errorf("sitemap directive needs to be put before a named variable")
				}

				file.Routes = append(file.Routes, &FileRoute{
					Type:        RawRoute,
					Path:        "GET /sitemap.xml",
					ContentType: "application/xml",
					Content:     spec.Names[0].String(),
				})
			default:
				return nil, fmt.Errorf("unknown page label: %s", label)
			}
		} else {
			path, err := param("path", page)
			if err != nil {
				return nil, err
			}
			path = "/" + cleanpath(path)

			var layout *FileLayout
			var currlpath string
			for _, lpath := range layoutpaths {
				if strings.HasPrefix(path, lpath) {
					currlpath = lpath
				}
			}
			if len(currlpath) > 0 {
				l := layouts[currlpath]
				layout = &l
			}

			fn, ok := page.Node.(*ast.FuncDecl)
			if !ok {
				return nil, fmt.Errorf("page %s is incorrectly placed, place it before a page function", path)
			}

			var typ FileRouteType = PageRoute
			if path == "/" {
				typ = HomeRoute
			}

			middlewarelist, _ := param("middlewares", page)
			middlewares := strings.Split(middlewarelist, ",")
			middlewares = slices.DeleteFunc(middlewares, func(m string) bool {
				return len(m) == 0
			})

			file.Routes = append(file.Routes, &FileRoute{
				Type:        typ,
				Path:        "GET " + path,
				Page:        fn.Name.String(),
				Layout:      layout,
				Middlewares: middlewares,
			})
			file.HasConstituents = true
		}
	}

	return file, nil
}
