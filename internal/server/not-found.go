package server

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/server/metadata"
)

type NotFound struct{}

func (*NotFound) Metadata() *metadata.Metadata {
	return &metadata.Metadata{
		Title: "Page Not Found",
	}
}

func (*NotFound) Page() html.Node {
	return html.Div(
		html.StyleAttr("min-height: 100dvh; display: flex; flex-direction: column; justify-content: center; place-items: center; font-family: system-ui, sans-serif;"),

		html.H1(
			html.StyleAttr("font-size: 2.25rem; font-weight: bold; margin: 1rem 0;"),

			html.Text("404"),
		),
		html.P(
			html.StyleAttr("margin: 0;"),

			html.Text("Page Not Found"),
		),
	)
}

var NotFoundPage = &NotFound{}

var ErrNotFound = errors.New("http not found")

type FileServer struct {
	fs       fs.FS
	handler  http.Handler
	notfound http.Handler
}

func (h *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clean := path.Clean("/" + r.URL.Path) // ensure at least "/"
	trimmed := strings.TrimPrefix(clean, "/")

	_, err := h.fs.Open(trimmed)
	if errors.Is(err, fs.ErrNotExist) {
		h.notfound.ServeHTTP(w, r)
	} else {
		h.handler.ServeHTTP(w, r)
	}
}

func NewFileServer(fs fs.FS, notfound http.Handler) *FileServer {
	return &FileServer{fs: fs, handler: http.FileServerFS(fs), notfound: notfound}
}
