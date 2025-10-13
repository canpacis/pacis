package server

import (
	"context"
	"log"
	"net/http"
	"net/url"

	payload "github.com/canpacis/http-payload"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/internal/util"
)

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
func Asset(app *App, name string) string {
	if app.options.env == Dev {
		return app.options.devserver + "/" + app.options.webfiles + "/" + name
	}
	entry, ok := app.entries[name]
	if !ok {
		log.Fatalf("failed to retrieve asset %s", name)
	}
	return entry
}

func HMR(app *App) html.Node {
	if app.options.env == Prod {
		return html.Fragment()
	}
	return html.Script(html.Type("module"), html.Src(app.options.devserver+"/@vite/client"))
}

func Async(comp html.Component, fallback html.Node) html.Node {
	id := util.PrefixedID("pacis")
	if fallback == nil {
		fallback = html.Fragment()
	}

	return html.Component(func(ctx context.Context) html.Node {
		serverctx, ok := ctx.(*serverctx)
		if ok {
			serverctx.async = append(serverctx.async, async{
				id:   id,
				comp: comp,
			})
		}
		return html.Slot(html.Name(id), fallback)
	})
}

func Data[T any](ctx context.Context) (*T, error) {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("Data helper used outside of server rendering context")
	}
	req := context.req.Clone(ctx)
	data := new(T)
	scanner := payload.NewPipeScanner(
		payload.NewPathScanner(req),
		payload.NewQueryScanner(req.URL.Query()),
		payload.NewCookieScanner(req.Cookies()),
		payload.NewHeaderScanner(&req.Header),
	)
	if err := scanner.Scan(data); err != nil {
		return nil, err
	}
	return data, nil
}

type RequestDetail struct {
	URL *url.URL
	Host,
	Method,
	Pattern,
	RemoteAddr,
	RequestURI string
	Cookies []*http.Cookie
}

func Detail(ctx context.Context) *RequestDetail {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("Detail helper used outside of server rendering context")
	}
	return &RequestDetail{
		URL:        context.req.URL,
		Host:       context.req.Host,
		Method:     context.req.Method,
		Pattern:    context.req.Pattern,
		RemoteAddr: context.req.RemoteAddr,
		RequestURI: context.req.RequestURI,
		Cookies:    context.req.Cookies(),
	}
}

func Redirect(ctx context.Context, to string) html.Node {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("Redirect node used outside of server rendering context")
	}
	context.redirect = &redirect{status: http.StatusFound, to: to}
	return html.Fragment()
}

func RedirectWith(ctx context.Context, to string, status int) html.Node {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("RedirectWith node used outside of server rendering context")
	}
	context.redirect = &redirect{status: status, to: to}
	return html.Fragment()
}
