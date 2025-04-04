package app

import (
	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
)

func HomePage(ctx *pages.PageContext) I {
	return Main(
		Class("container my-16 flex-1"),

		H1(
			Class("text-2xl font-bold leading-tight tracking-tighter sm:text-3xl md:text-4xl lg:leading-[1.1]"),

			Text("Build great UI's with Go"),
		),
		P(
			Class("max-w-2xl text-base font-light text-foreground sm:text-lg mt-4"),

			Text("Build stunning, modern UIs for Go applications with ease, intuitive components, flexible styling, and seamless performance."),
		),
		Div(
			Class("mt-8 flex gap-2"),

			Button(
				Replace(A),
				Href("/docs/introduction"),
				ButtonSizeSm,

				Text("Get Started"),
			),
			Button(
				ButtonSizeSm,
				ButtonVariantGhost,
				Replace(A),
				Href("/docs/components"),

				Text("See Components"),
			),
		),
	)
}

func Introduction(ctx *pages.PageContext) I {
	return Main(Text("Introduction"))
}

func Installation(ctx *pages.PageContext) I {
	return Main(Text("Installation"))
}

func AlertPage(ctx *pages.PageContext) I {
	return Main(Text("AlertPage"))
}
