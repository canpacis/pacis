package renderer

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strings"
)

var selfClosingTags = map[string]string{
	"area":   "area",
	"base":   "base",
	"br":     "br",
	"col":    "col",
	"embed":  "embed",
	"hr":     "hr",
	"img":    "img",
	"input":  "input",
	"link":   "link",
	"meta":   "meta",
	"param":  "param",
	"source": "source",
	"track":  "track",
	"wbr":    "wbr",
}

type noopwriter struct {
	w   io.Writer
	n   int
	err error
}

func (nw *noopwriter) write(p []byte) {
	if nw.err != nil {
		return
	}

	nw.n, nw.err = nw.w.Write(p)
}

func (nw *noopwriter) reset(w io.Writer) {
	nw.w = w
	nw.n = 0
	nw.err = nil
}

type Renderer interface {
	Render(io.Writer) error
}

type Attribute interface {
	Renderer
}

type HtmlAttribute struct {
	Key   string
	Value any
}

// Attribute renders only its value for key deduplication
func (a *HtmlAttribute) Render(w io.Writer) error {
	var str string

	switch value := a.Value.(type) {
	case bool:
		if value {
			str = "true"
		} else {
			str = "false"
		}
	case string:
		if a.Key == "class" {
			str = strings.ReplaceAll(value, "\"", "'")
		} else {
			str = html.EscapeString(value)
		}
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64, int, uint:
		str = fmt.Sprintf("%d", value)
	case interface{ String() string }:
		str = html.EscapeString(value.String())
	default:
		panic(fmt.Sprintf("unsupported attribute type: %t", value))
	}

	_, err := w.Write([]byte(str))
	return err
}

func (a HtmlAttribute) Attribute() string {
	return a.Key
}

func Attr(key string, value ...any) *HtmlAttribute {
	switch len(value) {
	case 0:
		return &HtmlAttribute{Key: key, Value: ""}
	case 1:
		return &HtmlAttribute{Key: key, Value: value[0]}
	default:
		panic("invalid number of attribute values")
	}
}

type Node interface {
	Renderer
	Node()
}

type Element struct {
	tag         string
	attrs       []Attribute
	children    []Node
	selfClosing bool

	*noopwriter
}

func (e *Element) Render(w io.Writer) error {
	e.reset(w)

	e.write(fmt.Appendf(nil, "<%s", e.tag))

	attrs := map[string][]Attribute{}

	// Extract duplicate attributes to be joined
	for _, attr := range e.attrs {
		htmlAttr, ok := attr.(*HtmlAttribute)
		if ok {
			attrs[htmlAttr.Key] = append(attrs[htmlAttr.Key], attr)
		}
		keyer, ok := attr.(interface{ Key() string })
		if ok {
			attrs[keyer.Key()] = append(attrs[keyer.Key()], attr)
		}
	}

	for key, attrs := range attrs {
		e.write(fmt.Appendf(nil, " %s=\"", key))

		for i, attr := range attrs {
			e.err = attr.Render(w)

			// seperate different attrs declarations with space
			if i < len(attrs)-1 {
				e.write([]byte(" "))
			}
		}

		e.write([]byte("\""))
	}

	if e.selfClosing {
		e.write([]byte(" />"))
		return e.err
	}
	e.write([]byte(">"))

	for _, child := range e.children {
		e.err = child.Render(w)
	}

	e.write(fmt.Appendf(nil, "</%s>", e.tag))

	return e.err
}

func El(tag string, constituents ...Renderer) *Element {
	children := []Node{}
	attrs := []Attribute{}

	for _, constituent := range constituents {
		switch constituent := constituent.(type) {
		case Text:
			children = append(children, &constituent)
		case Node:
			children = append(children, constituent)
		case Attribute:
			attrs = append(attrs, constituent)
		}
	}

	_, selfClosing := selfClosingTags[tag]
	return &Element{
		tag:         tag,
		attrs:       attrs,
		children:    children,
		noopwriter:  &noopwriter{},
		selfClosing: selfClosing,
	}
}

type Fragment struct {
	children []Node
}

func (f *Fragment) Render(w io.Writer) error {
	for _, child := range f.children {
		if err := child.Render(w); err != nil {
			return err
		}
	}
	return nil
}

func Frag(children ...Node) *Fragment {
	return &Fragment{children: children}
}

type Text string

func (t Text) Render(w io.Writer) error {
	_, err := w.Write([]byte(html.EscapeString(string(t))))
	return err
}

type RawNode struct {
	data   string
	escape bool
}

func (n *RawNode) Render(w io.Writer) (err error) {
	if n.escape {
		_, err = w.Write([]byte(html.EscapeString(n.data)))
	} else {
		_, err = w.Write([]byte(n.data))
	}
	return err
}

func Raw(data string) *RawNode {
	return &RawNode{data: data, escape: true}
}

func UnsafeRaw(data string) *RawNode {
	return &RawNode{data: data, escape: false}
}

func (*Element) Node()  {}
func (*Fragment) Node() {}
func (*Text) Node()     {}
func (*RawNode) Node()  {}

func GetAttr(el *Element, name string) (string, bool) {

	for _, attr := range el.attrs {
		var key string

		htmlAttr, ok := attr.(*HtmlAttribute)
		if ok {
			key = htmlAttr.Key
		} else {
			keyer, ok := attr.(interface{ Key() string })
			if ok {
				key = keyer.Key()
			}
		}

		if len(key) == 0 {
			return "", false
		}

		if key == name {
			buf := bytes.NewBuffer([]byte{})
			attr.Render(buf)
			return buf.String(), true
		}
	}

	return "", false
}
