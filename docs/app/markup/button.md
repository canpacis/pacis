# Button

{subtitle=""}
Displays a button or a component that looks like a button.

{plate="0"}
Plate

## Usage

```go
import (
	. "github.com/canpacis/pacis-ui/components"
)
```

```go
Button(Text("Button"))
```

## Examples

### Primary variant

{plate=0}
Plate

```go
Button(
	Text("Button"),
)
```

### Secondary variant

{plate=1}
Plate

```go
Button(
	ButtonVariantSecondary,
	
	Text("Button"),
)
```

### Outline variant

{plate=2}
Plate

```go
Button(
	ButtonVariantOutline,

	Text("Button"),
)
```

### Destructive variant

{plate=3}
Plate

```go
Button(
	ButtonVariantDestructive,

	Text("Button"),
)
```

### Ghost variant

{plate=4}
Plate

```go
Button(
	ButtonVariantGhost,

	Text("Button"),
)
```

### Link variant

{plate=5}
Plate

```go
Button(
	ButtonVariantLink,

	Text("Button"),
)
```

### Sizes

{plate=6}
Plate

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
Plate

```go
Button(
	ButtonSizeIcon,
	ButtonVariantOutline,

	icons.EllipsisVertical(),
)
```

### Button as link

{plate=8}
Plate

```go
Button(
	Replace(A),
	Href("#button-as-link"),
	ButtonVariantOutline,

	Text("This is a link"),
)
```

### With an event handler

{plate=9}
Plate

```go
Button(
	On("click", "alert('Clicked')"), // <- This comes from the components module

	Text("Press Me!"),
)
```