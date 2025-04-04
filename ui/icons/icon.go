//go:generate go run ./generate/main.go
package icons

import (
	"context"
	"io"
	"strconv"

	h "github.com/canpacis/pacis/ui/html"
)

func join(props []h.I, rest ...h.I) []h.I {
	return append(rest, props...)
}

type Width float64

func (Width) GetKey() string {
	return "width"
}

func (wd Width) IsEmpty() bool {
	return false
}

// Implements Deduper interface to deduplicate attribute
// and use the last provided value as the final attribte
func (Width) Dedupe() {}

func (wd Width) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(strconv.FormatFloat(float64(wd), 'f', -1, 64)))
	return err
}

type Height float64

func (Height) GetKey() string {
	return "height"
}

func (wd Height) IsEmpty() bool {
	return false
}

// Implements Deduper interface to deduplicate attribute
// and use the last provided value as the final attribte
func (Height) Dedupe() {}

func (wd Height) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(strconv.FormatFloat(float64(wd), 'f', -1, 64)))
	return err
}

type StrokeWidth float64

func (StrokeWidth) GetKey() string {
	return "stroke-width"
}

func (wd StrokeWidth) IsEmpty() bool {
	return false
}

// Implements Deduper interface to deduplicate attribute
// and use the last provided value as the final attribte
func (StrokeWidth) Dedupe() {}

func (wd StrokeWidth) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(strconv.FormatFloat(float64(wd), 'f', -1, 64)))
	return err
}

func Fill(fill string) h.Attribute {
	return h.Attr("fill", fill)
}

func Stroke(fill string) h.Attribute {
	return h.Attr("stroke", fill)
}

type SvgIcon struct {
	h.Element
}

func Icon(props ...h.I) SvgIcon {
	props = join(props,
		Width(24),
		Height(24),
		StrokeWidth(2),
		Fill("none"),
		Stroke("currentColor"),
		h.Attr("viewBox", "0 0 24 24"),
		h.Attr("stroke-linecap", "round"),
		h.Attr("stroke-linejoin", "round"),
	)
	return SvgIcon{Element: h.El("svg", props...)}
}
