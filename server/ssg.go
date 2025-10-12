package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/canpacis/pacis/html"
)

var chunkpool = sync.Pool{
	New: func() any {
		return &[]any{}
	},
}

type StaticRenderer struct {
	chunks *[]any
}

var bufpool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func (r *StaticRenderer) Build(node html.Node) error {
	buf := bufpool.New().(*bytes.Buffer)
	defer bufpool.Put(buf)

	for chunk := range node.Chunks() {
		switch chunk := chunk.(type) {
		case html.StaticChunk:
			if _, err := buf.Write(chunk); err != nil {
				return err
			}
		case html.DynamicChunk:
			*r.chunks = append(*r.chunks, buf.Bytes())
			bufpool.Put(buf)
			buf = bufpool.New().(*bytes.Buffer)
			*r.chunks = append(*r.chunks, chunk)
		default:
			return fmt.Errorf("invalid chunk type %T", chunk)
		}
	}
	*r.chunks = append(*r.chunks, buf.Bytes())
	return nil
}

func (r *StaticRenderer) Render(ctx context.Context, w io.Writer) error {
	bw := bufio.NewWriter(w)

	for _, chunk := range *r.chunks {
		switch chunk := chunk.(type) {
		case []byte:
			if _, err := bw.Write(chunk); err != nil {
				return err
			}
		case html.DynamicChunk:
			if err := chunk(ctx, bw); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid chunk type %t", chunk)
		}
	}
	return bw.Flush()
}

func (r *StaticRenderer) Release() {
	chunkpool.Put(r.chunks)
}

func (r *StaticRenderer) Clear() {
	r.Release()
	r.chunks = chunkpool.New().(*[]any)
}

func NewStaticRenderer() *StaticRenderer {
	return &StaticRenderer{
		chunks: chunkpool.New().(*[]any),
	}
}
