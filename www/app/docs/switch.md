# Switch

{.text-muted-foreground}
A control that allows the user to toggle between checked and not checked.

{plate=0}
```go
Switch()
```

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/components"
)
```

```go
Switch()
```

## Examples

### With label

{plate=1}
```go
Switch(Text("Label"))
```

### Default checked

{plate=2}
```go
Switch(Checked)
```

### With an event handler

{plate=3}
```go
Switch(On("changed", "alert($event.detail.checked)")),
```

## API

### Events

| Event | Description |
|---|---|
| `init` | Fires upon initialization and sends its initial state. |
| `changed` | Fires when the switch state changes. You can find the `boolean` value on the `$event.detail` object |

### Functions

| Signature | Description |
|---|---|
| `toggleSwitch(): void` | Toggles the siwtch state. |
| `isChecked(): boolean` | Returns the siwtch state. |

### Go Attributes

| Signature | Description |
|---|---|
| `ToggleSwitch` | Toggles the switch on click. |
| `ToggleSwitchOn(string)` | Toggles the switch upon given event. |

### State

You can reach to a switch\'s state outside of the component by providing an explicit id to it.

```go
Switch(ID("test"))
// Somewhere else
Div(X("text", "$switch_('test').isChecked()")) // <- use the api via the alpine magic
```

> Every switch, whether you provide an explicit id or not, is registered to this global store upon initialization.