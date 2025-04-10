package components

import h "github.com/canpacis/pacis/ui/html"

func Knot(props ...h.I) h.Element {
	return h.A(Join(props, h.Target("blank"), h.Rel("noreferer"))...)
}
