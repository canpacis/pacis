# Alert

{subtitle=""}
Displays a callout for user attention.

{plate="0"}
Plate

## Usage

```go
import (
	. "github.com/canpacis/pacis-ui/components"
	"github.com/canpacis/pacis-ui/icons"
)
```

```go
Alert(
	icons.Code(),
	AlertTitle(Text("Heads up!")),
	AlertDescription(Text("You can us Go tho create great UI's")),
)
```