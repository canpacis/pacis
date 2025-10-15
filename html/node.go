// Package html provides a lightweight HTML rendering system for Go, allowing the construction
// and rendering of HTML nodes and elements in a composable and type-safe manner.
//
// The package defines the Node interface, which represents any renderable HTML node, and provides
// implementations for text nodes, fragments (groups of nodes), and elements (HTML tags with attributes
// and children). Elements can be created with properties such as attributes and classes, and support
// both standard and void (self-closing) HTML elements.
//
// Key types and functions:
//   - Node: Interface for renderable HTML nodes.
//   - Text: Represents plain text content, automatically HTML-escaped.
//   - Frag: Represents a group of child nodes rendered in sequence.
//   - Element: Represents an HTML element with tag name, attributes, and children.
//   - Property: Interface for properties that can be applied to elements (e.g., Attribute).
//   - Attribute: Represents an HTML attribute key-value pair.
//   - El: Constructs a new Element with the given tag name, children, and properties.
//   - VoidEl: Constructs a new void (self-closing) Element.
//   - Fragment: Helper to create a Frag from a variadic list of nodes.
//
// Example usage:
//
//	div := html.El("div",
//	    html.Attr("id", "main"),
//	    html.Text("Hello, world!"),
//	    html.El("span", html.Text("Nested span")),
//	)
//
//	err := div.Render(context.Background(), os.Stdout)
package html

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"iter"
	"log"
	"slices"
	"strings"
	"sync"
)

// Item is an alias for any type.
type Item interface {
	Item()
}

// Node represents an element that can be rendered to an io.Writer within a given context.
// Implementations of Node should define the Render method to output their content.
type Node interface {
	Item
	Render(ChunkWriter)
	Release()
}

// Text represents a node containing plain text content within the HTML renderer.
type Text string

// Implements the Item interface.
func (Text) Item() {}

// Implements the Node interface.
func (t Text) Render(w ChunkWriter) {
	w.Write(StaticChunk(html.EscapeString(string(t))))
}

// Implements the Node interface.
func (Text) Release() {
	// no-op
}

// Textf formats according to a format specifier and returns the resulting Text.
func Textf(format string, a ...any) Text {
	return Text(fmt.Sprintf(format, a...))
}

// Frag represents a fragment of nodes, allowing multiple Node elements to be grouped together.
type Frag []Node

// Implements the Item interface.
func (Frag) Item() {}

// Implements the Node interface.
func (f Frag) Release() {
	for _, node := range f {
		node.Release()
	}
}

// Implements the Node interface.
func (f Frag) Render(w ChunkWriter) {
	for _, node := range f {
		node.Render(w)
	}
}

// Fragment creates a Frag from the provided nodes, allowing multiple nodes to be grouped together.
// It accepts a variadic number of Node arguments and returns a Frag containing them.
func Fragment(nodes ...Node) Frag {
	return Frag(nodes)
}

type PropertyLifeCycle int

const (
	LifeCycleImmediate = PropertyLifeCycle(iota)
	LifeCycleStatic
	LifeCycleDeferred
)

// Property represents an interface for applying a property to an Element.
// Implementations of Property should define how the property modifies or affects the given Element.
type Property interface {
	Item()
	LifeCycle() PropertyLifeCycle
}

// Attribute represents a key-value pair used as an attribute in an HTML node.
type Attribute struct {
	Key   string
	Value string
}

// Implements the Item interface.
func (*Attribute) Item() {}

func (*Attribute) LifeCycle() PropertyLifeCycle {
	return LifeCycleImmediate
}

// Implements the Propterty interface.
func (a *Attribute) Apply(el *Element) {
	if a.Key == "class" {
		el.AddClass(a.Value)
	} else {
		el.SetAttribute(a.Key, a.Value)
	}
}

type DeferredAttribute struct {
	key string
	fn  func(context.Context) string
}

// Implements the Item interface.
func (*DeferredAttribute) Item() {}

func (*DeferredAttribute) LifeCycle() PropertyLifeCycle {
	return LifeCycleDeferred
}

// Implements the Propterty interface.
func (a *DeferredAttribute) Apply(ctx context.Context, w io.Writer) error {
	var rhs = ""
	value := a.fn(ctx)
	if len(value) > 0 {
		rhs = "=" + "\"" + value + "\""
	}
	_, err := fmt.Fprintf(w, " %s%s", a.key, rhs)
	return err
}

