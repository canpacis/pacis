package html

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"slices"
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

// Represents a general renderable structure eg: an element, an attribute
type Renderer interface {
	Render(context.Context, io.Writer) error
}

// A shorter alias to Renderer
type I = Renderer

// Node type enumerator
type NodeType int

const (
	NodeText = NodeType(iota)
	NodeElement
	NodeFragment
)

// Represents an HTML node that is renderable. This
// can be anything from an element to a text node.
type Node interface {
	Renderer
	NodeType() NodeType
}

// Represents and HTML element, conceptually encompassing a
// Node in that it also renders a node but can also have attributes.
type Element interface {
	Node
	GetTag() string

	GetAttributes() []Attribute
	GetAttribute(string) (Attribute, bool)
	AddAttribute(Attribute)
	RemoveAttribute(string)

	GetNodes() []Node
	GetNode(int) (Node, bool)
	AddNode(Node)
	RemoveNode(int)

	GetElement(int) (Element, bool)
	GetElements() []Element
}

// ErrorSetter is an interface that allows a node to
// set its error so that it can be raised during rendering.
type ErrorSetter interface {
	SetError(error)
}

// Represents any kind of element attribute
type Attribute interface {
	Renderer
	GetKey() string
	IsEmpty() bool
}

type element struct {
	tag         string
	nodes       []Node
	attrs       []Attribute
	selfClosing bool

	err error
}

