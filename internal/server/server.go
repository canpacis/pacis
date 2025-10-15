package server

import (
	"context"
	"net/http"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/server/metadata"
)

type AsyncChunk struct {
	ID        string
	Component html.Component
}

type RedirectMark struct {
	To     string
	Status int
}

type Context struct {
	context.Context

	ResponseWriter http.ResponseWriter
	Request        *http.Request

	AsyncChunks  []AsyncChunk
	RedirectMark *RedirectMark
	NotFoundMark bool

	Metadata *metadata.Metadata
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Context:        r.Context(),
		ResponseWriter: w,
		Request:        r,
	}
}
