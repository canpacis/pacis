package app

import (
	"slices"

	. "github.com/canpacis/pacis/docs/components"
	"github.com/canpacis/pacis/pages"
	fonts "github.com/canpacis/pacis/pages/font"
	"github.com/canpacis/pacis/pages/i18n"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

var sans = fonts.New("Inter", fonts.WeightList{fonts.W100, fonts.W900}, fonts.Swap, fonts.Latin, fonts.LatinExt)

func Layout(ctx *pages.LayoutContext) I {
	return Html(
		Class(pages.Get[string](ctx, "theme")),

		Head(
			Meta(Name("title"), Content(i18n.Text("title").String(ctx))),
			Meta(Name("description"), Content(i18n.Text("desc").String(ctx))),
			Meta(Name("keywords"), Content(i18n.Text("keywords").String(ctx))),
			Meta(Name("robots"), Content("index, follow")),
			Meta(HttpEquiv("Content-Type"), Content("text/html; charset=utf-8")),
			Meta(HttpEquiv("language"), Content("English")),
			Meta(HttpEquiv("author"), Content("canpacis")),

			Meta(Property("og:type"), Content("website")),
			Meta(Property("og:url"), Content("https://ui.canpacis.com")),
			Meta(Property("og:title"), Content(i18n.Text("title").String(ctx))),
			Meta(Property("og:description"), Content(i18n.Text("desc").String(ctx))),
			Meta(Property("og:image"), Content("/public/banner.png")),

			Meta(Property("twitter:card"), Content("summary_large_image")),
			Meta(Property("twitter:url"), Content("https://ui.canpacis.com")),
			Meta(Property("twitter:title"), Content(i18n.Text("title").String(ctx))),
			Meta(Property("twitter:description"), Content(i18n.Text("desc").String(ctx))),
			Meta(Property("twitter:image"), Content("/public/banner.png")),

			fonts.Head(sans),
			ctx.Head(),
			Link(Href("/public/main.css"), Rel("stylesheet")),
			Link(Href("/public/favicon.png"), Rel("icon"), Type("image/png")),
			Title(i18n.Text("title")),
		),
		Body(
			Class("flex flex-col min-h-dvh"),

			AppHeader(),
			ctx.Outlet(),
			AppFooter(),
		),
	)
}

type NavLink struct {
	Href  string
	Label Node
}

type NavSection struct {
	Label Node
	Items []NavLink
}

func getNavItems(sections []NavSection) []NavLink {
	items := []NavLink{}

	for _, section := range sections {
		items = append(items, section.Items...)
	}
	return items
}

func getCurrentItem(ctx *pages.LayoutContext, items []NavLink) *NavLink {
	path := ctx.Request().URL.Path

	idx := slices.IndexFunc(items, func(item NavLink) bool {
		return item.Href == path
	})
	if idx < 0 {
		return nil
	}
	return &items[idx]
}

func getNextItem(ctx *pages.LayoutContext, items []NavLink) *NavLink {
	path := ctx.Request().URL.Path
	idx := slices.IndexFunc(items, func(item NavLink) bool {
		return item.Href == path
	})
	if idx < 0 {
		return nil
	}
	if idx >= len(items)-1 {
		return nil
	}

	return &items[idx+1]
}

func getPrevItem(ctx *pages.LayoutContext, items []NavLink) *NavLink {
	path := ctx.Request().URL.Path
	idx := slices.IndexFunc(items, func(item NavLink) bool {
		return item.Href == path
	})
	if idx <= 0 {
		return nil
	}

	return &items[idx-1]
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
				// {"/docs/select", Text("Select")},
				// {"/docs/seperator", Text("Seperator")},
				// {"/docs/sheet", Text("Sheet")},
				// {"/docs/slider", Text("Slider")},
				// {"/docs/switch", Text("Switch")},
				// {"/docs/tabs", Text("Tabs")},
				// {"/docs/textarea", Text("Textarea")},
				// {"/docs/toast", Text("Toast")},
				// {"/docs/tooltip", Text("Tooltip")},
			},
		},
	}
}

func Navigation(sections []NavSection, current *NavLink) Node {
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
						A(
							Class("rounded-md block text-sm w-full px-2.5 py-1.5 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 cursor-pointer"),
							If(current == nil, Href(item.Href)),
							IfFn(current != nil, func() Renderer {
								return If(current.Href != item.Href, Href(item.Href))
							}),

							item.Label,
						),
					)
				}),
			),
		)
	})
}

func DocLayout(ctx *pages.LayoutContext) I {
	sections := getNavSections()
	items := getNavItems(sections)
	current := getCurrentItem(ctx, items)
	prev := getPrevItem(ctx, items)
	next := getNextItem(ctx, items)

	return Main(
		Class("container flex flex-1 items-start gap-4"),

		Aside(
			Class("hidden flex-col gap-2 border-r border-dashed py-4 pr-2 sticky top-[var(--header-height)] h-[calc(100dvh-var(--header-height)-var(--footer-height))] min-w-none lg:flex lg:min-w-[240px]"),

			Navigation(sections, current),
		),
		Section(
			Class("py-8 flex-1 w-full ml-0 lg:ml-4 xl:ml-8"),

			IfFn(current != nil, func() Renderer {
				return Breadcrumb(
					Class("mb-4"),

					BreadcrumbItem(Text("Docs")),
					BreadcrumbSeperator(),
					BreadcrumbItem(current.Label),
				)
			}),
			ctx.Outlet(),
			Div(
				Class("flex gap-8 mb-[var(--footer-height)]"),

				Div(
					Class("flex mt-12 flex-3 w-full xl:w-fit"),

					IfFn(prev != nil, func() Renderer {
						return DocButton(prev.Href, false, prev.Label)
					}),
					IfFn(next != nil, func() Renderer {
						return DocButton(next.Href, true, next.Label)
					}),
				),
				Div(Class("flex-1 hidden xl:block")),
			),
		),
	)
}

func AppHeader() Element {
	links := []NavLink{
		{"/docs/introduction", Text("Docs")},
		{"/docs/components", Text("Components")},
	}
	sections := getNavSections()

	return Header(
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

				Img(Src("/public/logo.png"), Class("w-6")),
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
					Replace(Knot),
					Href("https://github.com/canpacis/pacis-ui"),

					GithubIcon(),
				),
				Button(
					ButtonSizeIcon,
					ButtonVariantGhost,
					ToggleColorScheme,

					icons.Sun(),
				),
			),
		),
	)
}

func GithubIcon(props ...I) Element {
	return El("svg",
		Join(
			props,
			Attr("viewBox", "0 0 100 100"),

			El("path",
				Attr("fill-rule", "evenodd"),
				Attr("clip-rule", "evenodd"),
				Attr("fill", "currentColor"),
				Attr("d", "M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z"),
			),
		)...,
	)
}

func AppFooter() Element {
	return Footer(
		Class("border-t border-dashed py-2 text-center h-[var(--footer-height)] fixed bottom-0 w-dvw bg-background"),

		P(Class("text-sm text-muted-foreground"), Text("Built by "), Knot(Href("https://canpacis.com"), Class("hover:underline"), Text("canpacis"))),
	)
}
