package server_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/internal"
	"github.com/canpacis/pacis/server"
	"github.com/stretchr/testify/assert"
)

type ResponseWriter struct {
	header http.Header
	buf    *bytes.Buffer
	status int
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: make(http.Header),
		buf:    new(bytes.Buffer),
	}
}

func TestRedirect(t *testing.T) {
	assert := assert.New(t)

	rw := NewResponseWriter()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	ctx := internal.NewContext(rw, r)
	server.Redirect(ctx, "/redirect")

	assert.NotNil(ctx.RedirectMark)
	assert.Equal(http.StatusFound, ctx.RedirectMark.Status)
	assert.Equal("/redirect", ctx.RedirectMark.To)
}

func TestNotFound(t *testing.T) {
	assert := assert.New(t)

	rw := NewResponseWriter()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	ctx := internal.NewContext(rw, r)
	server.NotFound(ctx)

	assert.True(ctx.NotFoundMark)
}

func TestSetCookie(t *testing.T) {
	assert := assert.New(t)

	rw := NewResponseWriter()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	ctx := internal.NewContext(rw, r)
	server.SetCookie(ctx, &http.Cookie{Name: "cookie", Value: "value"})

	assert.Equal("cookie=value", rw.header.Get("Set-Cookie"))
}

func TestAsync(t *testing.T) {
	assert := assert.New(t)

	rw := NewResponseWriter()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	ctx := internal.NewContext(rw, r)

	component := func(ctx context.Context) html.Node {
		return html.Fragment()
	}
	server.Async(component, nil)(ctx)
	assert.Equal(1, len(ctx.AsyncChunks))
}
