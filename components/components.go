package components

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	r "github.com/canpacis/pacis/renderer"
)

/*
	Provide x-data attributes to your elements.

Usage:

	Div(
		D{"open": false}
	) // <div x-data="{'open':false}"></div>
*/
type D map[string]any

func (d D) Render(ctx context.Context, w io.Writer) error {
	enc, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(strings.ReplaceAll(string(enc), "\"", "'")))
	return err
}

func (D) GetKey() string {
	return "x-data"
}

func (d D) GetValue() any {
	return d
}

/*
	Provide arbitrary x attributes to your elements.

Usage:

	Div(X("show", "false")) // <div x-show="false"></div>
*/
func X(key string, value ...any) r.Attribute {
	return r.Attr(fmt.Sprintf("x-%s", key), value...)
}

/*
	Provide arbitrary event attributes to your elements.

Usage:

	Button(On("click", "console.log('clicked')")) // <button x-on:click="console.log('clicked')"></button>
*/
func On(event string, handler string) r.Attribute {
	return r.Attr(fmt.Sprintf("x-on:%s", event), handler)
}

func join(props []r.I, rest ...r.I) []r.I {
	return append(props, rest...)
}

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

// Place this component inside the head tag
func AppHead() r.Node {
	return r.Frag(
		r.Link(r.Href("https://fonts.googleapis.com"), r.Rel("preconnect")),
		r.Link(r.Href("https://fonts.gstatic.com"), r.Rel("preconnect")),
		r.Link(r.Href("https://fonts.googleapis.com/css2?family=Inter:opsz,wght@14..32,100..900&display=swap"), r.Rel("stylesheet")),
		r.Link(r.Href("/public/main.css"), r.Rel("stylesheet")),
		r.Script(r.Src("https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js")),
		r.Script(r.Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js")),
		r.Script(r.Src("https://unpkg.com/embla-carousel/embla-carousel.umd.js")),
	)
}
