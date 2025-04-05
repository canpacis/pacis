# Dropdown

{.text-muted-foreground}
Displays a menu to the user — such as a set of actions or functions — triggered by a button.

{plate=0}
Plate

## Usage

```go
import (
	. "github.com/canpacis/pacis-ui/components"
)
```

```go
Dropdown(
  DropdownTrigger(
    Button(Text("Open Menu")),
  ),
  DropdownContent(
    DropdownItem(
      ID("item-1"),

      icons.User(),
      Text("Profile"),
    ),
    DropdownItem(
      ID("item-2"),

      icons.Settings(),
      Text("Settings"),
    ),
  ),
)
```