package icons

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	r "github.com/canpacis/pacis/renderer"
)

type IconOptions struct {
	Width       float64
	Height      float64
	Fill        string
	Stroke      string
	StrokeWidth float64
}

func Width(w float64) func(*IconOptions) {
	return func(io *IconOptions) {
		io.Width = w
	}
}

func Height(h float64) func(*IconOptions) {
	return func(io *IconOptions) {
		io.Height = h
	}
}

func Fill(fill string) func(*IconOptions) {
	return func(io *IconOptions) {
		io.Fill = fill
	}
}

func Stroke(stroke string) func(*IconOptions) {
	return func(io *IconOptions) {
		io.Stroke = stroke
	}
}

func StrokeWidth(w float64) func(*IconOptions) {
	return func(io *IconOptions) {
		io.StrokeWidth = w
	}
}

type Node struct {
	root    bool
	options []func(*IconOptions)
	name    string

	XMLName xml.Name

	CX             string `xml:"cx,attr"`
	CY             string `xml:"cy,attr"`
	X1             string `xml:"x1,attr"`
	X2             string `xml:"x2,attr"`
	Y1             string `xml:"y1,attr"`
	Y2             string `xml:"y2,attr"`
	R              string `xml:"r,attr"`
	D              string `xml:"d,attr"`
	Width          string `xml:"width,attr"`
	Height         string `xml:"height,attr"`
	Viewbox        string `xml:"viewBox,attr"`
	Fill           string `xml:"fill,attr"`
	Stroke         string `xml:"stroke,attr"`
	Points         string `xml:"points,attr"`
	StrokeWidth    string `xml:"stroke-width,attr"`
	StrokeLinecap  string `xml:"stroke-linecap,attr"`
	StrokeLinejoin string `xml:"stroke-linejoin,attr"`

	Nodes []Node `xml:",any"`
}

func (el *Node) Render(w io.Writer) error {
	if el.root {
		if err := el.init(); err != nil {
			return err
		}
	}
	attrs := []r.Renderer{}

	switch el.XMLName.Local {
	case "svg":
		attrs = append(attrs,
			r.Attr("width", el.Width),
			r.Attr("height", el.Height),
			r.Attr("viewBox", el.Viewbox),
			r.Attr("fill", el.Fill),
			r.Attr("stroke", el.Stroke),
			r.Attr("stroke-width", el.StrokeWidth),
			r.Attr("stroke-linecap", el.StrokeLinecap),
			r.Attr("stroke-linejoin", el.StrokeLinejoin),
		)
	case "circle":
		attrs = append(attrs,
			r.Attr("cx", el.CX),
			r.Attr("cy", el.CY),
			r.Attr("r", el.R),
			r.Attr("fill", el.Fill),
			r.Attr("stroke", el.Stroke),
		)
	case "path":
		attrs = append(attrs,
			r.Attr("d", el.D),
			r.Attr("fill", el.Fill),
			r.Attr("stroke", el.Stroke),
		)
	case "polyline":
		attrs = append(attrs,
			r.Attr("points", el.Points),
			r.Attr("fill", el.Fill),
			r.Attr("stroke", el.Stroke),
		)
	case "line":
		attrs = append(attrs,
			r.Attr("x1", el.X1),
			r.Attr("x2", el.X2),
			r.Attr("y1", el.Y1),
			r.Attr("y2", el.Y2),
			r.Attr("fill", el.Fill),
			r.Attr("stroke", el.Stroke),
		)
	case "rect", "group":
	default:
		return fmt.Errorf("invalid xml element %s", el.XMLName.Local)
	}

	for _, node := range el.Nodes {
		attrs = append(attrs, &node)
	}

	return r.El(el.XMLName.Local, attrs...).Render(w)
}

func (n *Node) init() error {
	file, err := os.Open(fmt.Sprintf("./lucide/icons/%s.svg", n.name))
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	var svg Node
	if err := decoder.Decode(&svg); err != nil {
		return err
	}

	opts := IconOptions{
		Width:       24,
		Height:      24,
		Fill:        "none",
		Stroke:      "currentColor",
		StrokeWidth: 2,
	}

	for _, fn := range n.options {
		fn(&opts)
	}

	*n = Node{
		XMLName:        svg.XMLName,
		Viewbox:        svg.Viewbox,
		Width:          fmt.Sprintf("%f", opts.Width),
		Height:         fmt.Sprintf("%f", opts.Height),
		Fill:           opts.Fill,
		Stroke:         opts.Stroke,
		StrokeWidth:    fmt.Sprintf("%f", opts.StrokeWidth),
		StrokeLinecap:  svg.StrokeLinecap,
		StrokeLinejoin: svg.StrokeLinejoin,
		Nodes:          svg.Nodes,
	}
	return nil
}

func (*Node) Node() {}

func Icon(name string, options ...func(*IconOptions)) *Node {
	return &Node{name: name, options: options, root: true}
}
