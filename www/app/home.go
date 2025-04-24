package app

import (
	_ "embed"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
)

//pacis:language default=en
//pacis:page path=/ middlewares=auth
func HomePage(ctx *pages.Context) I {
	// ctx.SetTitle("Homepage | Pacis")
	pages.SetHeader(ctx, pages.NewHeader("Host", "canpacis.com"))

	return Main(
		Class("container my-8 lg:my-16 flex-1 flex flex-col lg:flex-row items-start md:items-center gap-8 mt:0 lg:-mt-[var(--footer-height)]"),

		Div(
			Class("flex-0 lg:flex-3"),

			H1(
				Class("text-2xl font-bold leading-tight tracking-tighter sm:text-3xl md:text-4xl lg:leading-[1.1]"),

				Text("Build web applications "),
				Span(Class("relative inline-block px-2 py-1 rounded-lg bg-gradient-to-r from-blue-300 via-purple-300 to-pink-300 dark:from-blue-500 dark:via-purple-500 dark:to-pink-500"), Text("with Go")),
			),
			P(
				Class("max-w-2xl text-base font-light text-foreground sm:text-lg mt-4"),

				Text("Build stunning, modern UIs for Go applications with ease, intuitive components, flexible styling, and seamless performance."),
			),
			Div(
				Class("mt-8 flex gap-2"),

				Button(
					pages.Eager,
					Replace(pages.A),
					Href("/docs"),
					Class("!rounded-full"),
					ButtonSizeLg,

					Text("Get Started"),
				),
				Button(
					ButtonSizeLg,
					ButtonVariantGhost,
					pages.Eager,
					Replace(pages.A),
					Class("!rounded-full"),
					Href("/docs/components"),

					Text("See Components"),
				),
			),
		),
		Div(),
	)
}