func (e *element) Render(ctx context.Context, w io.Writer) error {
	if e.GetError() != nil {
		return e.GetError()
	}

	if _, err := w.Write(fmt.Appendf(nil, "<%s", e.tag)); err != nil {
		return err
	}

	attrs := map[string][]Attribute{}

	// Collapse duplicate definitions
	for _, attr := range e.attrs {
		key := attr.GetKey()
		attrs[key] = append(attrs[key], attr)
	}

	// TODO: Maybe refactor this to its own function
	for key, list := range attrs {
		// join class names with a space, duplicate other attributes
		if key == "class" {
			if _, err := w.Write(fmt.Appendf(nil, " %s=\"", key)); err != nil {
				return err
			}

			for i, class := range list {
				if i != 0 {
					if _, err := w.Write([]byte(" ")); err != nil {
						return err
					}
				}
				if err := class.Render(ctx, w); err != nil {
					return err
				}
			}
			if _, err := w.Write([]byte("\"")); err != nil {
				return err
			}
		} else {
			if len(list) == 0 {
				continue
			}
			// if the attribute is dedupable, pick the last element
			_, ok := list[0].(interface{ Dedupe() })
			if ok {
				attr := list[len(list)-1]
				if attr.IsEmpty() {
					if _, err := w.Write(fmt.Appendf(nil, " %s", key)); err != nil {
						return err
					}
				} else {
					if _, err := w.Write(fmt.Appendf(nil, " %s=\"", key)); err != nil {
						return err
					}
					if err := attr.Render(ctx, w); err != nil {
						return err
					}
					if _, err := w.Write([]byte("\"")); err != nil {
						return err
					}
				}
			} else {
				for _, attr := range list {
					if attr.IsEmpty() {
						if _, err := w.Write(fmt.Appendf(nil, " %s", key)); err != nil {
							return err
						}
					} else {
						if _, err := w.Write(fmt.Appendf(nil, " %s=\"", key)); err != nil {
							return err
						}
						if err := attr.Render(ctx, w); err != nil {
							return err
						}
						if _, err := w.Write([]byte("\"")); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	if e.selfClosing {
		_, err := w.Write([]byte(" />"))
		return err
	}
	if _, err := w.Write([]byte(">")); err != nil {
		return err
	}

	for _, node := range e.nodes {
		if err := node.Render(ctx, w); err != nil {
			return err
		}
	}

	if _, err := w.Write(fmt.Appendf(nil, "</%s>", e.tag)); err != nil {
		return err
	}
	return nil
}

func (e *element) NodeType() NodeType {
	return NodeElement
}

func (e *element) GetTag() string {
	return e.tag
}

func (e *element) GetAttributes() []Attribute {
	return e.attrs
}

func (e *element) GetAttribute(key string) (Attribute, bool) {
	for _, attr := range e.attrs {
		if attr.GetKey() == key {
			return attr, true
		}
	}
	return nil, false
}

func (e *element) AddAttribute(attr Attribute) {
	e.attrs = append(e.attrs, attr)
}

func (e *element) RemoveAttribute(key string) {
	idx := -1
	for i, attr := range e.attrs {
		if attr.GetKey() == key {
			idx = i
		}
	}
	if idx < 0 {
		e.err = fmt.Errorf("remove attribute: cannot remove attribute %s on element %s, attribute does not exist", key, e.tag)
	} else {
		e.attrs = slices.Delete(e.attrs, idx, idx+1)
	}
}

func (e *element) GetElement(i int) (Element, bool) {
	elements := e.GetElements()
	if len(elements) <= i {
		return nil, false
	}
	return elements[i], false
}

func (e *element) GetElements() []Element {
	elements := []Element{}

	for _, node := range e.nodes {
		element, ok := node.(Element)
		if ok {
			elements = append(elements, element)
		}
	}
	return elements
}

func (e *element) GetNode(i int) (Node, bool) {
	if len(e.nodes) <= i {
		return nil, false
	}
	return e.nodes[i], false
}

func (e *element) GetNodes() []Node {
	return e.nodes
}

func (e *element) AddNode(node Node) {
	e.nodes = append(e.nodes, node)
}

func (e *element) RemoveNode(i int) {
	if len(e.nodes) <= i {
		e.err = fmt.Errorf("remove: cannot remove node %d on element %s, index is out of bounds", i, e.tag)
		return
	}
	e.nodes = slices.Delete(e.nodes, i, i+1)
}

func (e *element) SetError(err error) {
	e.err = err
}

func (e *element) GetError() error {
	return e.err
}

// Creates an element with default html renderer
func El(tag string, items ...I) Element {
	_, selfClosing := selfClosingTags[tag]
	el := &element{tag: tag, selfClosing: selfClosing}

	for _, item := range items {
		switch item := item.(type) {
		case Attribute:
			el.attrs = append(el.attrs, item)
		case Node:
			el.nodes = append(el.nodes, item)
		default:
			fmt.Println(tag, item)
			panic(fmt.Sprintf("unknown item type %T", item))
		}
	}

	return el
}

// Clones an element. This does not clone the Render() method
// of that element, instead it creates a new element with the
// default html element renderer.
func Clone(elem Element, items ...I) Element {
	_, selfClosing := selfClosingTags[elem.GetTag()]
	el := &element{
		tag:         elem.GetTag(),
		nodes:       elem.GetNodes(),
		attrs:       elem.GetAttributes(),
		selfClosing: selfClosing,
	}

	for _, item := range items {
		switch item := item.(type) {
		case Attribute:
			el.attrs = append(el.attrs, item)
		case Node:
			el.nodes = append(el.nodes, item)
		default:
			panic(fmt.Sprintf("unknown item type %T", item))
		}
	}

	return el
}

// Represents a text node
type Text string

func (t Text) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(html.EscapeString(string(t))))
	return err
}

func (Text) NodeType() NodeType {
	return NodeText
}

// Create a text node with formatting
func Textf(format string, a ...any) Text {
	return Text(fmt.Sprintf(format, a...))
}

// Represents an unsafe raw text
// content, use it at your own risk.
type RawUnsafe string

func (t RawUnsafe) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(string(t)))
	return err
}

func (RawUnsafe) NodeType() NodeType {
	return NodeText
}

type attr struct {
	key   string
	value any
}

func (a attr) Render(ctx context.Context, w io.Writer) error {
	var str string

	switch value := a.value.(type) {
	case bool:
		if value {
			str = "true"
		} else {
			str = "false"
		}
	case string:
		if a.key == "class" {
			str = strings.ReplaceAll(value, "\"", "'")
		} else {
			str = html.EscapeString(value)
		}
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64, int, uint:
		str = fmt.Sprintf("%d", value)
	case interface{ String() string }:
		str = html.EscapeString(value.String())
	default:
		str = fmt.Sprintf("%v", value)
	}

	_, err := w.Write([]byte(str))
	return err
}

func (a *attr) GetKey() string {
	return a.key
}

func (a *attr) IsEmpty() bool {
	return a.value == nil
}

func (*attr) Dedupe() {}

func Attr(key string, value ...any) Attribute {
	switch len(value) {
	case 0:
		return &attr{key: key, value: nil}
	case 1:
		return &attr{key: key, value: value[0]}
	default:
		panic("attr expects no more than 2 parameters")
	}
}

type Fragment struct {
	children []I
}

func (f *Fragment) Render(ctx context.Context, w io.Writer) error {
	for _, child := range f.children {
		if err := child.Render(ctx, w); err != nil {
			return err
		}
	}
	return nil
}

func (Fragment) NodeType() NodeType {
	return NodeFragment
}

func Frag(children ...I) *Fragment {
	return &Fragment{children: children}
}

func Map[T any](items []T, fn func(T, int) I) I {
	mapped := []I{}

	for i, item := range items {
		mapped = append(mapped, fn(item, i))
	}

	return Frag(mapped...)
}

type case_[T comparable] struct {
	expr T
	node Node
}

func Case[T comparable](expr T, node Node) case_[T] {
	return case_[T]{expr, node}
}

func SwitchCase[T comparable](expr T, cases ...case_[T]) I {
	for _, c := range cases {
		if expr == c.expr {
			return c.node
		}
	}
	return Frag()
}

func If(cond bool, elem Renderer) Renderer {
	if cond {
		return elem
	}
	return Frag()
}

func IfFn(cond bool, fn func() Renderer) Renderer {
	if cond {
		return fn()
	}
	return Frag()
}

type Boundary struct {
	node     Node
	fallback func(error) Node
}

func (b *Boundary) Render(ctx context.Context, w io.Writer) error {
	buf := bytes.NewBuffer([]byte{})
	err := b.node.Render(ctx, w)
	if err == nil {
		_, err := io.Copy(w, buf)
		return err
	}
	return b.fallback(err).Render(ctx, w)
}

func (b Boundary) NodeType() NodeType {
	return NodeFragment
}

func Try(node Node, fallback func(error) Node) *Boundary {
	return &Boundary{node: node, fallback: fallback}
}

type ContextNode func(context.Context) Node

func (cn ContextNode) Render(ctx context.Context, w io.Writer) error {
	return cn(ctx).Render(ctx, w)
}

func (ContextNode) NodeType() NodeType {
	return NodeFragment
}

type ContextAttr func(context.Context) Attribute

func (ca ContextAttr) Render(ctx context.Context, w io.Writer) error {
	return ca(ctx).Render(ctx, w)
}

func (ca ContextAttr) GetKey() string {
	return ca(context.Background()).GetKey()
}

func (ContextAttr) IsEmpty() bool {
	return false
}

type Class string

func (ca Class) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(ca))
	return err
}

func (ca Class) GetKey() string {
	return "class"
}

func (Class) IsEmpty() bool {
	return false
}

func (Class) Dedupe() {}
