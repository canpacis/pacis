package components

import . "github.com/canpacis/pacis/ui/html"

func Knot(props ...I) Element {
	return A(Join(props, Target("blank"), Rel("noreferer"))...)
}
