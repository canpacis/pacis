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
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"iter"
	"slices"
	"strings"
)

// Item is an alias for any type.
type Item interface {
	Item()
}

// Node represents an element that can be rendered to an io.Writer within a given context.
// Implementations of Node should define the Render method to output their content.
type Node interface {
	Item
	Chunks() iter.Seq[Chunk]
}

type Chunk struct {
	fn   func(context.Context, io.Writer) error
	pure bool
}

func (c Chunk) Render(ctx context.Context, w io.Writer) error {
	return c.fn(ctx, w)
}

func (c Chunk) IsPure() bool {
	return c.pure
}

// Text represents a node containing plain text content within the HTML renderer.
type Text string

// Implements the Item interface.
func (Text) Item() {}

// Implements the Node interface.
func (t Text) Chunks() iter.Seq[Chunk] {
	return func(yield func(Chunk) bool) {
		yield(Chunk{
			fn: func(ctx context.Context, w io.Writer) error {
				_, err := w.Write([]byte(html.EscapeString(string(t))))
				return err
			},
			pure: true,
		})
	}
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
func (f Frag) Chunks() iter.Seq[Chunk] {
	return func(yield func(Chunk) bool) {
		for _, node := range f {
			for chunk := range node.Chunks() {
				if !yield(chunk) {
					return
				}
			}
		}
	}
}

// Fragment creates a Frag from the provided nodes, allowing multiple nodes to be grouped together.
// It accepts a variadic number of Node arguments and returns a Frag containing them.
func Fragment(nodes ...Node) Frag {
	return Frag(nodes)
}

// Property represents an interface for applying a property to an Element.
// Implementations of Property should define how the property modifies or affects the given Element.
type Property interface {
	Apply(*Element)
}

type Deferred func(context.Context) Property

// Implements the Item interface
func (Deferred) Item() {}

type Hook interface {
	Done(*Element)
}

// Attribute represents a key-value pair used as an attribute in an HTML node.
type Attribute struct {
	Key   string
	Value string
}

// Implements the Item interface.
func (*Attribute) Item() {}

// Implements the Propterty interface.
func (a *Attribute) Apply(el *Element) {
	if a.Key == "class" {
		el.ClassList.Add(a.Value)
	} else {
		el.SetAttribute(a.Key, a.Value)
	}
}

// Creates a new Attribute with given key and value.
func Attr(key string, value string) *Attribute {
	return &Attribute{Key: key, Value: value}
}

type Component func(context.Context) Node

// Implements the Item interface.
func (Component) Item() {}

// Implements the Node interface.
func (c Component) Chunks() iter.Seq[Chunk] {
	return func(yield func(Chunk) bool) {
		yield(Chunk{
			fn: func(ctx context.Context, w io.Writer) error {
				for chunk := range c(ctx).Chunks() {
					if err := chunk.Render(ctx, w); err != nil {
						return err
					}
				}
				return nil
			},
			pure: false,
		})
	}
}

// Element represents an HTML element node, containing the element's name,
// a map of its attributes, a slice of child nodes, and a flag indicating
// whether the element is a void (self-closing) element.
type Element struct {
	ClassList     *ClassList
	nodes         []Node
	deferreds     []Deferred
	attributes    map[string]string
	attributelist []*Attribute
	name          string
	void          bool
	ctx           context.Context
}

func (e *Element) Tag() string {
	if strings.HasPrefix(e.name, "!") {
		return e.name
	}
	return strings.ToLower(e.name)
}

func (e *Element) IsVoid() bool {
	return e.void
}

func (e *Element) Set(key, value any) {
	e.ctx = context.WithValue(e.ctx, key, value)
}

func (e *Element) Get(key any) any {
	return e.ctx.Value(key)
}

func (e *Element) SetAttribute(key, value string) {
	e.attributes[key] = value
	e.attributelist = append(e.attributelist, &Attribute{Key: key, Value: value})
}

func (e *Element) GetAttribute(key string) string {
	return e.attributes[key]
}

func (e *Element) GetAttributes() map[string]string {
	return e.attributes
}

func (e *Element) SetAttributes(list map[string]string) {
	e.attributes = list
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
func (e *Element) Chunks() iter.Seq[Chunk] {
	return func(yield func(Chunk) bool) {
		if !yield(Chunk{
			fn: func(ctx context.Context, w io.Writer) error {
				_, err := fmt.Fprintf(w, "<%s", e.Tag())
				return err
			},
			pure: true,
		}) {
			return
		}

		if len(e.deferreds) > 0 {
			for _, deferred := range e.deferreds {
				if !yield(Chunk{
					fn: func(ctx context.Context, w io.Writer) error {
						deferred(ctx).Apply(e)
						return nil
					},
					pure: false,
				}) {
					return
				}
			}
			if !yield(Chunk{
				fn: func(ctx context.Context, w io.Writer) error {
					if len(e.ClassList.Items) > 0 {
						e.SetAttribute("class", strings.Join(e.ClassList.Items, " "))
					}

					for _, attr := range e.attributelist {
						var rhs string
						if len(attr.Value) != 0 {
							rhs = "=" + "\"" + attr.Value + "\""
						}
						if _, err := fmt.Fprintf(w, " %s%s", attr.Key, rhs); err != nil {
							return err
						}
					}
					return nil
				},
				pure: false,
			}) {
				return
			}
		} else {
			if len(e.ClassList.Items) > 0 {
				e.SetAttribute("class", strings.Join(e.ClassList.Items, " "))
			}

			for _, attr := range e.attributelist {
				if !yield(Chunk{
					fn: func(ctx context.Context, w io.Writer) error {
						var rhs string
						if len(attr.Value) != 0 {
							rhs = "=" + "\"" + attr.Value + "\""
						}
						_, err := fmt.Fprintf(w, " %s%s", attr.Key, rhs)
						return err
					},
					pure: true,
				}) {
					return
				}
			}
		}

		if !yield(Chunk{
			fn: func(ctx context.Context, w io.Writer) error {
				_, err := fmt.Fprint(w, ">")
				return err
			},
			pure: true,
		}) {
			return
		}

		for _, node := range e.nodes {
			for chunk := range node.Chunks() {
				if !yield(chunk) {
					return
				}
			}
		}

		yield(Chunk{
			fn: func(ctx context.Context, w io.Writer) error {
				_, err := fmt.Fprintf(w, "</%s>", e.Tag())
				return err
			},
			pure: true,
		})
	}
}

// ClassList represents a collection of CSS class names associated with an HTML element.
type ClassList struct {
	Items []string
}

// Checks if a given class is present.
func (l *ClassList) Has(class string) bool {
	return slices.Contains(l.Items, class)
}

// Adds the give class.
func (l *ClassList) Add(class string) {
	l.Items = append(l.Items, class)
}

// Removes the given class.
func (l *ClassList) Remove(class string) {
	l.Items = slices.DeleteFunc(l.Items, func(v string) bool {
		return class == v
	})
}

// Toggles the given class, removing if it is present or adding if not.
func (l *ClassList) Toggle(class string) {
	if l.Has(class) {
		l.Remove(class)
	} else {
		l.Add(class)
	}
}

// El creates a new Element with the specified tag name and a variadic list of items,
// which can be either Node or Property types. Nodes are added as children of the element,
// while Properties are collected and applied to the element after all children are processed.
// Panics if an item of unknown type is provided.
// Returns a pointer to the constructed Element.
func El(name string, items ...Item) *Element {
	el := &Element{
		name:       name,
		attributes: make(map[string]string),
		ctx:        context.Background(),
		ClassList:  new(ClassList),
	}
	props := []Property{}
	hooks := []Hook{}

	for _, item := range items {
		switch item := item.(type) {
		case Deferred:
			el.deferreds = append(el.deferreds, item)
		case Node:
			el.nodes = append(el.nodes, item)
		case Property:
			props = append(props, item)
			if hook, ok := item.(Hook); ok {
				hooks = append(hooks, hook)
			}
		case Hook:
			hooks = append(hooks, item)
			if prop, ok := item.(Property); ok {
				props = append(props, prop)
			}
		default:
			panic(fmt.Sprintf("unknown item type: %T", item))
		}
	}

	for _, prop := range props {
		prop.Apply(el)
	}
	for _, hook := range hooks {
		hook.Done(el)
	}

	return el
}

// VoidEl creates a new void HTML element with the specified name and optional child items.
// Void elements are HTML elements that do not have closing tags (e.g., <img>, <br>, <input>).
// The function marks the created element as void and returns a pointer to the Element.
func VoidEl(name string, items ...Item) *Element {
	el := El(name, items...)
	el.void = true
	return el
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

// If returns the provided node if the condition is true; otherwise, it returns an empty Fragment node.
// This is useful for conditional rendering of nodes.
func If(cond bool, node Node) Node {
	if cond {
		return node
	}
	return Fragment()
}

// IfFn conditionally returns the result of the provided function as a Node.
// If cond is true, it calls and returns fn(); otherwise, it returns an empty Fragment Node.
// This is useful for conditional rendering of nodes.
func IfFn(cond bool, fn func() Node) Node {
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
	Fn      func() Node
	Default bool
}

// Switch iterates over the provided cases and returns the Node produced by the function
// of the first SwitchCase whose Value matches the given expr. If no cases match, it returns
// an empty Fragment Node. The generic type T must be comparable. Use Case and CaseFn
// functions for creating switch cases.
func Switch[T comparable](expr T, cases ...*SwitchCase[T]) Node {
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
func Case[T comparable](v T, node Node) *SwitchCase[T] {
	return &SwitchCase[T]{
		Value: v,
		Fn: func() Node {
			return node
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
func CaseFn[T comparable](v T, fn func() Node) *SwitchCase[T] {
	return &SwitchCase[T]{
		Value: v,
		Fn:    fn,
	}
}

// Default creates a new SwitchCase for the given node that acts as the default case.
// The returned SwitchCase will always return the provided node when matched.
// This function is generic over type T, which must be comparable.
func Default[T comparable](node Node) *SwitchCase[T] {
	return &SwitchCase[T]{
		Fn: func() Node {
			return node
		},
		Default: true,
	}
}

// Default creates a new SwitchCase for the given node that acts as the default case.
// The returned SwitchCase will always return the provided node when matched.
// The function fn is executed if the default case matches during a switch operation.
// This function is generic over type T, which must be comparable.
func DefaultFn[T comparable](fn func() Node) *SwitchCase[T] {
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

// Applies indentation to the JSONNode and returns it back.
func (n *JSONNode) WithIndent(indent string) *JSONNode {
	n.Indent = indent
	return n
}

// Implements the Node interface.
func (n *JSONNode) Render(ctx context.Context, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", n.Indent)
	return enc.Encode(n.Data)
}

// Creates a new JSONNode for serializing arbitrary json data.
func JSON(data any) *JSONNode {
	return &JSONNode{Data: data}
}
