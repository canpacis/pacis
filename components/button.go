package components

import (
	"io"

	"github.com/canpacis/pacis/renderer"
)

type ButtonSize int

const (
	ButtonSizeDefault = ButtonSize(iota)
	ButtonSizeSm
	ButtonSizeLg
	ButtonSizeIcon
)

func (s ButtonSize) Render(w io.Writer) error {
	class := ""

	switch s {
	case ButtonSizeDefault:
		class = "default-classes"
	case ButtonSizeSm:
		class = "sm-classes"
	case ButtonSizeLg:
		class = "lg-classes"
	case ButtonSizeIcon:
		class = "icon-classes"
	default:
		panic("invalid button size property")
	}

	return renderer.Class(class).Render(w)
}

func (ButtonSize) Key() string {
	return "class"
}

func Button(constituents ...renderer.Renderer) *renderer.Element {
	var size ButtonSize

	c := []renderer.Renderer{
		renderer.Class("base-classes"),
	}

	for _, cn := range constituents {
		switch cn := cn.(type) {
		case ButtonSize:
			size = cn
		default:
			c = append(c, cn)
		}
	}

	c = append(c, size)

	return renderer.Btn(c...)
}
