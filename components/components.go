// Package components provides utilities for manipulating HTML elements with custom properties,
// such as merging Tailwind CSS classes and handling child elements.
//
// The ItemsOf function combines multiple html.Item values, appending the TailwindMerge property
// to the resulting slice for consistent class merging behavior.
package components

import (
	"log"
	"maps"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/canpacis/pacis/html"
)

type Variant int

func (Variant) Item() {}

func (Variant) LifeCycle() html.PropertyLifeCycle {
	return html.LifeCycleImmediate
}

func (v Variant) Apply(el *html.Element) {
	el.Set("variant", v)
}

type VariantApplier struct {
	fn func(*html.Element, Variant)
}

func (*VariantApplier) Item() {}

func (*VariantApplier) LifeCycle() html.PropertyLifeCycle {
	return html.LifeCycleStatic
}

func (va *VariantApplier) Apply(el *html.Element) {
	va.fn(el, el.Get("variant").(Variant))
}

func NewVariantApplier(fn func(*html.Element, Variant)) *VariantApplier {
	return &VariantApplier{fn: fn}
}

type Size int

func (Size) Item() {}

func (Size) LifeCycle() html.PropertyLifeCycle {
	return html.LifeCycleImmediate
}

func (v Size) Apply(el *html.Element) {
	el.Set("size", v)
}

type SizeApplier struct {
	fn func(*html.Element, Size)
}

func (*SizeApplier) Item() {}

func (*SizeApplier) LifeCycle() html.PropertyLifeCycle {
	return html.LifeCycleStatic
}

func (va *SizeApplier) Apply(el *html.Element) {
	va.fn(el, el.Get("size").(Size))
}

func NewSizeApplier(fn func(*html.Element, Size)) *SizeApplier {
	return &SizeApplier{fn: fn}
}

/*
The TailwindMergeProperty type applies Tailwind CSS class merging to an element's class list,
ensuring that conflicting or duplicate classes are resolved according to Tailwind's rules.
*/
type TailwindMergeProperty struct{}

func (*TailwindMergeProperty) LifeCycle() html.PropertyLifeCycle {
	return html.LifeCycleImmediate
}

// Implements the html.Item interface
func (*TailwindMergeProperty) Item() {}

// Implements the html.Property interface
func (*TailwindMergeProperty) Apply(el *html.Element) {
	el.SetAttribute("class", twmerge.Merge(el.GetAttribute("class")))
}

// TailwindMerge is a global instance of TailwindMergeProperty used to merge
// conflicting or duplicate tailwind classes.
var TailwindMerge = &TailwindMergeProperty{}

/*
The AsChildProperty type allows an element to adopt the properties of its single child element,
merging class lists and attributes, and replacing the parent element with the child.
*/
type AsChildProperty struct{}

// Implements the html.Item interface
func (*AsChildProperty) Item() {}

// Implements the html.Hook interface
func (*AsChildProperty) Apply(el *html.Element) {
	nodes := el.GetNodes()
	if len(nodes) != 1 {
		log.Fatal("Exactly 1 child should be present in an element with AsChild property")
	}
	child, ok := nodes[0].(*html.Element)
	if !ok {
		log.Fatal("Non element node is passed to the element with AsChild propery")
	}

	// child.ClassList.Items = append(child.ClassList.Items, el.ClassList.Items...)
	attrs := map[string]string{}
	maps.Copy(attrs, child.GetAttributes())
	maps.Copy(attrs, el.GetAttributes())
	// TODO: Deferred properties are not copied this way
	child.SetAttributes(attrs)
	*el = *child
}

// AsChild is a global instance of AsChildHook used to provide hook functionality
// for child components within the package. It can be used to manage or modify
// behavior specific to child components.
var AsChild = &AsChildProperty{}

// Merges a list of html.Items with another. Puts the first list of items at the end.
func ItemsOf(passed []html.Item, items ...html.Item) []html.Item {
	itms := make([]html.Item, len(passed)+len(items)+1)
	copy(itms[:len(items)], items)
	copy(itms[len(items):], passed)
	itms[len(itms)-1] = TailwindMerge
	return itms
}
