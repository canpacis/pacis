package components

import (
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"strings"

	h "github.com/canpacis/pacis/ui/html"
)

func randid() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return "pacis-" + hex.EncodeToString(buf)
}

func readattr(attr h.Attribute) string {
	var buf bytes.Buffer
	attr.Render(context.Background(), &buf)
	return buf.String()
}

func fn(name string, args ...any) string {
	list := []string{}
	for _, arg := range args {
		if arg == nil {
			list = append(list, "null")
			continue
		}
		switch arg := arg.(type) {
		case string:
			if strings.Contains(arg, "\n") {
				list = append(list, "`"+arg+"`")
			} else {
				list = append(list, "'"+arg+"'")
			}
		case int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			list = append(list, fmt.Sprintf("%d", arg))
		case float32, float64:
			list = append(list, fmt.Sprintf("%f", arg))
		case bool:
			list = append(list, fmt.Sprintf("%t", arg))
		default:
			panic(fmt.Sprintf("cannot serialize function argument type %T", arg))
		}
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(list, ", "))
}

type GroupedClass struct {
	group     string
	class     h.Class
	isdefault bool
}

func (GroupedClass) Render(context.Context, io.Writer) error {
	return nil
}

type groupedclasses []GroupedClass

func (list groupedclasses) Render(ctx context.Context, w io.Writer) error {
	if len(list) == 0 {
		return nil
	}

	var def GroupedClass
	var selected GroupedClass
	for _, item := range list {
		if item.isdefault {
			def = item
		} else {
			selected = item
		}
	}
	if selected.group == "" {
		selected = def
	}
	return selected.class.Render(ctx, w)
}

func (groupedclasses) GetKey() string {
	return "class"
}

func (a groupedclasses) IsEmpty() bool {
	return false
}

func (groupedclasses) Dedupe() {}

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
func Join(props []h.I, rest ...h.I) []h.I {
	source := []h.I{}
	source = append(source, props...)
	source = append(source, rest...)

	result := []h.I{}

	groups := map[string]groupedclasses{}

	for _, prop := range source {
		grouped, ok := prop.(*GroupedClass)
		if ok {
			groups[grouped.group] = append(groups[grouped.group], *grouped)
		} else {
			result = append(result, prop)
		}
	}
	for _, group := range groups {
		result = append(result, &group)
	}

	return result
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
func X(key string, value ...any) h.Attribute {
	return h.Attr(fmt.Sprintf("x-%s", key), value...)
}

/*
	Provide arbitrary event attributes to your elements.

Usage:

	Button(On("click", "console.log('clicked')")) // <button x-on:click="console.log('clicked')"></button>
*/
func On(event string, handler string) h.Attribute {
	return h.Attr(fmt.Sprintf("x-on:%s", event), handler)
}

// Toggles color scheme upon a click event
var ToggleColorScheme = On("click", "$store.colorScheme.toggle()")

/*
An error handler element that you can use with error boundaries

Usage:

	Try(
		MightError(),
		ErrorText, // provides a simple error message on both frontend and the terminal
	)
*/
func ErrorText(err error) h.Node {
	log.Println(err.Error())

	return h.Div(
		h.Class("fixed inset-0 flex justify-center items-center z-80"),

		h.Div(h.Class("bg-neutral-800/60 absolute inset-0")),
		h.Div(
			h.Class("bg-neutral-800 text-red-600 rounded-sm p-4 relative z-90"),

			h.Textf("Failed to render: %s", err.Error()),
		),
	)
}

type AppHead struct {
	*h.Fragment
	prefix string
}

//go:embed dist
var dist embed.FS

// Provides the file system to serve statically
func (h AppHead) FS() fs.FS {
	content, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return content
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
	head := &AppHead{prefix: prefix, Fragment: h.Frag(
		h.Meta(h.Charset("UTF-8")),
		h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1.0")),
		h.Link(h.Href(fmt.Sprintf("%smain.css", prefix)), h.Rel("stylesheet")),
		h.Script(h.Src(fmt.Sprintf("%smain.js", prefix))),
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
	element func(items ...h.I) h.Element
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

func Replace(element func(items ...h.I) h.Element) *Replacer {
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

type ComponentAttribute int

func (ComponentAttribute) Render(context.Context, io.Writer) error {
	return nil
}

func (a ComponentAttribute) GetKey() string {
	switch a {
	case Clearable:
		return "clearable"
	case Open:
		return "open"
	default:
		return "invalid-input-attribute"
	}
}

func (ComponentAttribute) IsEmpty() bool {
	return true
}

const (
	Clearable = ComponentAttribute(iota)
	Open
)

func Init(handler string) h.Attribute {
	return On("init", handler)
}

func Changed(handler string) h.Attribute {
	return On("changed", handler)
}

func Opened(handler string) h.Attribute {
	return On("opened", handler)
}

func Closed(handler string) h.Attribute {
	return On("closed", handler)
}

func Dismissed(handler string) h.Attribute {
	return On("dismissed", handler)
}
