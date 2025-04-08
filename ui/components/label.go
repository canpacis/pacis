package components

import h "github.com/canpacis/pacis/ui/html"

func Label(text string, props ...h.I) h.Element {
	props = Join(props, h.Text(text), h.Class("text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 block [&>input]:mt-2"))
	return h.Lbl(props...)
}
