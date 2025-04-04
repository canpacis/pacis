package components

import (
	"image/color"
)

type Shade int

const (
	Shade50  = Shade(iota + 50)
	Shade100 = Shade(iota * 100)
	Shade200
	Shade300
	Shade400
	Shade500
	Shade600
	Shade700
	Shade800
	Shade900
	Shade950 = Shade(iota + 940)
)

type Palette [11]color.Color

func (p Palette) Color(shade Shade) color.Color {
	switch shade {
	case Shade50:
		return p[0]
	case Shade100:
		return p[1]
	case Shade200:
		return p[2]
	case Shade300:
		return p[3]
	case Shade400:
		return p[4]
	case Shade500:
		return p[5]
	case Shade600:
		return p[6]
	case Shade700:
		return p[7]
	case Shade800:
		return p[8]
	case Shade900:
		return p[9]
	case Shade950:
		return p[10]
	default:
		panic("invalid palette shade")
	}
}

var Slate = Palette{
	color.RGBA{248, 250, 252, 1},
	color.RGBA{241, 245, 249, 1},
	color.RGBA{226, 232, 240, 1},
	color.RGBA{203, 213, 225, 1},
	color.RGBA{148, 163, 184, 1},
	color.RGBA{100, 116, 139, 1},
	color.RGBA{71, 85, 105, 1},
	color.RGBA{51, 65, 85, 1},
	color.RGBA{30, 41, 59, 1},
	color.RGBA{15, 23, 42, 1},
	color.RGBA{2, 6, 23, 1},
}
