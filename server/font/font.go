// Package font provides utilities for managing and generating Google Fonts URLs,
// including font weights, subsets, and display options.
//
// It defines types for font weights, subsets, and display strategies, and offers
// functions to construct font URLs and HTML head nodes for font inclusion.
package font

import (
	"fmt"
	"strings"

	"github.com/canpacis/pacis/html"
)

/*
“Weight” (wght in CSS) is an axis found in many variable fonts. It controls the font file’s weight parameter.

https://fonts.google.com/knowledge/glossary/weight_axis
*/
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

// Implements the Stringer interface
func (w Weight) String() string {
	return fmt.Sprintf("%d", w)
}

// A set of font weights to load
type WeightList []Weight

// Implements the Stringer interface
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

/*
Represents a font subset e.g: latin or cyrillic.

Subsetting is the practice of creating a “subset” of a font—a file that contains a
custom (and usually limited) collection of glyphs.

https://fonts.google.com/knowledge/glossary/subsetting
*/
type Subset string

const (
	/*
		In typography, “Latin script” refers to the most widely adopted writing system in
		the world (and is also often used as a term to mean placeholder copy).

		https://fonts.google.com/knowledge/glossary/latin
	*/
	Latin = Subset("latin")
	/*
		In typography, “Latin script” refers to the most widely adopted writing system in
		the world (and is also often used as a term to mean placeholder copy).

		https://fonts.google.com/knowledge/glossary/latin
	*/
	LatinExt = Subset("latin-ext")
	/*
		Cyrillic is a writing system, named after the missionary work of St. Cyril in the
		first Bulgarian Empire. The original Cyrillic script was based on uppercase Greek letterforms.

		https://fonts.google.com/knowledge/glossary/cyrillic
	*/
	Cyrillic = Subset("cyrillic")
)

/*
The font-display property lets you control what happens while the font is still loading
or otherwise unavailable. Specifying a value other than the default auto is usually appropriate.

https://developers.google.com/fonts/docs/css2#optimizing_for_latency_and_file_size
*/
type Display string

const (
	/*
		The font display policy is user-agent-defined.

		https://www.w3.org/TR/css-fonts-4/#valdef-font-face-font-display-auto
	*/
	Auto = Display("auto")
	/*
		Gives the font face a short block period (3s is recommended in most cases) and an infinite swap period.

		https://www.w3.org/TR/css-fonts-4/#valdef-font-face-font-display-block
	*/
	Block = Display("block")
	/*
		Gives the font face an extremely small block period (100ms or less is recommended in
		most cases) and an infinite swap period.

		https://www.w3.org/TR/css-fonts-4/#valdef-font-face-font-display-swap
	*/
	Swap = Display("swap")
	/*
		Gives the font face an extremely small block period (100ms or less is recommended in
		most cases) and a short swap period (3s is recommended in most cases).

		https://www.w3.org/TR/css-fonts-4/#valdef-font-face-font-display-fallback
	*/
	Fallback = Display("fallback")
	/*
		If the font can be loaded "immediately" (such that it’s available to be used for
		the "first paint" of the text), the font is used.
		Otherwise, the font is treated as if its block period and swap period both expired
		before it finished loading. If the font is not used due to this, the user agent may
		choose to abort the font download, or download it with a very low priority. If the
		user agent believes it would be useful for the user, it may avoid even starting the
		font download, and proceed immediately to using a fallback font.

		An optional font must never cause the layout of the page to "jump" as it loads in.
		A user agent may choose to slightly delay rendering an element using an optional
		font to give it time to load from a possibly-slow local cache, but once the text
		has been painted to the screen with a fallback font instead, it must not be rendered
		with the optional font for the rest of the page’s lifetime.

		This value should be used for body text, or any other text where the chosen font is
		purely a decorative "nice-to-have". It should be used anytime it is more important
		that the web page render quickly on first visit, than it is that the user wait a
		longer time to see everything perfect immediately.

		https://www.w3.org/TR/css-fonts-4/#valdef-font-face-font-display-optional
	*/
	Optional = Display("optional")
)

/*
A typeface is the underlying visual design that can exist in many different typesetting
technologies and a font is one of these implementations. In other words, a typeface is
what you see and a font is what you use.

Decades of choosing type via software “font” menus has caused the majority of us to
think of “font” and “typeface” as interchangeable terms for the same thing, but the
difference is important for anyone serious about typography. Another useful analogy
is that a typeface is to a song as a font is to an MP3 file: It’s a manifestation of
the typeface/song, but the typeface/song exists outside of the format.

https://fonts.google.com/knowledge/glossary/font
*/
type Font struct {
	Name       string
	WeightList WeightList
	Subsets    []Subset
	Display    Display
	// TODO: Optical sizing and default weight list features
}

// Builds the URL for loading the font from Google Fonts
func (f Font) URL() string {
	return fmt.Sprintf(
		"https://fonts.googleapis.com/css2?family=%s%s&display=%s",
		strings.ReplaceAll(f.Name, " ", "+"),
		f.WeightList.String(),
		f.Display,
	)
}

/*
Creates a new Font struct with given name, weight, display and subset properties

https://fonts.google.com/
*/
func New(name string, weights WeightList, display Display, subsets ...Subset) *Font {
	return &Font{
		Name:       name,
		WeightList: weights,
		Subsets:    subsets,
		Display:    display,
	}
}

// Head generates an HTML fragment containing <link> elements for preconnecting to Google Fonts
// and for including the provided font stylesheets. It accepts a variadic list of Font pointers,
// and for each font, it creates a corresponding <link rel="stylesheet"> element referencing the font's URL.
func Head(fonts ...*Font) html.Node {
	return html.Fragment(
		html.Link(html.Href("https://fonts.googleapis.com"), html.Rel("preconnect")),
		html.Link(html.Href("https://fonts.gstatic.com"), html.Rel("preconnect")),
		html.Map(fonts, func(font *Font) html.Node {
			return html.Link(html.Href(font.URL()), html.Rel("stylesheet"))
		}),
	)
}
