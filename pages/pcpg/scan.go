package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path"
	"strings"
)

type directivekind string

const (
	layoutdir         = directivekind("layout")
	pagedir           = directivekind("page")
	redirectdir       = directivekind("redirect")
	languagedir       = directivekind("language")
	authenticationdir = directivekind("authentication")
	middlewaredir     = directivekind("middleware")
)

type directive struct {
	kind   directivekind
	params map[string]string
}

func parsedirective(src string) (*directive, error) {
	rest, _ := strings.CutPrefix(src, "//pacis:")
	parts := strings.Split(rest, " ")
	if len(parts) == 0 {
		return nil, errors.New("malformed pacis directive")
	}
	dir := &directive{
		params: map[string]string{},
	}

	switch parts[0] {
	case string(layoutdir):
		dir.kind = layoutdir
	case string(pagedir):
		dir.kind = pagedir
	case string(redirectdir):
		dir.kind = redirectdir
	case string(languagedir):
		dir.kind = languagedir
	case string(authenticationdir):
		dir.kind = authenticationdir
	case string(middlewaredir):
		dir.kind = middlewaredir
	default:
		return nil, fmt.Errorf("invalid pacis directive kind %s", parts[0])
	}

	for _, part := range parts[1:] {
		split := strings.Split(part, "=")
		if len(split) < 2 {
			return nil, fmt.Errorf("malformed pacis directive param %s", part)
		}
		dir.params[split[0]] = split[1]
	}

	return dir, nil
}

type def struct {
	dir  *directive
	decl ast.Decl
}

func scanfile(decls []ast.Decl) ([]def, error) {
	defs := []def{}
	for _, decl := range decls {
		fndecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			// TODO: Check var decl
			continue
		}
		if fndecl.Doc == nil {
			continue
		}

		for _, comment := range fndecl.Doc.List {
			dir, err := parsedirective(comment.Text)
			if err != nil {
				return nil, err
			}

			defs = append(defs, def{dir, decl})
		}
	}
	return defs, nil
}

func scan(target string, assetmap map[string]string) (*generator, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root := path.Join(wd, target)
	app := path.Join(root, "app")

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, app, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	pkg, ok := pkgs["app"]
	if !ok {
		return nil, errors.New("failed to locate the app package")
	}

	gen := &generator{
		deps: []string{
			"embed", "io/fs", "log", "net/http", "time",
			"github.com/canpacis/pacis/pages",
			"github.com/canpacis/pacis/pages/i18n",
			"github.com/canpacis/pacis/pages/middleware",
			"github.com/canpacis/pacis/ui/html",
			"golang.org/x/text/language",
		},
		embeds: []string{"messages", "static"},
		assets: map[string]string{},
		head: []struct {
			path string
			typ  string
		}{
			{`pages.Asset("main.css")`, "stylesheet"},
			{`pages.Asset("app.ts")`, "javascript"},
			{`pages.Asset("main.ts")`, "javascript"},
		},
		body: []struct {
			path string
			typ  string
		}{
			{`pages.Asset("stream.js")`, "javascript"},
		},
		routes: []route{
			home("HomePage", "Layout", "middleware.Logger", "locale", "auth", "middleware.Theme", "middleware.Gzip", "middleware.Tracer"),
			page("GET /docs/{slug}", "DocsPage", "pages.WrapLayout(DocLayout, Layout)", "middleware.Logger", "locale", "auth", "middleware.Theme", "middleware.Gzip", "middleware.Tracer"),
			page("GET /auth/login", "LoginPage", "Layout", "middleware.Logger", "locale", "auth", "middleware.Theme", "middleware.Gzip", "middleware.Tracer"),
			page("GET /auth/logout", "LogoutPage", "Layout", "middleware.Logger", "locale", "auth", "middleware.Theme", "middleware.Gzip", "middleware.Tracer"),
			page("GET /auth/callback", "AuthCallbackPage", "Layout", "middleware.Logger", "locale", "auth", "middleware.Theme", "middleware.Gzip", "middleware.Tracer"),
			redirect("GET /docs/", "/docs/introduction", "http.StatusFound", "middleware.Logger", "middleware.Tracer"),
			redirect("GET /docs/components", "/docs/alert", "http.StatusFound", "middleware.Logger", "middleware.Tracer"),
			raw("GET /robots.txt", "robots", "text/plain; charset=utf-8", "middleware.Logger", "middleware.Gzip"),
			raw("GET /sitemap.xml", "sitemap", "application/xml", "middleware.Logger", "middleware.Gzip"),
		},
	}

	maps.Copy(gen.assets, assetmap)

	for _, file := range pkg.Files {
		_, err := scanfile(file.Decls)
		if err != nil {
			return nil, err
		}
	}

	return gen, nil
}
