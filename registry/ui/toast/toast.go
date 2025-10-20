package toast

import (
	"fmt"
	"time"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/x"
)

func New(message string, duration time.Duration) html.Property {
	return x.Data(fmt.Sprintf("toast('%s', %d)", message, duration.Milliseconds()))
}

var Default = time.Second * 2

var Show = ShowOn("click")

func ShowOn(event string) *html.Attribute {
	return html.Attr("x-on:"+event, "show()")
}
