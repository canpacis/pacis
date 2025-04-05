# Dialog

{.text-muted-foreground}
A window overlaid on either the primary window or another dialog window, rendering the content underneath inert.

{plate=0}
Plate

## Usage

```go
import (
	. "github.com/canpacis/pacis-ui/components"
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
Plate

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