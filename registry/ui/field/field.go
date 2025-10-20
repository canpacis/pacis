package field

import (
	"fmt"

	"components/ui/separator"

	"github.com/canpacis/pacis/components"
	"github.com/canpacis/pacis/html"
)

func Set(items ...html.Item) html.Node {
	return html.Fieldset(
		components.ItemsOf(
			items,
			html.Data("slot", "field-set"),
			html.Class("flex flex-col gap-6 has-[>[data-slot=checkbox-group]]:gap-3 has-[>[data-slot=radio-group]]:gap-3"),
		)...,
	)
}

func Legend(items ...html.Item) html.Node {
	return html.Legend(
		components.ItemsOf(
			items,
			html.Data("slot", "field-legend"),
			html.Data("variant", "label"),
			html.Class("mb-3 font-medium data-[variant=legend]:text-base data-[variant=label]:text-sm"),
		)...,
	)
}

func Group(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "field-group"),
			html.Class("group/field-group @container/field-group flex w-full flex-col gap-7 data-[slot=checkbox-group]:gap-3 [&>[data-slot=field-group]]:gap-4"),
		)...,
	)
}

type Orientation = components.Variant

const (
	Vertical = Orientation(iota)
	Horizontal
	Responsive
)

var orientation = components.NewVariantApplier(func(el *html.Element, v components.Variant) {
	switch v {
	case Vertical:
		el.AddClass("flex-col [&>*]:w-full [&>.sr-only]:w-auto")
		el.SetAttribute("data-orientation", "vertical")
	case Horizontal:
		el.AddClass("flex-row items-center [&>[data-slot=field-label]]:flex-auto has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px has-[>[data-slot=field-content]]:items-start")
		el.SetAttribute("data-orientation", "horizontal")
	case Responsive:
		el.AddClass("@md/field-group:flex-row @md/field-group:items-center @md/field-group:[&>*]:w-auto flex-col [&>*]:w-full [&>.sr-only]:w-auto @md/field-group:[&>[data-slot=field-label]]:flex-auto @md/field-group:has-[>[data-slot=field-content]]:items-start @md/field-group:has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px")
		el.SetAttribute("data-orientation", "responsive")
	default:
		panic(fmt.Sprintf("invalid field orientation variant: %d", v))
	}
})

func New(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Role("group"),
			html.Data("slot", "field"),
			html.Class("group/field data-[invalid=true]:text-destructive flex w-full gap-3"),
			Vertical,
			orientation,
		)...,
	)
}

func Content(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "field-content"),
			html.Class("group/field-content flex flex-1 flex-col gap-1.5 leading-snug"),
		)...,
	)
}

func Label(items ...html.Item) html.Node {
	return html.Label(
		components.ItemsOf(
			items,
			html.Data("slot", "field-label"),
			html.Class("group/field-label peer/field-label flex w-fit gap-2 leading-snug group-data-[disabled=true]/field:opacity-50 has-[>[data-slot=field]]:w-full has-[>[data-slot=field]]:flex-col has-[>[data-slot=field]]:rounded-md has-[>[data-slot=field]]:border [&>[data-slot=field]]:p-4 has-data-[state=checked]:bg-primary/5 has-data-[state=checked]:border-primary dark:has-data-[state=checked]:bg-primary/10 text-sm font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"),
		)...,
	)
}

func Title(items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "field-label"),
			html.Class("flex w-fit items-center gap-2 text-sm font-medium leading-snug group-data-[disabled=true]/field:opacity-50"),
		)...,
	)
}

func Description(items ...html.Item) html.Node {
	return html.P(
		components.ItemsOf(
			items,
			html.Data("slot", "field-description"),
			html.Class("text-muted-foreground text-sm font-normal leading-normal group-has-[[data-orientation=horizontal]]/field:text-balance nth-last-2:-mt-1 last:mt-0 [[data-variant=legend]+&]:-mt-1.5 [&>a:hover]:text-primary [&>a]:underline [&>a]:underline-offset-4"),
		)...,
	)
}

func Separator(children html.Frag, items ...html.Item) html.Node {
	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "field-separator"),
			html.Class("relative -my-2 h-5 text-sm group-data-[variant=outline]/field-group:-mb-2"),

			separator.New(html.Class("absolute inset-0 top-1/2")),
			html.IfFn(children != nil, func() html.Item {
				return html.Span(html.Class("bg-background text-muted-foreground relative mx-auto block w-fit px-2"), html.Data("slot", "field-separator-content"), children)
			}),
		)...,
	)
}

func Error(errors []html.Node, children html.Frag, items ...html.Item) html.Node {
	if len(errors) == 0 {
		return html.Fragment()
	}

	var content = func() html.Node {
		if children != nil {
			return children
		}

		if len(errors) == 0 {
			return nil
		}

		if len(errors) == 1 {
			return errors[0]
		}

		return html.Ul(
			html.Class("ml-4 flex list-disc flex-col gap-1"),

			html.Map(errors, func(err html.Node) html.Node {
				return html.Li(err)
			}),
		)
	}()

	if content == nil {
		return html.Fragment()
	}

	return html.Div(
		components.ItemsOf(
			items,
			html.Data("slot", "field-error"),
			html.Class("text-destructive text-sm font-normal"),

			content,
		)...,
	)
}
