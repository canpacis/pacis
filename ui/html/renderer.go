package html

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"iter"
	"log"
	"slices"
	"strings"
	"sync"
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

	// GetAttributes() []*Attribute_
	GetAttribute(string) (*Attribute, bool)
	AddAttribute(*Attribute)
	RemoveAttribute(string)

	GetNodes() []Node
	GetNode(int) (Node, bool)
	AddNode(Node)
	RemoveNode(int)

	GetElement(int) (Element, bool)
	GetElements() []Element
}

type DedupeStrategy int

const (
	DedupeTakeLast = DedupeStrategy(iota)
	DedupeTakeFirst
	DedupeJoinComma
	DedupeJoinSpace
)

// Represents any kind of element attribute
type Attribute struct {
	strategy DedupeStrategy
	Name     string
	value    []any
	IsEmpty  bool

	rendered []byte
}

func (a Attribute) get(i int) string {
	raw := a.value[i]
	switch raw := raw.(type) {
	case string:
		return html.EscapeString(raw)
	case bool:
		if raw {
			return "true"
		}
		return "false"
	case int, uint:
		return fmt.Sprintf("%d", raw)
	default:
		return ""
	}
}

func (a Attribute) strings() []string {
	strs := []string{}
	for i := range a.value {
		strs = append(strs, a.get(i))
	}
	return strs
}

func (a Attribute) Value() string {
	switch a.strategy {
	case DedupeTakeFirst:
		return a.get(0)
	case DedupeTakeLast:
		if len(a.value) == 0 {
			return ""
		}
		return a.get(len(a.value) - 1)
	case DedupeJoinComma:
		return strings.Join(a.strings(), ",")
	case DedupeJoinSpace:
		return strings.Join(a.strings(), " ")
	default:
		return ""
	}
}

func (a Attribute) Raw() any {
	switch a.strategy {
	case DedupeTakeFirst:
		return a.value[0]
	case DedupeTakeLast:
		if len(a.value) == 0 {
			return ""
		}
		return a.value[len(a.value)-1]
	case DedupeJoinComma:
		return strings.Join(a.strings(), ",")
	case DedupeJoinSpace:
		return strings.Join(a.strings(), " ")
	default:
		return ""
	}
}

func (a *Attribute) Render(ctx context.Context, w io.Writer) error {
	if len(a.rendered) != 0 {
		_, err := w.Write(a.rendered)
		return err
	}
	var err error
	var buf = new(bytes.Buffer)
	_, err = buf.Write([]byte(" " + a.Name))
	if err != nil {
		return err
	}

	if !a.IsEmpty {
		_, err = buf.WriteString("=\"" + a.Value() + "\"")
		if err != nil {
			return err
		}
	}
	a.rendered = buf.Bytes()
	_, err = io.Copy(w, buf)
	return err
}

type element struct {
	tag         string
	nodes       []Node
	attrmap     map[string]*Attribute
	selfClosing bool

	err error
}

var bufpool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

func (e *element) Render(ctx context.Context, w io.Writer) error {
	var nocopy bool
	var buf *bytes.Buffer
	b, ok := w.(*bytes.Buffer)
	if ok {
		buf = b
		nocopy = true
	} else {
		buf = bufpool.Get().(*bytes.Buffer)
		buf.Reset()
	}

	var err error
	_, err = buf.WriteString("<" + e.tag)
	if err != nil {
		return err
	}

	for _, attr := range e.attrmap {
		if err := attr.Render(ctx, buf); err != nil {
			return err
		}
	}

	if e.selfClosing {
		_, err := buf.WriteString(" />")
		if err != nil {
			return err
		}
		if !nocopy {
			_, err = io.Copy(w, buf)
			bufpool.Put(buf)
			return err
		}
		return nil
	}
	_, err = buf.WriteString(">")
	if err != nil {
		return err
	}

	for _, node := range e.nodes {
		if err := node.Render(ctx, buf); err != nil {
			return err
		}
	}

	_, err = buf.WriteString("</" + e.tag + ">")
	if err != nil {
		return err
	}
	if !nocopy {
		_, err = io.Copy(w, buf)
		bufpool.Put(buf)
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

func (e *element) GetAttribute(key string) (*Attribute, bool) {
	attr, ok := e.attrmap[key]
	return attr, ok
}

var nullbuf = bytes.NewBuffer(make([]byte, 0, 1024))

func (e *element) AddAttribute(attr *Attribute) {
	attr.Render(context.Background(), nullbuf)
	nullbuf.Reset()
	existing, ok := e.attrmap[attr.Name]
	if ok {
		existing.value = append(existing.value, attr.value...)
	} else {
		e.attrmap[attr.Name] = attr
	}
}

func (e *element) RemoveAttribute(key string) {
	delete(e.attrmap, key)
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
	el := &element{tag: tag, selfClosing: selfClosing, attrmap: make(map[string]*Attribute, 4)}

	for _, item := range items {
		switch item := item.(type) {
		case *Attribute:
			el.AddAttribute(item)
		case Node:
			el.nodes = append(el.nodes, item)
		default:
			unwrapper, ok := item.(interface{ Unwrap() *Attribute })
			if ok {
				el.AddAttribute(unwrapper.Unwrap())
			} else {
				panic(fmt.Sprintf("unknown item type %T", item))
			}
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

func Attr(name string, value ...any) *Attribute {
	switch len(value) {
	case 0:
		return &Attribute{Name: name, IsEmpty: true}
	case 1:
		return &Attribute{Name: name, value: value}
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

func Iter[T any, K int](items any, fn func(T, K) I) I {
	mapped := []I{}

	seq, ok := items.(iter.Seq[T])
	if ok {
		i := 0
		for item := range seq {
			mapped = append(mapped, fn(item, K(i)))
			i++
		}
	} else {
		seq2, ok := items.(iter.Seq2[T, K])
		if ok {
			for item, k := range seq2 {
				mapped = append(mapped, fn(item, k))
			}
		} else {
			slc, ok := items.([]T)
			if ok {
				for i, item := range slc {
					mapped = append(mapped, fn(item, K(i)))
				}
			} else {
				mp, ok := items.(map[K]T)
				if ok {
					keys := []K{}
					for k := range mp {
						keys = append(keys, k)
					}
					slices.Sort(keys)
					for _, key := range keys {
						mapped = append(mapped, fn(mp[key], key))
					}
				} else {
					log.Fatalf("failed to iterate over type %T", items)
				}
			}
		}
	}

	return Frag(mapped...)
}

type case_[T comparable] struct {
	expr T
	node I
}

func Case[T comparable](expr T, node I) case_[T] {
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

func Class(class string) *Attribute {
	attr := Attr("class", class)
	attr.strategy = DedupeJoinSpace
	return attr
}
