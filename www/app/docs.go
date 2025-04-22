package app

import (
	"embed"
	"io"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	. "github.com/canpacis/pacis/www/app/components"
	parser "github.com/sivukhin/godjot/djot_parser"
)

type NavItem struct {
	Label Node
	Href  string
	Items []NavItem
}

func getnavitems() []NavItem {
	return []NavItem{
		{
			Label: Text("Getting Started"),
			Items: []NavItem{
				{Label: Text("Introduction"), Href: "/docs/getting-started/introduction"},
				{Label: Text("Installation"), Href: "/docs/getting-started/installation"},
				{Label: Text("Quick Start"), Href: "/docs/getting-started/quick-start"},
				{Label: Text("Contributing"), Href: "/docs/getting-started/contributing"},
			},
		},
		{
			Label: Text("Pages"),
			Items: []NavItem{
				{Label: Text("Overview"), Href: "/docs/pages/overview"},
				{Label: Text("Pages"), Href: "/docs/pages/pages"},
				{Label: Text("Layouts"), Href: "/docs/pages/layouts"},
				{Label: Text("Actions"), Href: "/docs/pages/actions"},
				{Label: Text("Font"), Href: "/docs/pages/font"},
				{Label: Text("I18n"), Href: "/docs/pages/i18n"},
				{Label: Text("Middleware"), Href: "/docs/pages/middleware"},
				{Label: Text("Prefetching"), Href: "/docs/pages/prefetching"},
				{Label: Text("Streaming"), Href: "/docs/pages/streaming"},
			},
		},
		{
			Label: Text("UI"),
			Items: []NavItem{
				{
					Label: Text("Templating"),
					Items: []NavItem{
						{Label: Text("Elements"), Href: "/docs/ui/templating/elements"},
						{Label: Text("Attributes"), Href: "/docs/ui/templating/attributes"},
						{Label: Text("Utilities"), Href: "/docs/ui/templating/utilities"},
						{Label: Text("Extending"), Href: "/docs/ui/templating/extending"},
					},
				},
				{
					Label: Text("Icons"),
					Items: []NavItem{
						{Label: Text("Icon Set"), Href: "/docs/ui/icons/icon-set"},
						{Label: Text("Custom Icons"), Href: "/docs/ui/icons/custom-icons"},
					},
				},
				{
					Label: Text("Components"),
					Items: []NavItem{
						{Label: Text("Overview"), Href: "/docs/ui/components/overview"},
						{Label: Text("Alert"), Href: "/docs/ui/components/alert"},
						{Label: Text("Avatar"), Href: "/docs/ui/components/avatar"},
						{Label: Text("Badge"), Href: "/docs/ui/components/badge"},
						{Label: Text("Button"), Href: "/docs/ui/components/button"},
						{Label: Text("Calendar"), Href: "/docs/ui/components/calendar"},
						{Label: Text("Card"), Href: "/docs/ui/components/card"},
						{Label: Text("Carousel"), Href: "/docs/ui/components/carousel"},
						{Label: Text("Checkbox"), Href: "/docs/ui/components/checkbox"},
						{Label: Text("Code"), Href: "/docs/ui/components/code"},
						{Label: Text("Collapsible"), Href: "/docs/ui/components/collapsible"},
						{Label: Text("Dialog"), Href: "/docs/ui/components/dialog"},
						{Label: Text("Dropdown"), Href: "/docs/ui/components/dropdown"},
						{Label: Text("Input"), Href: "/docs/ui/components/input"},
						{Label: Text("Label"), Href: "/docs/ui/components/label"},
						{Label: Text("Radio"), Href: "/docs/ui/components/radio"},
						{Label: Text("Select"), Href: "/docs/ui/components/select"},
						{Label: Text("Seperator"), Href: "/docs/ui/components/seperator"},
						{Label: Text("Sheet"), Href: "/docs/ui/components/sheet"},
						{Label: Text("Slider"), Href: "/docs/ui/components/slider"},
						{Label: Text("Switch"), Href: "/docs/ui/components/switch"},
						{Label: Text("Table"), Href: "/docs/ui/components/table"},
						{Label: Text("Tabs"), Href: "/docs/ui/components/tabs"},
						{Label: Text("Textarea"), Href: "/docs/ui/components/textarea"},
						{Label: Text("Toast"), Href: "/docs/ui/components/toast"},
						{Label: Text("Tooltip"), Href: "/docs/ui/components/tooltip"},
					},
				},
			},
		},
	}
}

func flattenitems(items []NavItem) []NavItem {
	result := []NavItem{}
	for _, item := range items {
		if len(item.Href) > 0 {
			result = append(result, item)
		}
		result = append(result, flattenitems(item.Items)...)
	}
	return result
}