func DeferredAttr(key string, fn func(context.Context) string) *DeferredAttribute {
	return &DeferredAttribute{key: key, fn: fn}
}

// Creates a new Attribute with given key and value.
func Attr(key string, value string) *Attribute {
	return &Attribute{Key: key, Value: value}
}

type Component func(context.Context) Node

// Implements the Item interface.
func (Component) Item() {}

// Implements the Node interface.
func (Component) Release() {
	// no-op
}

// Implements the Node interface.
func (c Component) Render(cw ChunkWriter) {
	cw.Write(DynamicChunk(func(ctx context.Context, w io.Writer) error {
		writer := &teecw{ChunkWriter: cw, fn: func(c Chunk) error {
			if err := Render(c, ctx, w); err != nil {
				return err
			}
			return nil
		}}
		c(ctx).Render(writer)
		return nil
	}))
}

// Element represents an HTML element node, containing the element's name,
// a map of its attributes, a slice of child nodes, and a flag indicating
// whether the element is a void (self-closing) element.
type Element struct {
	nodes         []Node
	properties    []Property
	attributelist []*Attribute
	name          string
	meta          map[string]any
}

func (e *Element) Tag() string {
	if strings.HasPrefix(e.name, "!") {
		return e.name
	}
	return strings.ToLower(e.name)
}

func (e *Element) Set(key string, value any) {
	e.meta[key] = value
}

func (e *Element) Get(key string) any {
	return e.meta[key]
}

func (e *Element) SetAttribute(key, value string) {
	for i, attr := range e.attributelist {
		if attr.Key == key {
			e.attributelist[i].Value = value
			return
		}
	}
	e.attributelist = append(e.attributelist, &Attribute{Key: key, Value: value})
}

func (e *Element) AddClass(class string) {
	attr := e.GetAttribute("class")
	if len(attr) == 0 {
		attr = class
	} else {
		attr = attr + " " + class
	}
	e.SetAttribute("class", attr)
}

