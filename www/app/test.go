package app

import (
	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/html"
)

//pacis:page path=/test middlewares=auth,limiter
func TestPage(ctx *pages.PageContext) I {
	return Div(
		Class("container"),

		&pages.Hook{},
	)
}
