package app

import (
	"fmt"
	"strconv"

	. "github.com/canpacis/pacis/docs/components"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
	"github.com/gosimple/slug"
	parser "github.com/sivukhin/godjot/djot_parser"
)

const imgsrc = "https://avatars.githubusercontent.com/u/37307107?s=400&u=54dd07c06503644ce385228881ea6e0b177c4d11&v=4"

type notification struct {
	title       string
	description string
}

var notifications = []notification{
	{"Your call has been confirmed.", "1 hour ago"},
	{"You have a new message!", "1 hour ago"},
	{"Your subscription is expiring soon!", "2 hours ago"},
}

var plates = map[string][]Node{
	"alert": {
		Alert(
			icons.Code(),
			AlertTitle(Text("Heads up!")),
			AlertDescription(Text("You can use Go to create great UI's")),
		),
	},
	"avatar": {
		Avatar(
			AvatarImage(Src(imgsrc)),
			AvatarFallback(Text("MC")),
		),
		Avatar(
			AvatarFallback(Text("MC")),
		),
		Frag(
			Avatar(
				AvatarSizeSm,

				AvatarImage(Src(imgsrc)),
				AvatarFallback(Text("MC")),
			),
			Avatar(
				AvatarImage(Src(imgsrc)),
				AvatarFallback(Text("MC")),
			),
			Avatar(
				AvatarSizeLg,

				AvatarImage(Src(imgsrc)),
				AvatarFallback(Text("MC")),
			),
		),
	},
	"badge": {
		Badge(Text("Badge")),
		Badge(BadgeVariantSecondary, Text("Badge")),
		Badge(BadgeVariantOutline, Text("Badge")),
		Badge(BadgeVariantDestructive, Text("Badge")),
	},
	"button": {
		Button(
			Text("Button"),
		),
		Button(
			ButtonVariantSecondary,
			Text("Button"),
		),
		Button(
			ButtonVariantOutline,
			Text("Button"),
		),
		Button(
			ButtonVariantDestructive,
			Text("Button"),
		),
		Button(
			ButtonVariantGhost,
			Text("Button"),
		),
		Button(
			ButtonVariantLink,
			Text("Button"),
		),
		Frag(
			Button(
				ButtonSizeSm,
				Text("Small"),
			),
			Button(
				Text("Default"),
			),
			Button(
				ButtonSizeLg,
				Text("Large"),
			),
		),
		Button(
			ButtonSizeIcon,
			ButtonVariantOutline,

			icons.EllipsisVertical(),
			Span(Text("Icon Button"), Class("sr-only")),
		),
		Button(
			Replace(A),
			Href("#button-as-link"),
			ButtonVariantOutline,

			Text("This is a link"),
		),
		Button(
			On("click", "alert('Clicked')"),

			Text("Press Me!"),
		),
	},
	"card": {
		Card(
			Class("w-fit sm:min-w-[380px]"),

			CardHeader(
				CardTitle(Text("Notifications")),
				CardDescription(Text("You have 3 unread messages.")),
			),
			CardContent(
				Class("grid gap-4"),

				Div(
					Class("flex items-center space-x-4 rounded-md border p-4"),

					icons.BellRing(),
					Div(
						Class("flex-1 space-y-1"),

						P(Class("text-sm font-medium leading-none line-clamp-1"), Text("Push Notifications")),
						P(
							Class("text-sm text-muted-foreground line-clamp-2"),

							Text("Send notifications to device."),
						),
					),
					Checkbox(Name("Enable Notifications"), Span(Class("sr-only"), Text("Enable Notifications"))),
				),
				Div(
					Map(notifications, func(n notification, i int) Node {
						return Div(
							Class("mb-4 grid grid-cols-[25px_1fr] items-start pb-4 last:mb-0 last:pb-0"),

							Span(Class("flex h-2 w-2 translate-y-1 rounded-full bg-sky-500")),
							Div(
								Class("space-y-1"),

								P(Class("text-sm font-medium leading-none"), Text(n.title)),
								P(Class("text-sm text-muted-foreground"), Text(n.description)),
							),
						)
					}),
				),
			),
			CardFooter(
				Button(
					Class("w-full"),

					icons.Check(),
					Text("Mark all as read"),
				),
			),
		),
	},
	"checkbox": {
		Checkbox(),
		Label("Label", Checkbox()),
		Checkbox(On("changed", "alert($event.detail.value)")),
	},
	"collapsible": {
		Collapsible(
			Class("min-w-[200px] flex flex-col gap-2 items-center"),

			CollapsibleTrigger(
				Button(Text("Trigger")),
			),
			CollapsibleContent(
				Div(Text("Collapsible Content")),
			),
		),
	},
	"dialog": {
		Dialog(
			DialogTrigger(
				Button(Text("Open Dialog")),
			),
			DialogContent(
				Class("max-w-[92dvw] sm:max-w-[420px]"),

				DialogHeader(
					DialogTitle(Text("Are you absolutely sure?")),
					DialogDescription(Text("This action cannot be undone. This will permanently delete your account and remove your data from our servers.")),
				),
			),
		),
		Dialog(
			DialogTrigger(
				Button(Text("Open Dialog")),
			),
			DialogContent(
				Class("max-w-[92dvw] sm:max-w-[420px]"),

				DialogHeader(
					DialogTitle(Text("Edit profile")),
					DialogDescription(Text("Make changes to your profile here. Click save when you're done.")),
				),
				Div(
					Class("flex flex-col gap-4 py-4"),

					Div(
						Class("grid grid-cols-4 items-center gap-4"),

						Label("Name", HtmlFor("name"), Class("text-right")),
						Input(ID("name"), Class("col-span-3")),
					),
					Div(
						Class("grid grid-cols-4 items-center gap-4"),

						Label("Username", HtmlFor("username"), Class("text-right")),
						Input(ID("username"), Class("col-span-3")),
					),
				),
				DialogFooter(
					Button(Type("submit"), Text("Save Changes")),
				),
			),
		),
	},
	"dropdown": {
		Dropdown(
			DropdownTrigger(
				Button(Text("Open Menu")),
			),
			DropdownContent(
				DropdownItem(
					ID("item-1"),

					icons.User(),
					Text("Profile"),
				),
				DropdownItem(
					ID("item-2"),

					icons.Settings(),
					Text("Settings"),
				),
			),
		),
	},
}

