package html_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/canpacis/pacis/ui/html"
)

func TestRendering(t *testing.T) {
	el := html.Div(
		html.Class("flex justify-center items-center"),
		html.Class("flex-col"),

		html.P(html.Text("Hello, World")),
		html.Map([]string{"Lorem", "ipsum", "dolor", "sit", "amet"}, func(item string, i int) html.I {
			return html.Span(html.Text(item))
		}),
	)

	el.Render(context.Background(), os.Stdout)
}

func BenchmarkRendering(b *testing.B) {
	el := html.Div(
		html.Class("flex justify-center items-center"),

		html.P(html.Text("Hello, World")),
		html.Map([]string{"Lorem", "ipsum", "dolor", "sit", "amet"}, func(item string, i int) html.I {
			return html.Span(html.Text(item))
		}),
	)

	for b.Loop() {
		var buf = new(bytes.Buffer)
		el.Render(context.Background(), buf)
	}
}
