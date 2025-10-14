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

func RedirectComponent(to string) html.Component {
	return func(ctx context.Context) html.Node {
		return Redirect(ctx, to)
	}
}

func RedirectWith(ctx context.Context, to string, status int) html.Node {
	context, ok := ctx.(*serverctx)
	if !ok {
		log.Fatal("RedirectWith node used outside of server rendering context")
	}
	context.redirect = &redirect{status: status, to: to}
	return html.Fragment()
}

func RedirectWithComponent(to string, status int) html.Component {
	return func(ctx context.Context) html.Node {
		return RedirectWith(ctx, to, status)
	}
}