func (e *Element) GetAttribute(key string) string {
	for _, attr := range e.attributelist {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}

func (e *Element) GetAttributes() map[string]string {
	attrs := map[string]string{}
	for _, attr := range e.attributelist {
		attrs[attr.Key] = attr.Value
	}
	return attrs
}

func (e *Element) SetAttributes(list map[string]string) {
	e.attributelist = []*Attribute{}
	for key, value := range list {
		e.attributelist = append(e.attributelist, &Attribute{Key: key, Value: value})
	}
}

func (e *Element) GetNodes() []Node {
	return e.nodes
}

func (e *Element) AppendNode(node Node) {
	e.nodes = append(e.nodes, node)
}

// Implements the Item interface.
func (*Element) Item() {}

// Implements the Node interface.
func (e *Element) Release() {
	elpool.Put(e)
}

var voidelements = []string{"!DOCTYPE", "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "source", "track", "wbr"}

// Implements the Node interface.
func (e *Element) Render(w ChunkWriter) {
	w.Write(StaticChunk(fmt.Appendf(nil, "<%s", e.Tag())))

	if len(e.properties) > 0 {
		for _, prop := range e.properties {
			applier, ok := prop.(interface {
				Apply(context.Context, io.Writer) error
			})
			if !ok {
				log.Fatalf("property with deferred life cycle (%T) is not implementing the applier interface correctly, add a Apply(context.Context, io.Writer) error method", prop)
			}
			w.Write(DynamicChunk(applier.Apply))
		}
	}

	for _, attr := range e.attributelist {
		var rhs string
		if len(attr.Value) != 0 {
			rhs = "=" + "\"" + attr.Value + "\""
		}
		w.Write(StaticChunk(fmt.Appendf(nil, " %s%s", attr.Key, rhs)))
	}

	w.Write(StaticChunk(">"))

	if slices.Contains(voidelements, e.Tag()) {
		return
	}

	for _, node := range e.nodes {
		node.Render(w)
	}

	w.Write(StaticChunk(fmt.Appendf(nil, "</%s>", e.Tag())))
}

func (e *Element) Clone() *Element {
	attributelist := make([]*Attribute, len(e.attributelist))
	copy(attributelist, e.attributelist)

	return &Element{
		nodes:         e.nodes,
		properties:    e.properties,
		attributelist: attributelist,
		name:          e.name,
		meta:          e.meta,
	}
}

var propspool = sync.Pool{
	New: func() any {
		return &[]Property{}
	},
}

var elpool = sync.Pool{
	New: func() any {
		return &Element{
			nodes:         []Node{},
			properties:    []Property{},
			attributelist: []*Attribute{},
			meta:          make(map[string]any),
		}
	},
}

// El creates a new Element with the specified tag name and a variadic list of items,
// which can be either Node or Property types. Nodes are added as children of the element,
// while Properties are collected and applied to the element after all children are processed.
// Panics if an item of unknown type is provided.
// Returns a pointer to the constructed Element.
func El(name string, items ...Item) *Element {
	el := elpool.New().(*Element)
	el.name = name

	immediate := propspool.New().(*[]Property)
	defer propspool.Put(immediate)
	static := propspool.New().(*[]Property)
	defer propspool.Put(static)

	for _, item := range items {
		switch item := item.(type) {
		case Node:
			el.nodes = append(el.nodes, item)
		case Property:
			cycle := item.LifeCycle()
			switch cycle {
			case LifeCycleImmediate:
				*immediate = append(*immediate, item)
			case LifeCycleStatic:
				*static = append(*static, item)
			case LifeCycleDeferred:
				// TODO: Deferred properties have a serious bug
				el.properties = append(el.properties, item)
			default:
				panic(fmt.Sprintf("illegal property lifecycle: %T", item))
			}
		default:
			panic(fmt.Sprintf("illegal item type: %T", item))
		}
	}

	for _, prop := range *immediate {
		applier, ok := prop.(interface{ Apply(*Element) })
		if !ok {
			panic(fmt.Sprintf("property with immediate life cycle (%T) is not implementing the applier interface correctly, add a Apply(*Element) method", prop))
		}
		applier.Apply(el)
	}

	for _, prop := range *static {
		applier, ok := prop.(interface{ Apply(*Element) })
		if !ok {
			panic(fmt.Sprintf("property with static life cycle (%T) is not implementing the applier interface correctly, add a Apply(*Element) method", prop))
		}
		applier.Apply(el)
	}
	return el
}

// VoidEl creates a new void HTML element with the specified name and optional child items.
// Void elements are HTML elements that do not have closing tags (e.g., <img>, <br>, <input>).
// The function marks the created element as void and returns a pointer to the Element.
func VoidEl(name string, items ...Item) *Element {
	return El(name, items...)
}

// Map applies the provided function fn to each element of the input slice s,
// converting each element to a Node. It returns a single Node that is a Fragment
// containing all resulting child Nodes.
//
// E is a generic type parameter representing the element type of the input slice.
//
// Parameters:
//   - s: A slice of elements of type E.
//   - fn: A function that takes an element of type E and returns a Node.
//
// Returns:
//   - A Node that is a Fragment containing all Nodes produced by applying fn to each element of s.
func Map[E any](s []E, fn func(E) Node) Node {
	children := make([]Node, len(s))
	for i, item := range s {
		children[i] = fn(item)
	}
	return Fragment(children...)
}

// MapIdx applies the provided function fn to each element of the input slice s,
// along with each element's index converting each element to a Node. It returns
// a single Node that is a Fragment containing all resulting child Nodes.
//
// E is a generic type parameter representing the element type of the input slice.
//
// Parameters:
//   - s: A slice of elements of type E.
//   - fn: A function that takes an element of type E, an index of type int and
//     returns a Node.
//
// Returns:
//   - A Node that is a Fragment containing all Nodes produced by applying fn to each element of s.
func MapIdx[E any](s []E, fn func(E, int) Node) Node {
	children := make([]Node, len(s))
	for i, item := range s {
		children[i] = fn(item, i)
	}
	return Fragment(children...)
}

func Iter[T any](s iter.Seq[T], fn func(T) Node) Node {
	children := []Node{}
	for value := range s {
		children = append(children, fn(value))
	}
	return Fragment(children...)
}

func Iter2[K any, V any](s iter.Seq2[K, V], fn func(K, V) Node) Node {
	children := []Node{}
	for key, value := range s {
		children = append(children, fn(key, value))
	}
	return Fragment(children...)
}

// If returns the provided node if the condition is true; otherwise, it returns an empty Fragment node.
// This is useful for conditional rendering of nodes.
func If(cond bool, item Item) Item {
	if cond {
		return item
	}
	return Fragment()
}

// IfFn conditionally returns the result of the provided function as a Node.
// If cond is true, it calls and returns fn(); otherwise, it returns an empty Fragment Node.
// This is useful for conditional rendering of nodes.
func IfFn(cond bool, fn func() Item) Item {
	if cond {
		return fn()
	}
	return Fragment()
}

// SwitchCase represents a single case in a switch-like construct.
// It holds a value of type T to match against and a function Fn that returns a Node
// to be executed if the case matches. The generic type T must be comparable.
type SwitchCase[T comparable] struct {
	Value   T
	Fn      func() Item
	Default bool
}

// Switch iterates over the provided cases and returns the Node produced by the function
// of the first SwitchCase whose Value matches the given expr. If no cases match, it returns
// an empty Fragment Node. The generic type T must be comparable. Use Case and CaseFn
// functions for creating switch cases.
func Switch[T comparable](expr T, cases ...*SwitchCase[T]) Item {
	var d *SwitchCase[T]
	for _, c := range cases {
		if expr == c.Value {
			return c.Fn()
		} else if c.Default {
			d = c
		}
	}

	if d != nil {
		return d.Fn()
	}
	return Fragment()
}

// Case creates a new SwitchCase for the given value and node.
// It associates the value 'v' of type T with a function that returns the provided Node.
// This is used in switch-like constructs for rendering or processing nodes.
//
// T must be a comparable type.
// Parameters:
//
//	v    - the value to match in the switch case
//	node - the Node to return if the case matches
//
// Returns:
//
//	A pointer to a SwitchCase[T] containing the value and associated node function.
func Case[T comparable](v T, item Item) *SwitchCase[T] {
	return &SwitchCase[T]{
		Value: v,
		Fn: func() Item {
			return item
		},
	}
}

// Case creates a new SwitchCase for the given value and node.
// This is used in switch-like constructs for rendering or processing nodes.
// The function fn is executed if the case matches during a switch operation.
//
// T must be a comparable type.
// Parameters:
//
//	v    - the value to match in the switch case
//	fn - the function to run if the case matches
//
// Returns:
//
//	A pointer to a SwitchCase[T] containing the value and associated node function.
func CaseFn[T comparable](v T, fn func() Item) *SwitchCase[T] {
	return &SwitchCase[T]{
		Value: v,
		Fn:    fn,
	}
}

// Default creates a new SwitchCase for the given node that acts as the default case.
// The returned SwitchCase will always return the provided node when matched.
// This function is generic over type T, which must be comparable.
func Default[T comparable](item Item) *SwitchCase[T] {
	return &SwitchCase[T]{
		Fn: func() Item {
			return item
		},
		Default: true,
	}
}

// Default creates a new SwitchCase for the given node that acts as the default case.
// The returned SwitchCase will always return the provided node when matched.
// The function fn is executed if the default case matches during a switch operation.
// This function is generic over type T, which must be comparable.
func DefaultFn[T comparable](fn func() Item) *SwitchCase[T] {
	return &SwitchCase[T]{
		Fn:      fn,
		Default: true,
	}
}

/*
JSONNode represents a node containing arbitrary data to be serialized as JSON,
along with an optional indentation string for formatting the output.

Usage:

	html.Pre(html.JSON(map[string]any{ ... }))
*/
type JSONNode struct {
	Data   any
	Indent string
}

// Implements the Item interface.
func (*JSONNode) Item() {}

// Implements the Node interface.
func (*JSONNode) Release() {
	// no-op
}

// Applies indentation to the JSONNode and returns it back.
func (n *JSONNode) WithIndent(indent string) *JSONNode {
	n.Indent = indent
	return n
}

// Implements the Node interface.
func (n *JSONNode) Render(w ChunkWriter) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", n.Indent)
	if err := enc.Encode(n.Data); err != nil {
		panic(err)
	}
	w.Write(StaticChunk(buf.Bytes()))
}

// Creates a new JSONNode for serializing arbitrary json data.
func JSON(data any) *JSONNode {
	return &JSONNode{Data: data}
}

type RawUnsafe string

// Implements the Item interface.
func (RawUnsafe) Item() {}

// Implements the Node interface.
func (RawUnsafe) Release() {
	// no-op
}

// Implements the Node interface.
func (t RawUnsafe) Render(w ChunkWriter) {
	w.Write(StaticChunk(string(t)))
}
