package server

import (
	"context"
	"log"
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
		log.Fatal("Data helper used outside of server rendering context")
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
	context, ok := ctx.(*server.Context)
	if !ok {
		log.Fatal("Detail helper used outside of server rendering context")
	}
	return &RequestDetail{
		URL:        context.Request.URL,
		Host:       context.Request.Host,
		Method:     context.Request.Method,
		Pattern:    context.Request.Pattern,
		RemoteAddr: context.Request.RemoteAddr,
		RequestURI: context.Request.RequestURI,
		Cookies:    context.Request.Cookies(),
	}
}

func Redirect(ctx context.Context, to string) html.Node {
	context, ok := ctx.(*server.Context)
	if !ok {
		log.Fatal("Redirect node used outside of server rendering context")
	}
	context.RedirectMark = &server.RedirectMark{Status: http.StatusFound, To: to}
	return html.Fragment()
}

func RedirectComponent(to string) html.Component {
	return func(ctx context.Context) html.Node {
		return Redirect(ctx, to)
	}
}

func RedirectWith(ctx context.Context, to string, status int) html.Node {
	context, ok := ctx.(*server.Context)
	if !ok {
		log.Fatal("RedirectWith node used outside of server rendering context")
	}
	context.RedirectMark = &server.RedirectMark{Status: status, To: to}
	return html.Fragment()
}

func RedirectWithComponent(to string, status int) html.Component {
	return func(ctx context.Context) html.Node {
		return RedirectWith(ctx, to, status)
	}
}

func NotFound(ctx context.Context) html.Node {
	context, ok := ctx.(*server.Context)
	if !ok {
		log.Fatal("NotFound node used outside of server rendering context")
	}
	context.NotFoundMark = true
	return html.Fragment()
}

func SetCookie(ctx context.Context, cookie *http.Cookie) html.Node {
	context, ok := ctx.(*server.Context)
	if !ok {
		log.Fatal("SetCookie node used outside of server rendering context")
	}
	http.SetCookie(context.ResponseWriter, cookie)
	return html.Fragment()
}

func SetCookieComponent(cookie *http.Cookie) html.Component {
	return func(ctx context.Context) html.Node {
		return SetCookie(ctx, cookie)
	}
}
