# Card

{subtitle=""}
Displays a card with header, content, and footer.

{plate=0}
Plate

## Usage

```go
import (
	. "github.com/canpacis/pacis-ui/components"
)
```

```go
Card(
  CardHeader(
    CardTitle(Text("Title")),
    CardDescription(Text("Description")),
  ),
  CardContent(
    // Content
  ),
  CardFooter(
    // Footer
  ),
),
```
## Examples

### Notifcation

```go
import (
	. "github.com/canpacis/pacis-ui/components"
	"github.com/canpacis/pacis-ui/icons"
)
```

```go
type Notification struct {
	Title       string
	Description string
}

var notifications = []Notification{
	{"Your call has been confirmed.", "1 hour ago"},
	{"You have a new message!", "1 hour ago"},
	{"Your subscription is expiring soon!", "2 hours ago"},
}
```

```go
Card(
  Class("w-[380px]"),

  CardHeader(
    CardTitle(Text("Notifications")),
    CardDescription(Text("You have 3 unread messages.")),
  ),
  CardContent(
    Class("grid gap-4"),

    Div(
      Class("flex items-center space-x-4 rounded-md border p-4"),

      icons.BellRing(),
      Div(
        Class("flex-1 space-y-1"),

        P(Class("text-sm font-medium leading-none"), Text("Push Notifications")),
        P(
          Class("text-sm text-muted-foreground"),

          Text("Send notifications to device."),
        ),
      ),
      Checkbox(),
    ),
    Div(
      Map(notifications, func(n Notification, i int) Node {
        return Div(
          Class("mb-4 grid grid-cols-[25px_1fr] items-start pb-4 last:mb-0 last:pb-0"),

          Span(Class("flex h-2 w-2 translate-y-1 rounded-full bg-sky-500")),
          Div(
            Class("space-y-1"),

            P(Class("text-sm font-medium leading-none"), Text(n.Title)),
            P(Class("text-sm text-muted-foreground"), Text(n.Description)),
          ),
        )
      }),
    ),
  ),
  CardFooter(
    Button(
      Class("w-full"),

      icons.Check(),
      Text("Mark all as read"),
    ),
  ),
)
```