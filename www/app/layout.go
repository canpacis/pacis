package app

import (
	_ "embed"
	"os"
	"time"

	"github.com/canpacis/pacis/pages"
	fonts "github.com/canpacis/pacis/pages/font"
	"github.com/canpacis/pacis/pages/i18n"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

var sans = fonts.New("Inter", fonts.WeightList{fonts.W100, fonts.W900}, fonts.Swap, fonts.Latin, fonts.LatinExt)
var mono = fonts.New("JetBrains Mono", fonts.WeightList{fonts.W100, fonts.W800}, fonts.Swap, fonts.Latin, fonts.LatinExt)

//pacis:page label=robots
//go:embed robots.txt
var robots []byte

//pacis:page label=sitemap
//go:embed sitemap.xml
var sitemap []byte

//pacis:layout path=/
func Layout(ctx *pages.LayoutContext) I {
	locale, err := i18n.Locale(ctx)
	if err != nil {
		ctx.Logger().Error("failed to get locale", "error", err)
	}

	title := i18n.Text("title").String(ctx)
	desc := i18n.Text("desc").String(ctx)
	keywords := i18n.Text("keywords").String(ctx)

	user := pages.Get[*User](ctx, "user")
	appurl := os.Getenv("AppURL")
	banner := appurl + pages.Asset("banner.webp")

	return Html(
		Class(pages.Get[string](ctx, "theme")),
		Lang(locale.String()),

		Head(
			Meta(Name("title"), Content(title)),
			Meta(Name("description"), Content(desc)),
			Meta(Name("keywords"), Content(keywords)),
			Meta(Name("robots"), Content("index, follow")),
			Meta(HttpEquiv("Content-Type"), Content("text/html; charset=utf-8")),
			Meta(HttpEquiv("language"), Content("English")),
			Meta(HttpEquiv("author"), Content("canpacis")),

			Meta(Property("twitter:card"), Content("summary_large_image")),
			Meta(Property("twitter:url"), Content(appurl)),
			Meta(Property("twitter:title"), Content(title)),
			Meta(Property("twitter:description"), Content(desc)),
			Meta(Property("twitter:image"), Content(banner)),

			Meta(Property("og:type"), Content("website")),
			Meta(Property("og:url"), Content(appurl)),
			Meta(Property("og:title"), Content(title)),
			Meta(Property("og:description"), Content(desc)),
			Meta(Property("og:image"), Content(banner)),

			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0")),

			IfFn(user != nil, func() Renderer {
				return Store("user", user)
			}),
			If(user == nil, Store("user", &User{})),

			Script(Src(pages.Asset("before.ts"))),
			fonts.Head(sans, mono),
			ctx.Head(),
			Link(Href(pages.Asset("favicon.webp")), Rel("icon"), Type("image/png")),
			Script(
				Defer,
				Src("https://analytics.ui.canpacis.com/script.js"),
				Data("website-id", "4ce94416-1fb6-4a90-b881-f2f27a9736f7"),
			),
			Title(Text(title)),
		),
		Body(
			Class("flex flex-col min-h-dvh"),

			AppHeader(user),
			ctx.Outlet(),
			AppFooter(),
			ctx.Body(),
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
				{"/docs/radio", Text("Radio")},
				{"/docs/select", Text("Select")},
				// {"/docs/seperator", Text("Seperator")},
				// {"/docs/sheet", Text("Sheet")},
				// {"/docs/slider", Text("Slider")},
				{"/docs/switch", Text("Switch")},
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
							pages.A(
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

func AppHeader(user *User) Element {
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
						Span(Class("sr-only"), Text("Toggle Sidebar")),
					),
				),
				SheetContent(
					Class("overflow-scroll"),

					Navigation(sections, nil),
				),
			),
			pages.A(
				pages.Eager,
				Class("flex gap-3 items-center"),
				Href("/"),

				Img(Src(pages.Asset("logo.webp")), Width("24"), Height("24"), Class("w-6"), Alt("logo")),
				P(Class("font-semibold inline"), Text("Pacis")),
			),
			Ul(
				Class("hidden gap-4 lg:flex"),

				Map(links, func(link NavLink, i int) Node {
					return Li(
						Class("text-sm text-muted-foreground"),

						pages.A(Href(link.Href), link.Label),
					)
				}),
			),
			Div(
				Class("flex gap-1 items-center ml-auto"),

				Tooltip(
					"Github",
					time.Second*1,
					Button(
						ButtonSizeIcon,
						ButtonVariantGhost,
						Replace(pages.A),
						Href("https://github.com/canpacis/pacis-ui"),

						GithubIcon(),
						Span(Class("sr-only"), Text("Github")),
					),
				),
				Tooltip(
					"Toggle Theme",
					time.Second*1,
					Button(
						ButtonSizeIcon,
						ButtonVariantGhost,
						ToggleColorScheme,

						icons.Sun(),
						Span(Class("sr-only"), Text("Toggle Theme")),
					),
				),
				IfFn(user != nil, func() Renderer {
					return Div(
						Class("ml-2"),

						Dropdown(
							DropdownTrigger(
								Span(
									Class("cursor-pointer"),

									Avatar(
										AvatarImage(Src(user.Picture)),
										AvatarFallback(Text("MC")),
									),
								),
							),
							DropdownContent(
								Anchor(VBottom, HEnd, 8),

								DropdownLabel(user.Email),
								DropdownItem(Href("/auth/logout"), Replace(A), icons.LogOut(), Text("Logout")),
							),
						),
					)
				}),
			),
		),
	)
}

func AppFooter() Element {
	return Footer(
		Class("border-t border-dashed py-2 text-center h-[var(--footer-height)] fixed bottom-0 w-dvw bg-background z-50"),

		P(Class("text-sm text-muted-foreground"), Text("Built by "), pages.A(Href("https://canpacis.com"), Class("hover:underline"), Text("canpacis"))),
	)
}
