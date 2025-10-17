package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	payload "github.com/canpacis/http-payload"
	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/internal/server"
	"github.com/canpacis/pacis/internal/util"
)

func Async(comp html.Component, fallback html.Node) html.Component {
	id := util.PrefixedID("pacis")
	if fallback == nil {
		fallback = html.Fragment()
	}

	return html.Component(func(ctx context.Context) html.Node {
		context, ok := ctx.(*server.Context)
		if ok {
			context.AsyncChunks = append(context.AsyncChunks, server.AsyncChunk{
				ID:        id,
				Component: comp,
			})
		}
		return html.Slot(html.Name(id), fallback)
	})
}

func Data[T any](ctx context.Context) (*T, error) {
	context, ok := ctx.(*server.Context)
	if !ok {
		return nil, fmt.Errorf("Data helper used outside of server rendering context")
	}
	req := context.Request.Clone(ctx)
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

func FormData[T any](r *http.Request) (*T, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	data := new(T)
	if err := payload.NewFormScanner(&r.PostForm).Scan(data); err != nil {
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

func Detail(ctx context.Context) (*RequestDetail, error) {
	context, ok := ctx.(*server.Context)
	if !ok {
		return nil, fmt.Errorf("Detail helper used outside of server rendering context")
	}
	return &RequestDetail{
		URL:        context.Request.URL,
		Host:       context.Request.Host,
		Method:     context.Request.Method,
		Pattern:    context.Request.Pattern,
		RemoteAddr: context.Request.RemoteAddr,
		RequestURI: context.Request.RequestURI,
		Cookies:    context.Request.Cookies(),
	}, nil
}

func Redirect(ctx context.Context, to string) html.Node {
	context, ok := ctx.(*server.Context)
	if ok {
		context.RedirectMark = &server.RedirectMark{Status: http.StatusFound, To: to}
	} else {
		slog.Error("Redirect node used outside of server rendering context")
	}
	return html.Fragment()
}

func RedirectComponent(to string) html.Component {
	return func(ctx context.Context) html.Node {
		return Redirect(ctx, to)
	}
}

func RedirectWith(ctx context.Context, to string, status int) html.Node {
	context, ok := ctx.(*server.Context)
	if ok {
		context.RedirectMark = &server.RedirectMark{Status: status, To: to}
	} else {
		slog.Error("RedirectWith node used outside of server rendering context")
	}
	return html.Fragment()
}

func RedirectWithComponent(to string, status int) html.Component {
	return func(ctx context.Context) html.Node {
		return RedirectWith(ctx, to, status)
	}
}

func NotFound(ctx context.Context) html.Node {
	context, ok := ctx.(*server.Context)
	if ok {
		context.NotFoundMark = true
	} else {
		slog.Error("NotFound node used outside of server rendering context")
	}
	return html.Fragment()
}

func SetCookie(ctx context.Context, cookie *http.Cookie) html.Node {
	context, ok := ctx.(*server.Context)
	if ok {
		http.SetCookie(context.ResponseWriter, cookie)
	} else {
		slog.Error("SetCookie node used outside of server rendering context")
	}
	return html.Fragment()
}

func SetCookieComponent(cookie *http.Cookie) html.Component {
	return func(ctx context.Context) html.Node {
		return SetCookie(ctx, cookie)
	}
}

func Form(name string, items ...html.Item) html.Node {
	return html.Form(append(items, html.Method("POST"), html.Action("?__action="+name))...)
}
