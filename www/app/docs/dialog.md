# Dialog

{.text-muted-foreground}
A window overlaid on either the primary window or another dialog window, rendering the content underneath inert.

{plate=0}
```go
Dialog(
  DialogTrigger(
    Button(Text("Open Dialog")),
  ),
  DialogContent(
    Class("max-w-[92dvw] sm:max-w-[420px]"),

    DialogHeader(
      DialogTitle(Text("Are you absolutely sure?")),
      DialogDescription(Text("This action cannot be undone. This will permanently delete your account and remove your data from our servers.")),
    ),
  ),
)
```

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/components"
)
```

```go
Dialog(
  DialogTrigger(
    Button(Text("Open Dialog")),
  ),
  DialogContent(
    DialogHeader(
      DialogTitle(Text("Are you absolutely sure?")),
      DialogDescription(Text("This action cannot be undone. This will permanently delete your account and remove your data from our servers.")),
    ),
  ),
)
```

## Examples

### Form dialog

{plate=1}
```go
Dialog(
  DialogTrigger(
    Button(Text("Open Dialog")),
  ),
  DialogContent(
    Class("sm:max-w-[425px]"),

    DialogHeader(
      DialogTitle(Text("Edit profile")),
      DialogDescription(Text("Make changes to your profile here. Click save when you're done.")),
    ),
    Div(
      Class("grid gap-4 py-4"),

      Div(
        Class("grid grid-cols-4, items-center gap-4"),

        Label(HtmlFor("name"), Class("text-right"), Text("Name")),
        Input(ID("name"), Class("col-span-3")),
      ),
      Div(
        Class("grid grid-cols-4, items-center gap-4"),

        Label(HtmlFor("username"), Class("text-right"), Text("Username")),
        Input(ID("username"), Class("col-span-3")),
      ),
    ),
    DialogFooter(
      Button(Type("submit"), Text("Save Changes")),
    ),
  ),
)
```

## API

### Events

| Event | Description |
|---|---|
| `init` | Fires upon initialization and sends dialogs open state. |
| `opened` | Fires when the dialog is opened. You can find the original event on the `$event.detail` object. |
| `closed` | Fires when the dialog is **explicitly** closed. You can find the original event and the data associated with the event on the `$event.detail` object. |
| `dismissed` | Fires when the dialog is dismissed rather than closed. You can find the original event on the `$event.detail` object. |

### Functions

| Name | Description |
|---|---|
| `openDialog()` | Opens the dialog. |
| `closeDialog(data: unknown)` | Closes the dialog with some data. |
| `dismissDialog()` | Dismisses the dialog. |

### Go Attributes

| Signature | Description |
|---|---|
| `OpenDialog` | Opens the dialog on click. |
| `OpenDialogOn(string)` | Opens the dialog upon given event. |
| `CloseDialog(D)` | Closes the dialog with the serializable data on click. |
| `CloseDialog(string, D)` | Closes the dialog with the serializble data upon given event. |
| `DismissDialog` | Dismisses the dialog on click. |
| `DismissDialogOn(string)` | Dismisses the dialog upon given event. |