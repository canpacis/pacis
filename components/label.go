package components

import r "github.com/canpacis/pacis-ui/renderer"

func Label(props ...r.I) r.Element {
	props = Join(props, r.Class("text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 block"))
	return r.Lbl(props...)
}
