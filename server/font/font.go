package font

import (
	"fmt"
	"strings"

	"github.com/canpacis/pacis/html"
)

type Weight int

const (
	_    = Weight(iota)
	W100 = Weight(iota * 100)
	W200
	W300
	W400
	W500
	W600
	W700
	W800
	W900
)

func (w Weight) String() string {
	return fmt.Sprintf("%d", w)
}

type WeightList []Weight

func (wl WeightList) String() string {
	min := 1000
	max := 0

	for _, weight := range wl {
		if weight < Weight(min) {
			min = int(weight)
		}
		if weight > Weight(max) {
			max = int(weight)
		}
	}

	if max <= min {
		return ""
	}

	return fmt.Sprintf(":wght@%d..%d", min, max)
}

type Subset string

const (
	Latin    = Subset("latin")
	LatinExt = Subset("latin-ext")
	Cyrillic = Subset("cyrillic")
)

type Display string

const (
	Swap     = Display("swap")
	Auto     = Display("auto")
	Block    = Display("block")
	Fallback = Display("fallback")
	Optional = Display("optional")
)

type Font struct {
	Name       string
	WeightList WeightList
	Subsets    []Subset
	Display    Display
	// TODO: Optical sizing and default weight list features
}

func (f Font) URL() string {
	return fmt.Sprintf(
		"https://fonts.googleapis.com/css2?family=%s%s&display=%s",
		strings.ReplaceAll(f.Name, " ", "+"),
		f.WeightList.String(),
		f.Display,
	)
}

func New(name string, weights WeightList, display Display, subsets ...Subset) *Font {
	return &Font{
		Name:       name,
		WeightList: weights,
		Subsets:    subsets,
		Display:    display,
	}
}

func Head(fonts ...*Font) html.Node {
	return html.Fragment(
		html.Link(html.Href("https://fonts.googleapis.com"), html.Rel("preconnect")),
		html.Link(html.Href("https://fonts.gstatic.com"), html.Rel("preconnect")),
		html.Map(fonts, func(font *Font) html.Node {
			return html.Link(html.Href(font.URL()), html.Rel("stylesheet"))
		}),
	)
}
