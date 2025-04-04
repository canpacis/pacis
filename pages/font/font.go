package fonts

import (
	"fmt"
	"strings"

	h "github.com/canpacis/pacis/ui/html"
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

	return fmt.Sprintf("%d..%d", min, max)
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
		"https://fonts.googleapis.com/css2?family=%s:wght@%s&display=%s",
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

func Head(fonts ...*Font) *h.Fragment {
	return h.Frag(
		h.Link(h.Href("https://fonts.googleapis.com"), h.Rel("preconnect")),
		h.Link(h.Href("https://fonts.gstatic.com"), h.Rel("preconnect")),
		h.Map(fonts, func(font *Font, i int) h.Node {
			return h.Link(h.Href(font.URL()), h.Rel("stylesheet"))
		}),
	)
}
