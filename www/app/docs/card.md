# Card

{subtitle=""}
Displays a card with header, content, and footer.

{plate=0}
```go
Card(
  Class("w-fit sm:min-w-[380px]"),

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

        P(Class("text-sm font-medium leading-none line-clamp-1"), Text("Push Notifications")),
        P(
          Class("text-sm text-muted-foreground line-clamp-2"),

          Text("Send notifications to device."),
        ),
      ),
      Checkbox(Name("Enable Notifications"), Span(Class("sr-only"), Text("Enable Notifications"))),
    ),
    Div(
      Map(notifications, func(n notification, i int) Node {
        return Div(
          Class("mb-4 grid grid-cols-[25px_1fr] items-start pb-4 last:mb-0 last:pb-0"),

          Span(Class("flex h-2 w-2 translate-y-1 rounded-full bg-sky-500")),
          Div(
            Class("space-y-1"),

            P(Class("text-sm font-medium leading-none"), Text(n.title)),
            P(Class("text-sm text-muted-foreground"), Text(n.description)),
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

## Usage

```go
import (
	. "github.com/canpacis/pacis/ui/components"
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

{plate=0}
```go
import (
	. "github.com/canpacis/pacis/ui/html"
	. "github.com/canpacis/pacis/ui/components"
	"github.com/canpacis/pacis/ui/icons"
)

// ...

type Notification struct {
	Title       string
	Description string
}

var notifications = []Notification{
	{"Your call has been confirmed.", "1 hour ago"},
	{"You have a new message!", "1 hour ago"},
	{"Your subscription is expiring soon!", "2 hours ago"},
}

//...

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