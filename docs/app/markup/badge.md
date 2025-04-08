# Badge

{subtitle=""}
Displays a badge or a component that looks like a badge.

{plate=0}
```go
Badge(Text("Badge"))
```

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/components"
)
```

```go
Badge(Text("Badge"))
```

## Examples

### Primary

{plate=0}
```go
Badge(Text("Badge"))
```

### Secondary

{plate=1}
```go
Badge(BadgeVariantSecondary, Text("Secondary"))
```

### Outline

{plate=2}
```go
Badge(BadgeVariantOutline, Text("Outline"))
```

### Destructive

{plate=3}
```go
Badge(BadgeVariantDestructive, Text("Destructive"))
```