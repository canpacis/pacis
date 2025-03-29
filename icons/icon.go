//go:generate go run ./generate/main.go
package icons

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	r "github.com/canpacis/pacis/renderer"
)

/*
An error handler element that you can use with error boundaries

Usage:

	Try(
		MightError(),
		ErrorText, // provides a simple error message on both frontend and the terminal
	)
*/
func ErrorText(err error) r.Node {
	log.Println(err.Error())

	return r.Div(
		r.Class("fixed inset-0 flex justify-center items-center z-80"),

		r.Div(r.Class("bg-neutral-800/60 absolute inset-0")),
		r.Div(
			r.Class("bg-neutral-800 text-red-600 rounded-sm p-4 relative z-90"),

			r.Textf("Failed to render: %s", err.Error()),
		),
	)
}

type Width float64

func (Width) GetKey() string {
	return "width"
}

func (wd Width) GetValue() any {
	return float64(wd)
}

func (Width) Dedupe() {}

func (wd Width) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(strconv.FormatFloat(float64(wd), 'f', -1, 64)))
	return err
}

type Height float64

func (Height) GetKey() string {
	return "height"
}

func (wd Height) GetValue() any {
	return float64(wd)
}

func (Height) Dedupe() {}

func (wd Height) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(strconv.FormatFloat(float64(wd), 'f', -1, 64)))
	return err
}

type StrokeWidth float64

func (StrokeWidth) GetKey() string {
	return "stroke-width"
}

func (wd StrokeWidth) GetValue() any {
	return float64(wd)
}

func (StrokeWidth) Dedupe() {}

func (wd StrokeWidth) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(strconv.FormatFloat(float64(wd), 'f', -1, 64)))
	return err
}

func Fill(fill string) r.Attribute {
	return r.Attr("fill", fill)
}

func Stroke(fill string) r.Attribute {
	return r.Attr("stroke", fill)
}

type SvgIcon struct {
	r.Element
	path string

	Content []byte `xml:",innerxml"`
}

func Icon(path string, items ...r.I) r.Node {
	props := []r.I{
		Width(24),
		Height(24),
		StrokeWidth(2),
		Fill("none"),
		Stroke("currentColor"),
		r.Attr("viewBox", "0 0 24 24"),
		r.Attr("stroke-linecap", "round"),
		r.Attr("stroke-linejoin", "round"),
	}
	props = append(props, items...)

	icon := SvgIcon{path: path, Element: r.El("svg", props...)}

	file, err := os.OpenFile(fmt.Sprintf("./lucide/icons/%s.svg", path), os.O_RDONLY, 0o644)
	if err != nil {
		errset, ok := icon.Element.(r.ErrorSetter)
		if ok {
			errset.SetError(err)
			return ErrorText(err)
		} else {
			panic(err)
		}
	}
	defer file.Close()

	err = xml.NewDecoder(file).Decode(&icon)
	if err != nil {
		errset, ok := icon.Element.(r.ErrorSetter)
		if ok {
			errset.SetError(err)
		} else {
			panic(err)
		}
	}
	icon.AddNode(r.RawUnsafe(icon.Content))

	return r.Try(icon, ErrorText)
}
