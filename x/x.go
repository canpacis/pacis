package x

import (
	"fmt"
	"strings"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/internal/util"
)

/*
Everything in Alpine starts with the x-data directive.

x-data defines a chunk of HTML as an Alpine component and provides the reactive data for that component to reference.

https://alpinejs.dev/directives/data
*/
type DataProperty struct {
	id   string
	data any
}

// Implements html.Property interface
func (*DataProperty) LifeCycle() html.PropertyLifeCycle {
	return html.LifeCycleStatic
}

// Implements html.Node interface
func (*DataProperty) Item() {}

// Implements html.Property interface
func (p *DataProperty) Apply(el *html.Element) {
	el.AppendNode(html.Script(html.Type("application/json"), html.ID(p.id), html.JSON(p.data)))
	el.SetAttribute("x-data", fmt.Sprintf("data('%s')", p.id))
}

/*
Everything in Alpine starts with the x-data directive.

x-data defines a chunk of HTML as an Alpine component and provides the reactive data for that component to reference.

https://alpinejs.dev/directives/data
*/
func Data(data any) *DataProperty {
	return &DataProperty{data: data, id: util.PrefixedID("pacis")}
}

/*
The x-init directive allows you to hook into the initialization phase of any element in Alpine.

https://alpinejs.dev/directives/init
*/
func Init(value string) *html.Attribute {
	return html.Attr("x-init", value)
}

/*
x-show is one of the most useful and powerful directives in Alpine. It provides an expressive way to show and hide DOM elements.

https://alpinejs.dev/directives/show
*/
func Show(value string) *html.Attribute {
	return html.Attr("x-show", value)
}

/*
x-bind allows you to set HTML attributes on elements based on the result of JavaScript expressions.

https://alpinejs.dev/directives/bind
*/
func Bind(attr, value string) *html.Attribute {
	return html.Attr(":"+attr, value)
}

// Alpine directive modifiers to customize the behavior of your event listeners.
type EventModifier = string

