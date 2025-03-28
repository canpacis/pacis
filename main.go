package main

import (
	"fmt"
	"net/http"

	. "github.com/canpacis/pacis/components"
	icn "github.com/canpacis/pacis/icons"
	. "github.com/canpacis/pacis/renderer"
)

type Notification struct {
	Title       string
	Description string
	Unread      bool
}

func main() {
	notifications := []Notification{
		{Title: "Your call has been confirmed.", Description: "1 hour ago", Unread: true},
		{Title: "You have a new message!", Description: "1 hour ago"},
		{Title: "Your subscription is expiring soon!", Description: " hours ago"},
	}

	body := Main(
		Class("flex flex-col gap-4 my-8"),

		Div(
			Class("container flex flex-col gap-2"),

			Div(
				Class("flex gap-2 items-center"),

				Alert(
					icn.Terminal(),
					AlertTitle(Text("Heads Up!")),
					AlertDescription(
						Text("Lorem ipsum dolor sit amet consectetur"),
					),
				),
				Avatar(
					AvatarImage(
						Src("https://cloud.appwrite.io/v1/storage/buckets/67dc16070005054cb3c3/files/67dc16220011a2637b55/view?project=ksalt&mode=admin"),
					),
					AvatarFallback(Text("MC")),
				),
			),

			Div(
				Class("flex gap-2 items-start"),

				Card(
					Class("max-w-[340px]"),

					CardHeader(
						CardTitle(
							Class("flex gap-2"),

							icn.Bell(icn.Width(18), icn.Height(18)),
							Text("Notifications"),
						),
						CardDescription(Text("You have 3 unread messages.")),
					),
					CardContent(
						Class("grid gap-4"),

						Div(
							Map(notifications, func(item Notification, i int) Node {
								return Div(
									Class("mb-4 grid grid-cols-[25px_1fr] items-start pb-4 last:mb-0 last:pb-0"),

									Span(Class("flex h-2 w-2 translate-y-1 rounded-full bg-sky-500")),
									Div(
										Class("space-y-1"),

										P(
											Class("text-sm font-medium leading-none"),

											Text(item.Title),
											If(item.Unread, Badge(Class("ml-2"), Text("New"))),
										),
										P(Class("text-sm text-muted-foreground"), Text(item.Description)),
									),
								)
							}),
						),
						Checkbox(Text("Label")),
					),
					CardFooter(
						Button(
							Class("w-full"),

							Text("Mark all as read"),
						),
					),
				),
				Div(
					Class("flex gap-2 items-center"),

					Collapsible(
						CollapsibleTrigger(Button(Text("Click to collapse"))),
						CollapsibleContent(
							Div(Text("Content")),
						),
					),
					icn.Search(),
					icn.Search(icn.Width(50), icn.Height(50), icn.Stroke("red")),
					icn.Search(icn.Width(50), icn.Height(50), icn.StrokeWidth(3)),
					icn.Search(icn.Width(18), icn.Height(18), icn.StrokeWidth(1.4)),

					Dropdown(
						On("select", "console.log($event.detail)"),
						On("close", "console.log('closed')"),
						On("dismiss", "console.log('dismissed')"),

						DropdownTrigger(
							Button(Text("Dropdown")),
						),
						DropdownContent(
							Anchor(VBottom, HCenter, 8),

							DropdownItem(
								ID("item-1"),

								icn.UserRound(),
								Text("Item 1"),
							),
							DropdownItem(
								ID("item-2"),

								icn.Settings(),
								Text("Item 2"),
							),
						),
					),
				),
			),
			Div(
				Dialog(
					DialogTrigger(
						Button(Text("Open Dialog")),
					),
					DialogContent(
						Class("sm:max-w-[425px]"),

						DialogHeader(
							DialogTitle(Text("Edit profile")),
							DialogDescription(Text("Make changes to your profile here. Click save when you're done.")),
						),
						Div(
							Class("grid gap-4 py-4"),

							Div(
								Class("grid grid-cols-4 items-center gap-4"),

								Label(HtmlFor("name"), Class("text-right"), Text("Name")),
								Input(ID("name"), Class("col-span-3")),
							),
							Div(
								Class("grid grid-cols-4 items-center gap-4"),

								Label(HtmlFor("username"), Class("text-right"), Text("Username")),
								Input(ID("username"), Class("col-span-3"), Value("canpacis")),
							),
						),
						DialogFooter(
							Button(
								On("click", "closeDialog"),

								Text("Save changes"),
							),
						),
					),
				),
			),
		),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		html := Html(
			Head(AppHead()),
			Body(body),
		)
		fmt.Println(html.Render(r.Context(), w))
	})

	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	http.ListenAndServe(":8080", mux)
}
