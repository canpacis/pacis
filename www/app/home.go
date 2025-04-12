package app

import (
	_ "embed"

	"github.com/canpacis/pacis/pages"
	fonts "github.com/canpacis/pacis/pages/font"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
)

var sans = fonts.New("Inter", fonts.WeightList{fonts.W100, fonts.W900}, fonts.Swap, fonts.Latin, fonts.LatinExt)
var mono = fonts.New("JetBrains Mono", fonts.WeightList{fonts.W100, fonts.W800}, fonts.Swap, fonts.Latin, fonts.LatinExt)

//go:embed robots.txt
var robots []byte

//go:embed sitemap.xml
var sitemap []byte

func Init() {
	pages.Robots(robots)
	pages.Sitemap(sitemap)
}

//pacis:layout path=/
func HomeLayout(ctx *pages.LayoutContext) I {
	return Html(
		Class(pages.Get[string](ctx, "theme")),

		Head(
			Meta(Name("title"), Content("Title")),
			Meta(Name("description"), Content("Description")),
			Meta(Name("keywords"), Content("Keywords")),
			Meta(Name("robots"), Content("index, follow")),
			Meta(HttpEquiv("Content-Type"), Content("text/html; charset=utf-8")),
			Meta(HttpEquiv("language"), Content("English")),
			Meta(HttpEquiv("author"), Content("canpacis")),

			Meta(Property("og:type"), Content("website")),
			Meta(Property("og:url"), Content("https://ui.canpacis.com")),
			Meta(Property("og:title"), Content("Title")),
			Meta(Property("og:description"), Content("Description")),
			Meta(Property("og:image"), Content("/public/banner.webp")),

			Meta(Property("twitter:card"), Content("summary_large_image")),
			Meta(Property("twitter:url"), Content("https://ui.canpacis.com")),
			Meta(Property("twitter:title"), Content("Title")),
			Meta(Property("twitter:description"), Content("Description")),
			Meta(Property("twitter:image"), Content("/public/banner.webp")),

			fonts.Head(sans, mono),
			ctx.Head(),
			Link(Href(pages.Asset("favicon.webp")), Rel("icon"), Type("image/png")),
			Title(Text("Title")),
		),
		Body(
			Class("flex flex-col min-h-dvh"),

			// AppHeader(),
			ctx.Outlet(),
			// AppFooter(),
		),
	)
}

//pacis:page path=/
func HomePage(ctx *pages.PageContext) I {
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
					Replace(A),
					Href("/docs/introduction"),
					Class("!rounded-full"),
					ButtonSizeLg,

					Text("Get Started"),
				),
				Button(
					ButtonSizeLg,
					ButtonVariantGhost,
					Replace(A),
					Class("!rounded-full"),
					Href("/docs/components"),

					Text("See Components"),
				),
			),
		),
		Div(),
	)
}
