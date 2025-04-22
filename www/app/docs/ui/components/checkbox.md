# Checkbox

{.text-muted-foreground}
A control that allows the user to toggle between checked and not checked.

{plate=0}
```go
Checkbox()
```

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/components"
)
```

```go
Checkbox()
```

## Examples

### With label

{plate=1}
```go
Checkbox(Text("Label"))
```

### Default checked

{plate=2}
```go
Checkbox(Checked)
```

### With an event handler

{plate=3}
```go
Checkbox(On("changed", "alert($event.detail.checked)")),
```

## API

### Events

| Event | Description |
|---|---|
| `init` | Fires upon initialization and sends its initial state. |
| `changed` | Fires when the checkbox state changes. You can find the `boolean` value on the `$event.detail` object |

### Functions

| Signature | Description |
|---|---|
| `toggleCheckbox(): void` | Toggles the checkbox state. |
| `isChecked(): boolean` | Returns the checkbox state. |

### Go Attributes

| Signature | Description |
|---|---|
| `ToggleCheckbox` | Toggles the checkbox on click. |
| `ToggleChecboxOn(string)` | Toggles the checkbox upon given event. |

### State

You can reach to a checkbox\'s state outside of the component by providing an explicit id to it.

```go
Checkbox(ID("test"))
// Somewhere else
Div(X("text", "$checkbox('test').isChecked()")) // <- use the api via the alpine magic
```

> Every checkbox, whether you provide an explicit id or not, is registered to this global store upon initialization.