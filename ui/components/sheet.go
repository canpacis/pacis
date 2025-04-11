package components

import (
	h "github.com/canpacis/pacis/ui/html"
)

var (
	SheetVariantLeft = &GroupedClass{
		"sheet-variant",
		"top-0 bottom-0 left-0 h-dvh w-3/4 border-r sm:max-w-sm",
		true,
	}
	SheetVariantTop = &GroupedClass{
		"sheet-variant",
		"left-0 right-0 top-0 border-b",
		false,
	}
	SheetVariantRight = &GroupedClass{
		"sheet-variant",
		"top-0 bottom-0 right-0 h-dvh w-3/4 border-l sm:max-w-sm",
		false,
	}
	SheetVariantBottom = &GroupedClass{
		"sheet-variant",
		"left-0 right-0 bottom-0 border-t",
		false,
	}
)

func Sheet(props ...h.I) h.Element {
	el := h.Div(props...)

	_, isopen := el.GetAttribute("open")
	if isopen {
		el.RemoveAttribute("open")
	}
	id := getid(el)
	el.AddAttribute(X("data", fn("sheet", isopen, id)))
	return el
}

func SheetTrigger(trigger h.Element) h.Element {
	trigger.AddAttribute(OpenSheet)
	return trigger
}

func SheetContent(props ...h.I) h.Node {
	props = Join(
		props,
		SheetVariantLeft,
		X("cloak"),
		X("show", "open"),
		X("trap.noscroll", "open"),
		CloseSheetOn("click.outside"),
		h.Data(":state", "open ? 'opened' : 'closed'"),
		h.Class("fixed z-50 gap-4 bg-background p-6 shadow-lg transition ease-in-out"),
	)

	content := h.Div(props...)

	for _, prop := range props {
		grouped, ok := prop.(*groupedclasses)
		if ok {
			variant := grouped.Candidate()

			switch variant {
			case SheetVariantLeft:
				content.AddAttribute(X("transition:enter", "-translate-x-[100%]"))
				content.AddAttribute(X("transition:leave", "-translate-x-[100%]"))
			case SheetVariantTop:
				content.AddAttribute(X("transition:enter", "-translate-y-[100%]"))
				content.AddAttribute(X("transition:leave", "-translate-y-[100%]"))
			case SheetVariantRight:
				content.AddAttribute(X("transition:enter", "translate-x-[100%]"))
				content.AddAttribute(X("transition:leave", "translate-x-[100%]"))
			case SheetVariantBottom:
				content.AddAttribute(X("transition:enter", "translate-y-[100%]"))
				content.AddAttribute(X("transition:leave", "translate-y-[100%]"))
			}
		}
	}

	return h.Frag(
		// Overlay
		h.Div(
			X("cloak"),
			X("show", "open"),
			h.Data(":state", "open ? 'opened' : 'closed'"),
			h.Class("fixed h-dvh inset-0 z-50 bg-black/80"),
			X("transition.opacity"),
		),
		content,
	)
}

func OpenSheetOn(event string) h.Attribute {
	return On(event, "openSheet()")
}

var OpenSheet = OpenSheetOn("click")

func CloseSheetOn(event string) h.Attribute {
	return On(event, "closeSheet()")
}

var CloseSheet = CloseSheetOn("click")
