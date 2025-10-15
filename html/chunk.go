package html

import (
	"context"
	"fmt"
	"io"
)

type Chunk interface {
	chunk()
}

type StaticChunk []byte

type DynamicChunk func(context.Context, io.Writer) error

func (StaticChunk) chunk()  {}
func (DynamicChunk) chunk() {}

func Render(chunk Chunk, ctx context.Context, w io.Writer) error {
	switch chunk := chunk.(type) {
	case DynamicChunk:
		return chunk(ctx, w)
	case StaticChunk:
		_, err := w.Write(chunk)
		return err
	default:
		return fmt.Errorf("invalid chunk type %t", chunk)
	}
}

type ChunkWriter interface {
	Write(...Chunk)
	Chunks() []Chunk
}

type cw struct {
	buf []Chunk
}

func (cw *cw) Write(chunks ...Chunk) {
	cw.buf = append(cw.buf, chunks...)
}

func (cw cw) Chunks() []Chunk {
	return cw.buf
}

func NewChunkWriter() ChunkWriter {
	return &cw{buf: []Chunk{}}
}

type teecw struct {
	cw ChunkWriter
	fn func(Chunk) error
}

func (cw *teecw) Write(chunks ...Chunk) {
	for _, chunk := range chunks {
		cw.fn(chunk)
	}
	cw.cw.Write(chunks...)
}

func (cw *teecw) Chunks() []Chunk {
	return cw.cw.Chunks()
}
