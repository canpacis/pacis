package server

import (
	"context"
	"log"
	"net/http"

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