const (
	/*
		.prevent is the equivalent of calling .preventDefault() inside a listener on the browser event object.

		Usage:

			Form(
				x.On("submit", "console.log('submitted')", x.Prevent),
				Action("/foo"),

				Button(Text("Submit")),
			)

		In the above example, with the .prevent, clicking the button will NOT submit the form to the /foo
		endpoint. Instead, Alpine's listener will handle it and "prevent" the event from being handled any further.
	*/
	Prevent EventModifier = "prevent"

	/*
		Similar to .prevent, .stop is the equivalent of calling .stopPropagation() inside a listener on the browser event object.

		Usage:

			Div(
				x.On("click", "console.log('I will not get logged')"),

				Button(
					x.On("click", "", x.Stop),

					Text("Click Me"),
				),
			)

		In the above example, clicking the button WON'T log the message. This is because we are stopping the propagation of the event
		immediately and not allowing it to "bubble" up to the <div> with the @click listener on it.
	*/
	Stop EventModifier = "stop"

	/*
		.outside is a convenience helper for listening for a click outside of the element it is attached to. Here's a simple dropdown component example to demonstrate:

		Usage:

			Div(
				x.Data(map[string]any{"open": false}),

				Button(
					x.On("click", "open = !open")

					Text("Toggle"),
				),
				Div(
					x.Show("open"),
					x.On("click", "open = false", x.Outside),

					Text("Contents..."),
				),
			)

		In the above example, after showing the dropdown contents by clicking the "Toggle" button, you can close the dropdown by clicking anywhere on
		the page outside the content. This is because .outside is listening for clicks that DON'T originate from the element it's registered on.
		It's worth noting that the .outside expression will only be evaluated when the element it's registered on is visible on the page. Otherwise,
		there would be nasty race conditions where clicking the "Toggle" button would also fire the @click.outside handler when it is not visible.
	*/
	Outside EventModifier = "outside"

	/*
		When the .window modifier is present, Alpine will register the event listener on the root window object on the page instead of the element itself.

		Usage:

			Div(x.On("keyup", "", "escape", x.Window))

		The above snippet will listen for the "escape" key to be pressed ANYWHERE on the page. Adding .window to listeners is extremely useful for these
		sorts of cases where a small part of your markup is concerned with events that take place on the entire page.
	*/
	Window EventModifier = "window"
	/*
	  .document works similarly to .window only it registers listeners on the document global, instead of the window global.
	*/
	Document EventModifier = "document"
	/*
		By adding .once to a listener, you are ensuring that the handler is only called ONCE.

		Usage:

			Button(x.On("click", "console.log('I will only log once')", x.Once))
	*/
	Once EventModifier = "once"
	/*
		Sometimes it is useful to "debounce" an event handler so that it only is called after a certain period of inactivity (250 milliseconds by default).
		For example if you have a search field that fires network requests as the user types into it, adding a debounce will prevent the network requests
		from firing on every single keystroke.

		Usage:

			Input(x.On("input", "fetchResults", x.Debounce))

		Now, instead of calling fetchResults after every keystroke, fetchResults will only be called after 250 milliseconds of no keystrokes.
		If you wish to lengthen or shorten the debounce time, you can do so by trailing a duration after the .debounce modifier like so:

		Usage:

			Input(x.On("input", "fetchResults", x.Debounce, "500ms"))

		Now, fetchResults will only be called after 500 milliseconds of inactivity.
	*/
	Debounce EventModifier = "debounce"
	/*
		.throttle is similar to .debounce except it will release a handler call every 250 milliseconds instead of deferring it indefinitely. This is
		useful for cases where there may be repeated and prolonged event firing and using .debounce won't work because you want to still handle
		the event every so often.

		Usage:

			Div(x.On("scroll", "handleScroll", x.Window, x.Throttle))

		The above example is a great use case of throttling. Without .throttle, the handleScroll method would be fired hundreds of times as the user
		scrolls down a page. This can really slow down a site. By adding .throttle, we are ensuring that handleScroll only gets called every 250
		milliseconds.

		Just like with .debounce, you can add a custom duration to your throttled event:

		Usage:

			Div(x.On("scroll", "handleScroll", x.Window, x.Throttle, "750mx"))

		Now, handleScroll will only be called every 750 milliseconds.
	*/
	Throttle EventModifier = "throttle"
	/*
		By adding .self to an event listener, you are ensuring that the event originated on the element it is declared on, and not from a child element.

		Usage:

			Button(
				x.On("click", "handleClick", x.Self),

				Text("Click Me"),
				Img(Src("...")),
			)

		In the above example, we have an <img> tag inside the <button> tag. Normally, any click originating within the <button> element
		(like on <img> for example), would be picked up by a @click listener on the button. However, in this case, because we've added a .self, only
		clicking the button itself will call handleClick. Only clicks originating on the <img> element will not be handled.
	*/
	Self EventModifier = "self"
	/*
		Usage:

			Div(x.On("custom-event", "handleCustomEvent", x.Camel))

		Sometimes you may want to listen for camelCased events such as customEvent in our example. Because camelCasing inside HTML attributes
		is not supported, adding the .camel modifier is necessary for Alpine to camelCase the event name internally. By adding .camel in the
		above example, Alpine is now listening for customEvent instead of custom-event.

	*/
	Camel EventModifier = "camel"
	/*
		Usage:

			Div(x.On("custom-event", "handleCustomEvent", x.Dot))

		Similar to the .camelCase modifier there may be situations where you want to listen for events that have dots in their name
		(like custom.event). Since dots within the event name are reserved by Alpine you need to write them with dashes and add the .dot modifier.
		In the code example above custom-event.dot will correspond to the event name custom.event.
	*/
	Dot EventModifier = "dot"
	/*
		Browsers optimize scrolling on pages to be fast and smooth even when JavaScript is being executed on the page. However, improperly
		implemented touch and wheel listeners can block this optimization and cause poor site performance. If you are listening for touch
		events, it's important to add .passive to your listeners to not block scroll performance.

		Usage:

			Div(x.On("touchstart", "...", x.Passive))

		https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#improving_scrolling_performance_with_passive_listeners
	*/
	Passive EventModifier = "passive"

	/*
		Add this modifier if you want to execute this listener in the event's capturing phase, e.g. before the event bubbles from the target element up the DOM.

		Usage:

			Div(
				x.On("click", "console.log('I will log first')", x.Capture)

				Button(x.On("click", "console.log('I will log second')")),
			)
	*/
	Capture EventModifier = "capture"
)

