# Avatar

{subtitle=""}
An image element with a fallback for representing the user.

{plate=0}
```go
Avatar(
	AvatarImage(Src(imgsrc)),
	AvatarFallback(Text("MC")),
)
```

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/components"
)
```

```go
Avatar(
	AvatarImage(Src(...)),
	AvatarFallback(Text("MC"))
)
```

## Examples

### With Image

{plate=0}
```go
Avatar(
	AvatarImage(Src(imgsrc)),
	AvatarFallback(Text("MC")),
)
```

### Without Image

{plate=1}
```go
Avatar(
	AvatarFallback(Text("MC")),
)
```

### Sizes

{plate=2}
```go
Frag(
	Avatar(
		AvatarSizeSm,

		AvatarImage(Src(imgsrc)),
		AvatarFallback(Text("MC")),
	),
	Avatar(
		AvatarImage(Src(imgsrc)),
		AvatarFallback(Text("MC")),
	),
	Avatar(
		AvatarSizeLg,

		AvatarImage(Src(imgsrc)),
		AvatarFallback(Text("MC")),
	),
)
```

## API

### Events

| Event | Description |
|---|---|
| `error` | Fires when given source image fails to load. |