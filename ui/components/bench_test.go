package components_test

import (
	"bytes"
	"context"
	"testing"

	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

type NavLink struct {
	Href  string
	Label Node
}

type NavSection struct {
	Label Node
	Items []NavLink
}

func getNavSections() []NavSection {
	return []NavSection{
		{
			Label: Text("Getting Started"),
			Items: []NavLink{
				{"/docs/introduction", Text("Introduction")},
				// {"/docs/installation", Text("Installation")},
				// {"/docs/quick-start", Text("Quick Start")},
				// {"/docs/syntax-usage", Text("Syntax & Usage")},
				// {"/docs/roadmap", Text("Roadmap")},
			},
		},
		{
			Label: Text("Components"),
			Items: []NavLink{
				{"/docs/alert", Text("Alert")},
				{"/docs/avatar", Text("Avatar")},
				{"/docs/badge", Text("Badge")},
				{"/docs/button", Text("Button")},
				{"/docs/card", Text("Card")},
				// {"/docs/carousel", Text("Carousel")},
				{"/docs/checkbox", Text("Checkbox")},
				{"/docs/collapsible", Text("Collapsible")},
				{"/docs/dialog", Text("Dialog")},
				{"/docs/dropdown", Text("Dropdown")},
				{"/docs/input", Text("Input")},
				{"/docs/label", Text("Label")},
				// {"/docs/radio", Text("Radio")},
				{"/docs/select", Text("Select")},
				// {"/docs/seperator", Text("Seperator")},
				// {"/docs/sheet", Text("Sheet")},
				// {"/docs/slider", Text("Slider")},
				// {"/docs/switch", Text("Switch")},
				{"/docs/tabs", Text("Tabs")},
				// {"/docs/textarea", Text("Textarea")},
				// {"/docs/toast", Text("Toast")},
				// {"/docs/tooltip", Text("Tooltip")},
			},
		},
	}
}

func Navigation(sections []NavSection, current *NavLink) Node {
	iscurr := func(href string) bool {
		if current == nil {
			return false
		}
		return current.Href == href
	}

	return Map(sections, func(heading NavSection, i int) Node {
		return Div(
			Class("mb-4"),

			H2(
				Class("font-semibold text-sm text-muted-foreground mb-2"),

				heading.Label,
			),

			Ul(
				Map(heading.Items, func(item NavLink, i int) Node {
					return Li(
						If(!iscurr(item.Href),
							A(
								Class("rounded-md block text-sm w-full px-2.5 py-1.5 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 cursor-pointer"),
								Href(item.Href),

								item.Label,
							),
						),
						If(iscurr(item.Href),
							Span(
								Class("rounded-md block text-sm w-full px-2.5 py-1.5 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 cursor-pointer"),

								item.Label,
							),
						),
					)
				}),
			),
		)
	})
}

func BenchmarkRender(b *testing.B) {
	sections := getNavSections()
	links := []NavLink{
		{"/docs/introduction", Text("Docs")},
		{"/docs/components", Text("Components")},
	}

	for b.Loop() {
		html := Html(
			Class("dark"),
			Lang("en"),

			Head(
				Meta(Name("title"), Content("Title")),
				Meta(Name("description"), Content("Description")),
				Meta(Name("keywords"), Content("keywords")),
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

				Link(Href("/public/main.css"), Rel("stylesheet")),
				Link(Href("/public/favicon.webp"), Rel("icon"), Type("image/png")),
				Title(Text("Title")),
			),
			Body(
				Class("flex flex-col min-h-dvh"),

				Header(
					Class("py-3 border-b border-dashed sticky top-0 z-50 w-full bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 h-[var(--header-height)]"),

					Div(
						Class("flex container items-center gap-4 lg:gap-8 h-full"),

						Sheet(
							Class("block lg:hidden"),

							SheetTrigger(
								Button(
									ButtonSizeIcon,
									ButtonVariantGhost,

									icons.PanelLeft(),
									Span(Class("sr-only"), Text("Toggle Sidebar")),
								),
							),
							SheetContent(
								Class("overflow-scroll"),

								Navigation(sections, nil),
							),
						),
						A(
							Class("flex gap-3 items-center"),
							Href("/"),

							Img(Src("/public/logo.webp"), Width("24"), Height("24"), Class("w-6"), Alt("logo")),
							P(Class("font-semibold inline"), Text("Pacis")),
						),
						Ul(
							Class("hidden gap-4 lg:flex"),

							Map(links, func(link NavLink, i int) Node {
								return Li(
									Class("text-sm text-muted-foreground"),

									A(Href(link.Href), link.Label),
								)
							}),
						),
						Div(
							Class("flex gap-1 items-center ml-auto"),

							Button(
								ButtonSizeIcon,
								ButtonVariantGhost,
								ToggleColorScheme,

								icons.Sun(),
								Span(Class("sr-only"), Text("Toggle Theme")),
							),
						),
					),
				),
				Main(
					Class("container flex flex-1 items-start gap-4"),

					Aside(
						Class("hidden flex-col gap-2 border-r border-dashed py-4 pr-2 sticky overflow-auto top-[var(--header-height)] h-[calc(100dvh-var(--header-height)-var(--footer-height))] min-w-none lg:flex lg:min-w-[240px]"),

						Navigation(sections, nil),
					),
					Section(
						Class("py-8 flex-1 w-full ml-0 lg:ml-4 xl:ml-8"),

						Button(ButtonSizeDefault),
						Button(ButtonSizeSm),
						Button(ButtonSizeLg),
						Button(ButtonSizeIcon),
						Button(ButtonVariantDefault),
						Button(ButtonVariantDestructive),
						Button(ButtonVariantGhost),
						Button(ButtonVariantLink),
						Button(ButtonVariantOutline),
						Button(ButtonVariantSecondary),

						Button(ButtonSizeSm, ButtonVariantDefault),
						Button(ButtonSizeSm, ButtonVariantDestructive),
						Button(ButtonSizeSm, ButtonVariantGhost),
						Button(ButtonSizeSm, ButtonVariantLink),
						Button(ButtonSizeSm, ButtonVariantOutline),
						Button(ButtonSizeSm, ButtonVariantSecondary),
						Button(ButtonSizeLg, ButtonVariantDefault),
						Button(ButtonSizeLg, ButtonVariantDestructive),
						Button(ButtonSizeLg, ButtonVariantGhost),
						Button(ButtonSizeLg, ButtonVariantLink),
						Button(ButtonSizeLg, ButtonVariantOutline),
						Button(ButtonSizeLg, ButtonVariantSecondary),
						Button(ButtonSizeIcon, ButtonVariantDefault),
						Button(ButtonSizeIcon, ButtonVariantDestructive),
						Button(ButtonSizeIcon, ButtonVariantGhost),
						Button(ButtonSizeIcon, ButtonVariantLink),
						Button(ButtonSizeIcon, ButtonVariantOutline),
						Button(ButtonSizeIcon, ButtonVariantSecondary),
						IfFn(true, func() Renderer {
							return Breadcrumb(
								Class("mb-4"),

								BreadcrumbItem(Text("Docs")),
								BreadcrumbSeperator(),
								BreadcrumbItem(Text("Label")),
							)
						}),
						Div(
							Class("flex gap-8 flex-col-reverse xl:flex-row"),

							Div(
								Class("text-sm flex-1 h-fit leading-6 relative xl:sticky xl:top-[calc(var(--header-height)+2rem)]"),

								P(Class("font-semibold text-primary"), Text("On This Page")),
							),
						),
						Div(
							Class("flex gap-8 mb-[var(--footer-height)]"),

							Div(
								Class("flex mt-12 flex-3 w-full xl:w-fit"),

								IfFn(true, func() Renderer {
									return Button(Text("true"))
								}),
								IfFn(false, func() Renderer {
									return Button(Text("false"))
								}),
							),
							Div(Class("flex-1 hidden xl:block")),
						),
					),
				),
				Footer(
					Class("border-t border-dashed py-2 text-center h-[var(--footer-height)] fixed bottom-0 w-dvw bg-background"),

					P(Class("text-sm text-muted-foreground"), Text("Built by "), A(Href("https://canpacis.com"), Class("hover:underline"), Text("canpacis"))),
				),
			),
		)

		var buf = new(bytes.Buffer)
		html.Render(context.Background(), buf)
	}
}
