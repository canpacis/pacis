package components

import r "github.com/canpacis/pacis/renderer"

func Label(props ...r.I) r.Element {
	props = join(props, r.Class("text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"))
	return r.Lbl(props...)
}
