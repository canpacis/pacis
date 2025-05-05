package app

import (
	"errors"
	"os"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

//pacis:page label=not-found
func NotFoundPage(ctx *pages.Context) I {
	// ctx.SetTitle("Not Found | Pacis")

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

var ErrGenericAppError = errors.New("app error")

type AppError struct {
	code   AppErrorCode
	err    error
	status int
}

func (ae *AppError) Error() string {
	return ae.err.Error()
}

func (ae *AppError) SetError(err error) {
	ae.err = err
}

func (ae *AppError) Status() int {
	return ae.status
}

func (ae *AppError) SetStatus(status int) {
	ae.status = status
}

func (ae *AppError) Unwrap() error {
	return ae.err
}

//pacis:page label=error
func (p *AppError) Page(ctx *pages.Context) I {
	// ctx.SetTitle("Error | Pacis")
	var message string

	if os.Getenv("ENVIRONMENT") == "development" {
		message = p.Error()
	} else {
		message = "We don't know what happened"
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
				p.code,
				Case(InvalidAuthStateError, Text("There was an error with the auth state")),
				Case(AuthExchangeError, Text("Failed to exchange the auth token")),
				Case(UnknownError, Text(message)),
			),
		),
		Button(
			Replace(pages.A),
			Href("/"),
			Class("!rounded-full"),

			Text("Go Home"),
		),
	)
}

func NewAppError(code AppErrorCode, err error, status int) *AppError {
	return &AppError{code: code, err: err, status: status}
}