func getnextitem(ctx *pages.LayoutContext, items []NavItem) *NavItem {
	path := ctx.Request().URL.Path
	items = flattenitems(items)
	idx := slices.IndexFunc(items, func(item NavItem) bool {
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

func getprevitem(ctx *pages.LayoutContext, items []NavItem) *NavItem {
	path := ctx.Request().URL.Path
	items = flattenitems(items)
	idx := slices.IndexFunc(items, func(item NavItem) bool {
		return item.Href == path
	})
	if idx <= 0 {
		return nil
	}

	return &items[idx-1]
}

func parsepath(p string) []string {
	result := []string{}
	dir, file := path.Split(p)
	if len(file) > 0 {
		result = append(result, parsepath(strings.TrimRight(dir, "/"))...)
		result = append(result, file)
	}
	return result
}

func capitalize(s string) string {
	if len(s) < 3 {
		return strings.ToUpper(s)
	}
	return strings.ReplaceAll(strings.ToUpper(string(s[0]))+string(s[1:]), "-", " ")
}

func NavItemUI(item NavItem, i int) I {
	if len(item.Href) == 0 {
		return Div(
			Class("mb-4"),

			H2(
				Class("font-semibold text-sm text-muted-foreground mb-2"),

				item.Label,
			),
			Div(Map(item.Items, NavItemUI)),
		)
	}

	return Li(
		pages.A(
			Class("rounded-md block text-sm w-full px-2.5 py-1.5 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 cursor-pointer"),
			Href(item.Href),

			item.Label,
		),
	)
}

//pacis:layout path=/docs
//pacis:redirect from=/docs/ to=/docs/getting-started/introduction
//pacis:redirect from=/docs/components to=/docs/ui/components/overview
func DocLayout(ctx *pages.LayoutContext) I {
	items := getnavitems()
	prev := getprevitem(ctx, items)
	next := getnextitem(ctx, items)
	parts := parsepath(ctx.Request().URL.Path)

	return Main(
		Class("container flex flex-1 items-start gap-4"),

		Aside(
			Class("hidden flex-col gap-2 border-r border-dashed py-4 pr-2 sticky overflow-auto top-[var(--header-height)] h-[calc(100dvh-var(--header-height)-var(--footer-height))] min-w-none lg:flex lg:min-w-[240px]"),

			Map(items, NavItemUI),
		),
		Section(
			Class("py-8 flex-1 flex flex-col w-full ml-0 lg:ml-4 xl:ml-8"),

			Breadcrumb(
				Class("mb-4"),

				Map(parts, func(part string, i int) I {
					return Frag(
						BreadcrumbItem(Text(capitalize(part))),
						If(i != len(parts)-1, BreadcrumbSeperator()),
					)
				}),
			),
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

type DocPage struct {
	Title    string
	Nodes    []I
	Contents []*TableOfContentItem
}

type DocDir struct {
	Title string
	Pages map[string]*DocPage
	Dirs  map[string]*DocDir
}

func ExtractDocs(fsys fs.ReadDirFS, root string) (*DocDir, error) {
	entries, err := fsys.ReadDir(root)
	if err != nil {
		return nil, err
	}

	dir := &DocDir{
		Title: path.Base(root),
		Pages: make(map[string]*DocPage),
		Dirs:  make(map[string]*DocDir),
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			subdir, err := ExtractDocs(fsys, path.Join(root, name))
			if err != nil {
				return nil, err
			}
			dir.Dirs[subdir.Title] = subdir
		} else {
			slug, _ := strings.CutSuffix(name, path.Ext(name))
			file, err := fsys.Open(path.Join(root, entry.Name()))
			if err != nil {
				return nil, err
			}
			defer file.Close()

			src, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			ast := parser.BuildDjotAst(src)
			nodes := []I{Class("flex-3")}
			markup, err := RenderMarkup(ast[0], slug)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, markup)
			contents := ExtractTitles(ast[0])

			dir.Pages[slug] = &DocPage{
				Title:    capitalize(slug),
				Nodes:    nodes,
				Contents: contents,
			}
		}
	}

	return dir, nil
}

//go:embed docs
var docsfs embed.FS

var dir *DocDir

func InitDocs() error {
	var err error
	dir, err = ExtractDocs(docsfs, "docs")
	return err
}

//pacis:page path=/docs/getting-started/{slug} middlewares=auth,limiter
func GettingStartedPage(ctx *pages.PageContext) I {
	slug := ctx.Request().PathValue("slug")
	page, ok := dir.Dirs["getting-started"].Pages[slug]

	if ok {
		ctx.SetTitle("Pacis Docs | Getting Started > " + page.Title)
	} else {
		return ctx.NotFound()
	}

	return DocPageUI(page)
}

//pacis:page path=/docs/pages/{slug} middlewares=auth,limiter
func PagesPage(ctx *pages.PageContext) I {
	slug := ctx.Request().PathValue("slug")
	page, ok := dir.Dirs["pages"].Pages[slug]

	if ok {
		ctx.SetTitle("Pacis Docs | Pages > " + page.Title)
	} else {
		return ctx.NotFound()
	}

	return DocPageUI(page)
}

//pacis:page path=/docs/ui/templating/{slug} middlewares=auth,limiter
func TemplatingPage(ctx *pages.PageContext) I {
	slug := ctx.Request().PathValue("slug")
	page, ok := dir.Dirs["ui"].Dirs["templating"].Pages[slug]

	if ok {
		ctx.SetTitle("Pacis Docs | UI > Templating > " + page.Title)
	} else {
		return ctx.NotFound()
	}

	return DocPageUI(page)
}

//pacis:page path=/docs/ui/components/{slug} middlewares=auth,limiter
func ComponentsPage(ctx *pages.PageContext) I {
	slug := ctx.Request().PathValue("slug")
	page, ok := dir.Dirs["ui"].Dirs["components"].Pages[slug]

	if ok {
		ctx.SetTitle("Pacis Docs | UI > Components > " + page.Title)
	} else {
		return ctx.NotFound()
	}

	return DocPageUI(page)
}

func DocPageUI(page *DocPage) Node {
	return Div(
		Class("flex gap-8 flex-col-reverse xl:flex-row"),

		Div(page.Nodes...),
		Div(
			Class("text-sm flex-1 h-fit leading-6 relative xl:sticky xl:top-[calc(var(--header-height)+2rem)]"),

			P(Class("font-semibold text-primary"), Text("On This Page")),
			Map(page.Contents, func(item *TableOfContentItem, i int) I {
				return P(
					If(item.Order > 2, Class("ml-4")),

					pages.A(
						Href(item.Href),
						Class("font-light text-muted-foreground hover:text-primary"),

						item.Label,
					),
				)
			}),
		),
	)
}
