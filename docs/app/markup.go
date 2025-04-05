package app

import (
	"os"
	"strings"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/html"
	parser "github.com/sivukhin/godjot/djot_parser"
)

func MarkupPage(path string) pages.Page {
	source, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ast := parser.BuildDjotAst(source)
	nodes := []I{Class("flex-3")}
	nodes = append(nodes, RenderMarkup(ast[0], strings.TrimSuffix(strings.TrimPrefix(path, "./docs/app/markup/"), ".md")))
	headings := ExtractTitles(ast[0])

	return func(pc *pages.PageContext) I {
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
}