/*
x-on allows you to easily run code on dispatched DOM events.

https://alpinejs.dev/directives/on
*/
func On(event, handler string, modifiers ...EventModifier) *html.Attribute {
	var mods string
	if len(modifiers) > 0 {
		mods = "." + strings.Join(modifiers, ".")
	}
	return html.Attr(fmt.Sprintf("@%s%s", event, mods), handler)
}

/*
x-text sets the text content of an element to the result of a given expression.

https://alpinejs.dev/directives/text
*/
func Text(value string) *html.Attribute {
	return html.Attr("x-text", value)
}

/*
x-model allows you to bind the value of an input element to Alpine data.

https://alpinejs.dev/directives/model
*/
func Model(value string) *html.Attribute {
	return html.Attr("x-model", value)
}

/*
x-modelable allows you to expose any Alpine property as the target of the x-model directive.

https://alpinejs.dev/directives/modelable
*/
func Modelable(value string) *html.Attribute {
	return html.Attr("x-modelable", value)
}

/*
Alpine's x-for directive allows you to create DOM elements by iterating through a list.

https://alpinejs.dev/directives/for
*/
func For(value string) *html.Attribute {
	return html.Attr("x-for", value)
}

/*
x-effect is a useful directive for re-evaluating an expression when one of its dependencies change. You can think of it as a
watcher where you don't have to specify what property to watch, it will watch all properties used within it.

https://alpinejs.dev/directives/effect
*/
func Effect(value string) *html.Attribute {
	return html.Attr("x-effect", value)
}

/*
By default, Alpine will crawl and initialize the entire DOM tree of an element containing x-init or x-data.

https://alpinejs.dev/directives/ignore
*/
var Ignore = html.Attr("x-ignore", "")

/*
x-ref in combination with $refs is a useful utility for easily accessing DOM elements directly. It's most useful as a
replacement for APIs like getElementById and querySelector.

https://alpinejs.dev/directives/ref
*/
func Ref(value string) *html.Attribute {
	return html.Attr("x-ref", value)
}

/*
Sometimes, when you're using AlpineJS for a part of your template, there is a "blip" where you might see your
uninitialized template after the page loads, but before Alpine loads.

x-cloak addresses this scenario by hiding the element it's attached to until Alpine is fully loaded on the page.

https://alpinejs.dev/directives/cloak
*/
var Cloak = html.Attr("x-cloak", "")

/*
The x-teleport directive allows you to transport part of your Alpine template to another part of the DOM on
the page entirely. This is useful for things like modals (especially nesting them), where it's helpful to
break out of the z-index of the current Alpine component.

https://alpinejs.dev/directives/teleport
*/
func Teleport(value string) *html.Attribute {
	return html.Attr("x-teleport", value)
}

/*
x-if is used for toggling elements on the page, similarly to x-show, however it completely adds and removes
the element it's applied to rather than just changing its CSS display property to "none". Because of this
difference in behavior, x-if should not be applied directly to the element, but instead to a <template> tag
that encloses the element. This way, Alpine can keep a record of the element once it's removed from the page.

https://alpinejs.dev/directives/if
*/
func If(value string) *html.Attribute {
	return html.Attr("x-if", value)
}

/*
x-id allows you to declare a new "scope" for any new IDs generated using $id(). It accepts an array of strings
(ID names) and adds a suffix to each $id('...') generated within it that is unique to other IDs on the page.
x-id is meant to be used in conjunction with the $id(...) magic.
*/
func ID(value string) *html.Attribute {
	return html.Attr("x-id", value)
}
