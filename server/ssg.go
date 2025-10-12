package server

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync"

	"github.com/canpacis/pacis/html"
)

type StaticRenderer struct {
	chunks []any
}

var bufpool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func (r *StaticRenderer) Build(node html.Node) error {
	ctx := context.Background()
	buf := bufpool.New().(*bytes.Buffer)
	defer bufpool.Put(buf)

	for chunk := range node.Chunks() {
		if chunk.IsPure() {
			if err := chunk.Render(ctx, buf); err != nil {
				return err
			}
		} else {
			r.chunks = append(r.chunks, buf.Bytes())
			bufpool.Put(buf)
			buf = bufpool.New().(*bytes.Buffer)
			r.chunks = append(r.chunks, chunk.Render)
		}
	}
	r.chunks = append(r.chunks, buf.Bytes())

	return nil
}

func (r *StaticRenderer) Render(ctx context.Context, w io.Writer) error {
	for _, chunk := range r.chunks {
		switch chunk := chunk.(type) {
		case []byte:
			if _, err := w.Write(chunk); err != nil {
				return err
			}
		case func(context.Context, io.Writer) error:
			if err := chunk(ctx, w); err != nil {
				return err
			}
		default:
			log.Fatalf("Invalid static chunk type: %T", chunk)
		}
	}
	return nil
}
