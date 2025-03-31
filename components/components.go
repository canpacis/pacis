package components

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"

	r "github.com/canpacis/pacis-ui/renderer"
)

/*
	Joins a prop list with rest. Puts the props at the end for correct attribute deduplication.

Usage:

	func Component(props ...I) Element {
		return Div(
			Join(
				props,
				Class( ... )
			)...
		)
	}
*/
func Join(props []r.I, rest ...r.I) []r.I {
	return append(rest, props...)
}

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

func (d D) IsEmpty() bool {
	return d == nil
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

type AppHead struct {
	*r.Fragment
	prefix string
}

//go:embed public
var public embed.FS

// Provides an http handler for public assets
func (h AppHead) Handler() http.Handler {
	content, err := fs.Sub(public, "public")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(content))
}

/*
	Place this component inside the head tag and use the handler to add static files to your server

Usage:

	head := CreateHead("/public/") // include the '/' before and after
	http.Handle("/public/", head.Handler())

	html := Html(
		Head(head)
		Body( ... )
	)
*/
func CreateHead(prefix string) *AppHead {
	head := &AppHead{prefix: prefix, Fragment: r.Frag(
		r.Meta(r.Charset("UTF-8")),
		r.Meta(r.Name("viewport"), r.Content("width=device-width, initial-scale=1.0")),
		r.Link(r.Href(fmt.Sprintf("%smain.css", prefix)), r.Rel("stylesheet")),
		r.Script(r.Src(fmt.Sprintf("%smain.js", prefix))),
		r.Script(r.Src(fmt.Sprintf("%scarousel.js", prefix))),
		r.Script(r.Src(fmt.Sprintf("%sfocus.js", prefix))),
		r.Script(r.Src(fmt.Sprintf("%sanchor.js", prefix))),
		r.Script(r.Src(fmt.Sprintf("%salpine.js", prefix))),
	)}

	return head
}

/*
Provides vertical positioning for anchored elements like tooltips and dropdown menus.
See components.AnchorPosition type or components.Anchor function for usage details.

Available options are;
  - components.VTop: Position on top of an element
  - components.VBottom: Position on the bottom of an element
*/
type VPos int

const (
	VTop = VPos(iota)
	VBottom
)

/*
Provides horizontal positioning for anchored elements like tooltips and dropdown menus.
See components.AnchorPosition type or components.Anchor function for usage details.

Available options are;
  - components.HStart: Position at the start of an element
  - components.HCenter: Position at the center of an element
  - components.HEnd: Position at the end of an element
*/
type HPos int

const (
	HStart = HPos(iota)
	HCenter
	HEnd
)

/*
	Provides anchor positioning attributes to given element

Usage:

	Dropdown(
		DropdownTrigger( ... )
		DropdownContent(
			// Pass this to content elements
			Anchor(VBottom, HCenter, 12)
			// Positions content at bottom center of the trigger, offsetted 12 pixels
		)
	)
*/
type AnchorPosition struct {
	vpos   VPos
	hpos   HPos
	offset int
}

func (a AnchorPosition) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte("$refs.anchor"))
	return err
}

func (a AnchorPosition) GetKey() string {
	key := "x-anchor"

	switch a.vpos {
	case VTop:
		key += ".top"
	case VBottom:
		key += ".bottom"
	default:
		panic("invalid vertical position for anchor")
	}

	switch a.hpos {
	case HStart:
		key += "-start"
	case HCenter:
		key += "-center"
	case HEnd:
		key += "-end"
	default:
		panic("invalid horizontal position for anchor")
	}

	key += fmt.Sprintf(".offset.%d", a.offset)

	return key
}

func (a AnchorPosition) IsEmpty() bool {
	return false
}

/*
	Provides anchor positioning attributes to given element

Usage:

	Dropdown(
		DropdownTrigger( ... )
		DropdownContent(
			// Pass this to content elements
			Anchor(VBottom, HCenter, 12)
			// Positions content at bottom center of the trigger, offsetted 12 pixels
		)
	)
*/
func Anchor(v VPos, h HPos, offset int) AnchorPosition {
	return AnchorPosition{vpos: v, hpos: h, offset: offset}
}

// Implements Deduper interface to deduplicate attribute
// and use the last provided value as the final attribte
func (a AnchorPosition) Dedupe() {}

type Replacer struct {
	element func(items ...r.I) r.Element
}

func (*Replacer) Render(context.Context, io.Writer) error {
	return nil
}

func (*Replacer) GetKey() string {
	return "replace"
}

func (*Replacer) IsEmpty() bool {
	return true
}

func Replace(element func(items ...r.I) r.Element) *Replacer {
	return &Replacer{element: element}
}

type Orientation int

const (
	OHorizontal = Orientation(iota)
	OVertical
)

func (o Orientation) String() string {
	switch o {
	case OHorizontal:
		return "horizontal"
	case OVertical:
		return "vertical"
	default:
		panic("invalid orientation value")
	}
}
