package app

import (
	"embed"
	"path"
	"slices"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	. "github.com/canpacis/pacis/www/app/components"
	parser "github.com/sivukhin/godjot/djot_parser"
)

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

//pacis:layout path=/docs
func DocLayout(ctx *pages.LayoutContext) I {
	sections := getNavSections()
	items := getNavItems(sections)
	current := getCurrentItem(ctx, items)
	prev := getPrevItem(ctx, items)
	next := getNextItem(ctx, items)

	return Main(
		Class("container flex flex-1 items-start gap-4"),

		Aside(
			Class("hidden flex-col gap-2 border-r border-dashed py-4 pr-2 sticky overflow-auto top-[var(--header-height)] h-[calc(100dvh-var(--header-height)-var(--footer-height))] min-w-none lg:flex lg:min-w-[240px]"),

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

//go:embed docs
var docs embed.FS

//pacis:page path=/docs/{slug}
func DocsPage(ctx *pages.PageContext) I {
	slug := ctx.Request().PathValue("slug")
	source, err := docs.ReadFile(path.Join("docs", slug+".md"))
	if err != nil {
		return pages.NotFoundPage(ctx)
	}

	ast := parser.BuildDjotAst(source)
	nodes := []I{Class("flex-3")}
	nodes = append(nodes, RenderMarkup(ast[0], slug))
	headings := ExtractTitles(ast[0])

	return Div(
		Class("flex gap-8 flex-col-reverse xl:flex-row"),

		Div(nodes...),
		Div(
			Class("text-sm flex-1 h-fit leading-6 relative xl:sticky xl:top-[calc(var(--header-height)+2rem)]"),

			P(Class("font-semibold text-primary"), Text("On This Page")),
			Map(headings, func(item TableOfContentItem, i int) Node {
				return P(
					If(item.Order > 2, Class("ml-4")),

					A(
						Href(item.Href),
						Class("font-light text-muted-foreground hover:text-primary"),

						item.Label,
					),
				)
			}),
		),
	)
}
