package components

import (
	"log"
	"maps"
	"strings"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/canpacis/pacis/html"
)

type TailwindMergeProperty struct{}

func (*TailwindMergeProperty) Apply(el *html.Element) {
	el.ClassList.Items = strings.Split(twmerge.Merge(strings.Join(el.ClassList.Items, " ")), " ")
}
func (*TailwindMergeProperty) Item() {}

var TailwindMerge = &TailwindMergeProperty{}

func ItemsOf(passed []html.Item, items ...html.Item) []html.Item {
	i := []html.Item{}
	i = append(i, items...)
	i = append(i, passed...)
	i = append(i, TailwindMerge)
	return i
}

type AsChildProperty struct{}

func (*AsChildProperty) Apply(el *html.Element) {
	if len(el.Children) != 1 {
		log.Fatal("Exactly 1 child should be present in an element with AsChild property")
	}
	child, ok := el.Children[0].(*html.Element)
	if !ok {
		log.Fatal("Non element node is passed to the element with AsChild propery")
	}

	child.ClassList.Items = append(child.ClassList.Items, el.ClassList.Items...)
	maps.Copy(child.Attributes, el.Attributes)
	*el = *child
}

func (*AsChildProperty) Item() {}

var AsChild = &AsChildProperty{}
