# Checkbox

{.text-muted-foreground}
A control that allows the user to toggle between checked and not checked.

{plate=0}
Plate

## Usage

```go
import (
	. "github.com/canpacis/pacis-ui/components"
)
```

```go
Checkbox()
```

## Examples

### With label

{plate=1}
Plate

```go
Checkbox(Text("Label"))
```

### With an event handler

{plate=2}
Plate

```go
Checkbox(On("changed", "alert($event.detail.value)")),
```