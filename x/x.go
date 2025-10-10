package x

import (
	"fmt"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/internal/util"
)

type DataProperty struct {
	id   string
	data any
}

func (*DataProperty) Item() {}

func (p *DataProperty) Apply(el *html.Element) {
	el.Children = append(el.Children, html.Script(html.Type("application/json"), html.ID(p.id), html.JSON(p.data)))
	el.SetAttribute("x-data", fmt.Sprintf("data('%s')", p.id))
}

func Data(data any) *DataProperty {
	return &DataProperty{data: data, id: util.PrefixedID("pacis")}
}

func Show(value string) *html.Attribute {
	return html.Attr("x-show", value)
}

func Bind(attr, value string) *html.Attribute {
	return html.Attr(":"+attr, value)
}

func Text(value string) *html.Attribute {
	return html.Attr("x-text", value)
}

func Model(value string) *html.Attribute {
	return html.Attr("x-model", value)
}

func Effect(value string) *html.Attribute {
	return html.Attr("x-effect", value)
}

var Ignore = html.Attr("x-ignore", "")

var Cloak = html.Attr("x-cloak", "")

func If(value string) *html.Attribute {
	return html.Attr("x-if", value)
}
