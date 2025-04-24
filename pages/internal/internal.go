package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/scanner/structd"
)

type ctxkey string

func Set(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, ctxkey(fmt.Sprintf("%s:%s", "pacis", key)), value)
}

func Get[T any](ctx context.Context, key string) T {
	value := ctx.Value(ctxkey(fmt.Sprintf("%s:%s", "pacis", key)))
	cast, ok := value.(T)
	if !ok {
		var v T
		log.Fatalf("failed to cast ctx key '%s' to %T\n", key, v)
		return v
	}
	return cast
}

type ContextScanner struct {
	ctx context.Context
}

func (s *ContextScanner) Scan(v any) error {
	return structd.New(s, "context").Decode(v)
}

func (s *ContextScanner) Get(key string) any {
	return Get[any](s.ctx, key)
}

func NewContextScanner(ctx context.Context) *ContextScanner {
	return &ContextScanner{ctx: ctx}
}

type StreamWriter struct {
	Renderer html.I

	buf       *bytes.Buffer
	chunksize int
	w         http.ResponseWriter
	f         http.Flusher
}

func (s *StreamWriter) Write(p []byte) (int, error) {
	if s.buf.Len() < s.chunksize {
		return s.buf.Write(p)
	}

	n, err := s.buf.Write(p)
	if err != nil {
		return n, err
	}
	m, err := io.Copy(s.w, s.buf)
	s.f.Flush()
	s.buf.Reset()
	return int(m), err
}

func (s *StreamWriter) Flush() {
	if s.buf.Len() != 0 {
		io.Copy(s.w, s.buf)
		s.f.Flush()
		s.buf.Reset()
	}
}

func NewStreamWriter(renderer html.I, w http.ResponseWriter) *StreamWriter {
	return &StreamWriter{
		Renderer:  renderer,
		buf:       new(bytes.Buffer),
		chunksize: 500,
		w:         w,
		f:         w.(http.Flusher),
	}
}

func Render(ctx context.Context, sw *StreamWriter) error {
	sw.w.Header().Set("Content-Type", "text/html")
	if err := sw.Renderer.Render(ctx, sw); err != nil {
		return err
	}
	sw.Flush()

	// TODO: Async streaming
	// // TODO: Maybe carry this to the hooks api
	// size := int(ctx.chsize.Load())
	// if size == 0 {
	// 	return
	// }

	// ctx.elemch = make(chan h.Element, size)
	// ctx.ready.Store(true)

	// for range size {
	// 	select {
	// 	case <-ctx.Done():
	// 		// client disconnected
	// 		return
	// 	case el := <-ctx.elemch:
	// 		el.Render(ctx, w)
	// 		s.Flush()
	// 	}
	// }
	return nil
}
