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
			ToastContainer(),
			AppFooter(),
			ctx.Body(),
		),
	)
}

func AppHeader(user *User) Element {
	links := []NavItem{
		{Label: Text("Docs"), Href: "/docs"},
		{Label: Text("Components"), Href: "/docs/components"},
	}
	items := getnavitems()

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

					Map(items, NavItemUI),
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

				Map(links, func(link NavItem, i int) I {
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
					"Discord",
					time.Second*1,
					Button(
						ButtonSizeIcon,
						ButtonVariantGhost,
						Replace(pages.A),
						Href("https://discord.gg/QnXQjYZrJU"),

						DiscordIcon(),
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
		Class("border-t border-dashed py-2 text-center h-[var(--footer-height)] fixed bottom-0 w-dvw bg-background z-40"),

		P(Class("text-sm text-muted-foreground"), Text("Built by "), pages.A(Href("https://canpacis.com"), Class("hover:underline"), Text("canpacis"))),
	)
}
