# Introduction

{.text-muted-foreground}
PacisUI is a set of utilities that are mainly UI components that help you build beautiful web interfaces with the Go programming language.

I plan to make PacisUI a part of a bigger set of tools that is Pacis, but for now I only have this documentation.

## What this is

This is a UI library built with [Go Language](https://go.dev/), [TailwindCSS](https://tailwindcss.com/) and [Alpine.js](https://alpinejs.dev/) and it comes in two pieces:

### The Renderer

PacisUI comes with its own html renderer for html elements and their attributes. Think of it like [templ](https://templ.guide/) or [gomponents](https://www.gomponents.com/). If you are familiar with the latter, you will find the PacisUI renderer familiar as well. It looks something like this;

```go
Div(
  ID("my-div") // An attribute

  P(Text("Hello, World!")) // A child with a child
)
```

You compose your html with go functions with PacisUI. If you are not sure about writing your html with Go functions, give it a try anyway, it might be for you and I believe you will find it very expressive and liberating. 

> Visit the [Syntax & Usage](/docs/syntax-usage) page to dive deep in this subject.

### The Components

The second piece and the focal point of this library is the components. These are, styled, interactive components that you would mainly be using.

A web app built with PacisUI ideally would use both these components and the builtin html elements along with some custom styling with tailwind. 

> The truth about frontend development is that these libraries are not a 'be all' solution. It will still require a considerable effort to create something beautiful.

### Icons

A *secret* third piece of the puzzle is the icons. PacisUI comes with a prebuilt icon library that is [Lucide](https://lucide.dev/). Icons are an essential part of UI development and I believe an out of the box solution is always needed.

That being said, you can always bring your own icons in the form of fonts, SVGs or plain images (although I recommend SVGs). I plan to create a better API to interact with *your* custom SVG icons in the future.

## How it works

A simple overview of this library is that it has nothing more than a bunch of functions that create a meaninful UI's.

The renderer provides some primitive interfaces and other stuff around it (like the components, icons or your own stuff) consume them.

These primitives are:

- `renderer.Renderer`: an interface that any renderer implements. If you have worked with [templ](https://templ.guide/) before, the signature will look familiar.


```go
type Renderer interface {
  Render(context.Context, io.Writer) error
}
```

- `renderer.I`: an alias to `renderer.Renderer` for ease of use.

```go
type I Renderer
```

- `renderer.Node`: represents an HTML node that is renderable, this can be anything from an element to a text node.

```go
type NodeType int

const (
	NodeText = NodeType(iota)
	NodeElement
	NodeFragment
)

type Node interface {
  Renderer
  NodeType() NodeType
}
```

- `renderer.Element`: represents and HTML element but not attributes or texts.

```go
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
```

- `renderer.Attribute`: represents any kind of element attribute.

```go
type Attribute interface {
	Renderer
	GetKey() string
	IsEmpty() bool
}
```

By composing these primitives and building up on them, you can create very well designed UI's with great developer experience. If you love Go like I do, building user interfaces with PacisUI should feel like a breath of fresh air.

> PacisUI is neither complete nor production ready yet but I am working on it. But this documentation site it built with it so it should give you an idea.