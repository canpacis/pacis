# Dropdown

{.text-muted-foreground}
Displays a menu to the user — such as a set of actions or functions — triggered by a button.

{plate=0}
```go
Dropdown(
  DropdownTrigger(
    Button(Text("Open Menu")),
  ),
  DropdownContent(
    DropdownItem(
      ID("profile"),

      icons.User(),
      Text("Profile"),
    ),
    DropdownItem(
      ID("settings"),

      icons.Settings(),
      Text("Settings"),
    ),
  ),
)
```

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/html"
	. "github.com/canpacis/pacis/ui/components"
)
```

```go
Dropdown(
  DropdownTrigger(
    Button(Text("Open Menu")),
  ),
  DropdownContent(
    DropdownItem(
      ID("item-id"),

      Text("Dropdown Item"),
    ),
  ),
)
```