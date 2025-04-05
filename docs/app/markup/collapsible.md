# Colapsible

{.text-muted-foreground}
An interactive component which expands/collapses a panel.

{plate=0}
Plate

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