// Package components provides utilities for manipulating HTML elements with custom properties,
// such as merging Tailwind CSS classes and handling child elements.
//
// The ItemsOf function combines multiple html.Item values, appending the TailwindMerge property
// to the resulting slice for consistent class merging behavior.
package components

import (
	"log"
	"maps"
	"strings"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/canpacis/pacis/html"
)

/*
The TailwindMergeProperty type applies Tailwind CSS class merging to an element's class list,
ensuring that conflicting or duplicate classes are resolved according to Tailwind's rules.
*/
type TailwindMergeProperty struct{}

// Implements the html.Property interface
func (*TailwindMergeProperty) Apply(el *html.Element) {
	el.ClassList.Items = strings.Split(twmerge.Merge(strings.Join(el.ClassList.Items, " ")), " ")
}

// Implements the html.Item interface
func (*TailwindMergeProperty) Item() {}

// TailwindMerge is a global instance of TailwindMergeProperty used to merge
// conflicting or duplicate tailwind classes.
var TailwindMerge = &TailwindMergeProperty{}

/*
The AsChildHook type allows an element to adopt the properties of its single child element,
merging class lists and attributes, and replacing the parent element with the child.
*/
type AsChildHook struct{}

// Implements the html.Item interface
func (*AsChildHook) Item() {}

// Implements the html.Hook interface
func (*AsChildHook) Hook(el *html.Element) {
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

// AsChild is a global instance of AsChildHook used to provide hook functionality
// for child components within the package. It can be used to manage or modify
// behavior specific to child components.
var AsChild = &AsChildHook{}

// Merges a list of html.Items with another. Puts the first list of items at the end.
func ItemsOf(passed []html.Item, items ...html.Item) []html.Item {
	i := []html.Item{}
	i = append(i, items...)
	i = append(i, passed...)
	i = append(i, TailwindMerge)
	return i
}
