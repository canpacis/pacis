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
					icn.Icon("terminal"),
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

							icn.Icon("bell", icn.Width(18), icn.Height(18)),
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
					icn.Icon("search"),
					icn.Icon("search", icn.Width(50), icn.Height(50), icn.Stroke("red")),
					icn.Icon("search", icn.Width(50), icn.Height(50), icn.StrokeWidth(3)),
					icn.Icon("search", icn.Width(18), icn.Height(18), icn.StrokeWidth(1.4)),

					Dropdown(
						On("select", "console.log($event.detail)"),
						On("close", "console.log('closed')"),
						On("dismiss", "console.log('dismissed')"),

						DropdownTrigger(
							Button(Text("Dropdown")),
						),
						DropdownContent(
							DropdownItem(
								ID("item-1"),

								icn.Icon("user-round"),
								Text("Item 1"),
							),
							DropdownItem(
								ID("item-2"),

								icn.Icon("settings"),
								Text("Item 2"),
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
			Head(
				Link(Href("https://fonts.googleapis.com"), Rel("preconnect")),
				Link(Href("https://fonts.gstatic.com"), Rel("preconnect")),
				Link(Href("https://fonts.googleapis.com/css2?family=Inter:opsz,wght@14..32,100..900&display=swap"), Rel("stylesheet")),
				Link(Href("/public/main.css"), Rel("stylesheet")),
				Link(Href("https://rsms.me/inter/inter.css"), Rel("stylesheet")),
				Script(Src("https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js")),
				Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js")),
				Script(Src("https://unpkg.com/embla-carousel/embla-carousel.umd.js")),
			),
			Body(body),
		)
		fmt.Println(html.Render(r.Context(), w))
	})

	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	http.ListenAndServe(":8080", mux)
}