func RenderMarkup(node parser.TreeNode[parser.DjotNode], name string) Node {
	children := []I{}

	for _, child := range node.Children {
		children = append(children, RenderMarkup(child, name))
	}

	switch node.Type {
	case parser.DocumentNode:
		nodes := []Node{}
		for _, child := range children {
			node, ok := child.(Node)
			if ok {
				nodes = append(nodes, node)
			}
		}
		return Frag(nodes...)
	case parser.SectionNode:
		return Section(Join(children, Class("space-y-2"))...)
	case parser.TextNode:
		return Text(node.FullText())
	case parser.ParagraphNode:
		plate := node.Attributes.Get("plate")
		if len(plate) == 0 {
			return P(Join(children, Class(node.Attributes.Get("class")))...)
		}
		idx, err := strconv.Atoi(plate)
		if err != nil {
			panic(err)
		}

		plates, ok := plates[name]
		if ok {
			return Plate(plates[idx])
		}
		return Textf("Unknown plate %d in %s", idx, name)
	case parser.LinkNode:
		return A(Join(children, Class("text-sky-600 hover:text-sky-700 hover:underline"), Href(node.Attributes.Get("href")))...)
	case parser.CodeNode:
		return Code(string(node.FullText()), node.Attributes.Get("class"), Class("my-8"))
	case parser.StrongNode:
		return Span(Join(children, Class("font-semibold"))...)
	case parser.UnorderedListNode:
		return Ul(Join(children, Class("list-disc list-inside ml-4 leading-relaxed"))...)
	case parser.ListItemNode:
		return Li(children...)
	case parser.VerbatimNode:
		return Span(Join(children, Class("px-2 py-1 mx-1 inline bg-secondary font-mono text-sm rounded-sm font-semibold leading-8"))...)
	case parser.QuoteNode:
		return Div(
			Join(children,
				icons.Info(Class("size-4 shrink-0")),

				Class("border rounded-md p-4 text-sm flex gap-2 text-muted-foreground items-center my-4"),
			)...,
		)
	case parser.HeadingNode:
		level := len(node.Attributes.Get("$HeadingLevelKey"))
		txt := node.FullText()
		id := slug.Make(string(txt))

		switch level {
		case 1:
			return H1(Join(children, Class("scroll-m-20 text-3xl font-bold tracking-tight"))...)
		case 2:
			return Div(
				Class("my-4"),

				H2(
					Join(children, ID(id), Class("scroll-m-20 text-xl font-bold tracking-tight"))...,
				),
				Seperator(OHorizontal),
			)
		case 3:
			return H3(Join(children, ID(id), Class("scroll-m-20 text-lg font-bold mt-6"))...)
		default:
			return Text("unknown heading")
		}
	default:
		return Textf("unknown node type %s", node.Type)
	}
}

type TableOfContentItem struct {
	Label Node
	Href  string
	Order int
}

func ExtractTitles(node parser.TreeNode[parser.DjotNode]) []TableOfContentItem {
	links := []TableOfContentItem{}

	switch node.Type {
	case parser.DocumentNode, parser.SectionNode:
		for _, child := range node.Children {
			links = append(links, ExtractTitles(child)...)
		}
	case parser.HeadingNode:
		order := len(node.Attributes.Get("$HeadingLevelKey"))
		txt := node.FullText()
		id := slug.Make(string(txt))

		if order > 1 {
			links = append(links, TableOfContentItem{Text(txt), fmt.Sprintf("#%s", id), order})
		}
	}

	return links
}
