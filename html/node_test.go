package html_test

import (
	"bytes"
	"context"
	"slices"
	"testing"

	"github.com/canpacis/pacis/html"
	"github.com/stretchr/testify/assert"
)

type ChunkTest struct {
	Node     html.Node
	N        int
	Impure   bool
	Rendered string
}

func (c ChunkTest) Assert(a *assert.Assertions) {
	chunks := slices.Collect(c.Node.Chunks())

	if c.N > 0 {
		a.Equal(c.N, len(chunks))
	}

	buf := new(bytes.Buffer)

	for _, chunk := range chunks {
		if err := html.Render(chunk, context.Background(), buf); err != nil {
			a.Fail(err.Error())
		}
	}

	a.Equal(c.Rendered, buf.String())
}

func TestNodeChunks(t *testing.T) {
	tests := []ChunkTest{
		{
			Node:     html.Text("Hello, World!"),
			Rendered: "Hello, World!",
			N:        1,
		},
		{
			Node:     html.Div(),
			Rendered: "<div></div>",
		},
		{
			Node:     html.Doctype,
			Rendered: "<!DOCTYPE html>",
		},
		{
			Node:     html.Input(html.P()),
			Rendered: "<input>",
		},
		{
			Node:     html.Div(html.ID("app"), html.Data("app", "pacis")),
			Rendered: `<div id="app" data-app="pacis"></div>`,
		},
		{
			Node:     html.Div(html.P(html.Text("Hello, World!"))),
			Rendered: "<div><p>Hello, World!</p></div>",
		},
		{
			Node:     html.Fragment(html.Head(), html.Body()),
			Rendered: "<head></head><body></body>",
		},
		{
			Node: html.Component(func(ctx context.Context) html.Node {
				return html.Div(html.ID("app"))
			}),
			Rendered: `<div id="app"></div>`,
			Impure:   true,
		},
		{
			Node:     html.Body(html.Deferred(func(ctx context.Context) html.Property { return html.Class("dark") })),
			Rendered: `<body class="dark"></body>`,
			Impure:   true,
		},
		{
			Node:     html.Body(html.ID("app"), html.Deferred(func(ctx context.Context) html.Property { return html.Class("dark") })),
			Rendered: `<body id="app" class="dark"></body>`,
			Impure:   true,
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		test.Assert(assert)
	}
}
