# Button

{subtitle=""}
Displays a button or a component that looks like a button.

{plate="0"}
```go
Button(Text("Button"))
```

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/components"
)
```

```go
Button(Text("Button"))
```

## Examples

### Primary variant

{plate=0}
```go
Button(
	Text("Button"),
)
```

### Secondary variant

{plate=1}
```go
Button(
	ButtonVariantSecondary,
	
	Text("Secondary"),
)
```

### Outline variant

{plate=2}
```go
Button(
	ButtonVariantOutline,

	Text("Outline"),
)
```

### Destructive variant

{plate=3}
```go
Button(
	ButtonVariantDestructive,

	Text("Destructive"),
)
```

### Ghost variant

{plate=4}
```go
Button(
	ButtonVariantGhost,

	Text("Ghost"),
)
```

### Link variant

{plate=5}
```go
Button(
	ButtonVariantLink,

	Text("Link"),
)
```

### Sizes

{plate=6}
```go
Frag( // <- HTML Fragment
	Button(
		ButtonSizeSm,

		Text("Small"),
	),
	Button(
		Text("Default"),
	),
	Button(
		ButtonSizeLg,

		Text("Large"),
	)
)
```

### Icon

{plate=7}
```go
Button(
	ButtonSizeIcon,
	ButtonVariantOutline,

	icons.EllipsisVertical(),
)
```

### Button as link

{plate=8}
```go
Button(
	Replace(A), // <- Replace with an anchor tag
	Href("#button-as-link"), // <- Provide an href
	ButtonVariantOutline,

	Text("This is a link"),
)
```

### With an event handler

{plate=9}
```go
Button(
	On("click", "alert('Clicked')"), // <- This comes from the components module

	Text("Press Me!"),
)
```

## API

### Events

> All the events of a regular DOM button element is passed to this node