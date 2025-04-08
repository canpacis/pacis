# Colapsible

{.text-muted-foreground}
An interactive component which expands/collapses a panel.

{plate=0}
```go
Collapsible(
  Class("min-w-[200px] flex flex-col gap-2 items-center"),

  CollapsibleTrigger(
    Button(Text("Trigger")),
  ),
  CollapsibleContent(
    Div(Text("Collapsible Content")),
  ),
)
```

## Usage

```go
import (
	. "github.com/canpacis/pacis-ui/components"
)
```

```go
Collapsible(
  CollapsibleTrigger(
    Button(Text("Trigger")),
  ),
  CollapsibleContent(
    Div(Text("Collapsible Content")),
  ),
)
```

## Examples

### Simple collapsible

{plate=0}
```go
Collapsible(
  Class("min-w-[200px] flex flex-col gap-2 items-center"),

  CollapsibleTrigger(
    Button(Text("Trigger")),
  ),
  CollapsibleContent(
    Div(Text("Collapsible Content")),
  ),
)
```

## API

### Events

| Event | Description |
|---|---|
| `init` | Fires upon initialization and sends its initial state. |
| `changed` | Fires when the collapse state changes. You can find the `boolean` value on the `$event.detail` object |

### Functions

| Name | Description |
|---|---|
| `toggleCollapse()` | Toggles the collapse state. |

### Go Attributes

| Signature | Description |
|---|---|
| `ToggleCollapse` | Toggles the collapsible on click. |
| `ToggleCollapseOn(string)` | Toggles the collapsible upon given event. |