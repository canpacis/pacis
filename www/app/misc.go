package app

import (
	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

//pacis:page label=not-found
func NotFoundPage(ctx *pages.PageContext) I {
	ctx.SetTitle("Not Found | Pacis")

	return Div(
		Class("flex flex-col container gap-6 flex-1 items-center justify-center"),

		H1(
			Class("text-xl md:text-3xl font-thin flex items-end gap-2 leading-7"),

			icons.FileSearch(Class("size-7"), icons.StrokeWidth(1)),
			Span(Text("Not Found!")),
		),
		P(
			Text("We couldn't find the page you were looking for"),
		),
		Button(
			Replace(pages.A),
			pages.Eager,
			Href("/"),
			Class("!rounded-full"),

			Text("Go Home"),
		),
	)
}

type AppErrorCode int

const (
	UnknownError = AppErrorCode(iota)
	InvalidAuthStateError
	AuthExchangeError
)

type AppError struct {
	Code        AppErrorCode
	Description string
}

func (ae *AppError) Error() string {
	return ae.Description
}

//pacis:page label=error
func ErrorPage(ctx *pages.PageContext) I {
	ctx.SetTitle("Error | Pacis")

	err, ok := pages.SafeGet[*AppError](ctx, "error")
	var code AppErrorCode
	if ok {
		code = err.Code
	}

	return Div(
		Class("flex flex-col container gap-6 flex-1 items-center justify-center"),

		H1(
			Class("text-xl md:text-3xl font-thin flex items-end gap-2 leading-7"),

			icons.TriangleAlert(Class("size-7"), icons.StrokeWidth(1)),
			Span(Text("Error!")),
		),
		P(
			SwitchCase(
				code,
				Case(InvalidAuthStateError, Text("There was an error with the auth state")),
				Case(AuthExchangeError, Text("Failed to exchange the auth token")),
				Case(UnknownError, Text("We don't know what happened")),
			),
		),
		Button(
			Replace(pages.A),
			pages.Eager,
			Href("/"),
			Class("!rounded-full"),

			Text("Go Home"),
		),
	)
}
