package app

import (
	_ "embed"
	"os"
	"strings"
	"time"

	"github.com/canpacis/pacis/pages"
	fonts "github.com/canpacis/pacis/pages/font"
	"github.com/canpacis/pacis/pages/i18n"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
	"golang.org/x/text/language"
)

var sans = fonts.New("Inter", fonts.WeightList{fonts.W100, fonts.W900}, fonts.Swap, fonts.Latin, fonts.LatinExt)
var mono = fonts.New("JetBrains Mono", fonts.WeightList{fonts.W100, fonts.W800}, fonts.Swap, fonts.Latin, fonts.LatinExt)

//pacis:page label=robots
//go:embed robots.txt
var robots []byte

//pacis:page label=sitemap
//go:embed sitemap.xml
var sitemap []byte

type Layout struct {
	User   *User         `context:"user"`
	Locale *language.Tag `context:"locale"`
	Theme  string        `context:"theme"`
}

//pacis:layout path=/
func (l *Layout) Layout(ctx *pages.Context) I {
	pages.SetHeader(ctx, pages.NewHeader("Host", "canpacis.com"))

	title := i18n.Text("title").String(ctx)
	desc := i18n.Text("desc").String(ctx)
	keywords := strings.Split(i18n.Text("keywords").String(ctx), ",")

	appurl := os.Getenv("AppURL")
	banner := appurl + pages.Asset("banner.webp")

	pages.SetMetadata(ctx, &pages.Metadata{
		Title:       title,
		Description: desc,
		Keywords:    keywords,
		Robots:      "index, follow",
		Authors:     []string{"canpacis"},
		Language:    l.Locale.String(),
		Twitter: &pages.MetadataTwitter{
			Card:        "summary_large_image",
			URL:         appurl,
			Title:       title,
			Description: desc,
			Image:       banner,
		},
		OpenGraph: &pages.MetadataOG{
			Type:        "website",
			URL:         appurl,
			Title:       title,
			Description: desc,
			Image:       banner,
		},
	})

	return Html(
		Class(l.Theme),
		Lang(l.Locale.String()),

		Head(
			IfFn(l.User != nil, func() Renderer {
				return Store("user", l.User)
			}),
			If(l.User == nil, Store("user", &User{})),

			fonts.Head(sans, mono),
			pages.Head(ctx),
			Link(Href(pages.Asset("favicon.webp")), Rel("icon"), Type("image/png")),
			Script(
				Defer,
				Src("https://analytics.ui.canpacis.com/script.js"),
				Data("website-id", "4ce94416-1fb6-4a90-b881-f2f27a9736f7"),
			),
		),
		Body(
			Class("flex flex-col min-h-dvh overflow-x-hidden"),

			AppHeader(l.User),
			pages.Outlet(ctx),
			ToastContainer(),
			AppFooter(),
			pages.Body(ctx),
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
				Class("flex gap-3 items-center focusable p-2"),
				Href("/"),

				Img(Src(pages.Asset("logo.webp")), Width("24"), Height("24"), Class("w-6"), Alt("logo")),
				P(Class("font-semibold inline"), Text("Pacis")),
			),
			Ul(
				Class("hidden gap-4 lg:flex"),

				Map(links, func(link NavItem, i int) I {
					return Li(
						Class("text-sm text-muted-foreground"),

						pages.A(Href(link.Href), Class("focusable"), link.Label),
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

		P(Class("text-sm text-muted-foreground"), Text("Built by "), pages.A(Href("https://canpacis.com"), Class("hover:underline focusable"), Text("canpacis"))),
	)
}
