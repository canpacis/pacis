package metadata

import (
	"fmt"
	"maps"
	"net/url"
	"strings"

	"github.com/canpacis/pacis/html"
)

type Author struct {
	Name string
	URL  string
}

type Alternate struct {
	Title    string
	Href     string
	HrefLang string
	Type     string
	Media    string
}

type Robots struct {
	Index   bool
	Follow  bool
	Nocache bool
	Other   map[string]map[string]any
}

func (r *Robots) Node() html.Node {
	var build = func(data map[string]any) string {
		list := []string{}
		for key, value := range data {
			boolean, ok := value.(bool)
			if ok && boolean {
				list = append(list, key)
				continue
			}
			list = append(list, fmt.Sprintf("%s:%s", key, value))
		}
		return strings.Join(list, ", ")
	}

	return html.Fragment(
		html.Meta(html.Name("robots"), html.Content(build(map[string]any{"index": r.Index, "follow": r.Follow, "nocache": r.Nocache}))),
		html.Iter2(maps.All(r.Other), func(key string, value map[string]any) html.Node {
			return html.Meta(html.Name(key), html.Content(build(value)))
		}),
	)
}

type OpenGraphMedia struct {
	URL    string
	Width  int
	Height int
	Alt    string
}

type OpenGraph struct {
	Title       string
	Description string
	URL         string
	SiteName    string
	Locale      string
	Type        string
	Images      []OpenGraphMedia
	Videos      []OpenGraphMedia
	Audio       []OpenGraphMedia
}

func (o *OpenGraph) Node() html.Node {
	return html.Fragment(
		html.If(len(o.Title) > 0, html.Meta(html.PropertyAttr("og:title"), html.Content(o.Title))).(html.Node),
		html.If(len(o.Description) > 0, html.Meta(html.PropertyAttr("og:description"), html.Content(o.Description))).(html.Node),
		html.If(len(o.URL) > 0, html.Meta(html.PropertyAttr("og:url"), html.Content(o.URL))).(html.Node),
		html.If(len(o.SiteName) > 0, html.Meta(html.PropertyAttr("og:site_name"), html.Content(o.SiteName))).(html.Node),
		html.If(len(o.Locale) > 0, html.Meta(html.PropertyAttr("og:locale"), html.Content(o.Locale))).(html.Node),
		html.If(len(o.Type) > 0, html.Meta(html.PropertyAttr("og:type"), html.Content(o.Type))).(html.Node),
		html.Map(o.Images, func(image OpenGraphMedia) html.Node {
			return html.Fragment(
				html.Meta(html.PropertyAttr("og:image"), html.Content(image.URL)),
				html.If(image.Width > 0, html.Meta(html.PropertyAttr("og:image:width"), html.Content(fmt.Sprintf("%d", image.Width)))).(html.Node),
				html.If(image.Height > 0, html.Meta(html.PropertyAttr("og:image:height"), html.Content(fmt.Sprintf("%d", image.Height)))).(html.Node),
				html.If(len(image.Alt) > 0, html.Meta(html.PropertyAttr("og:image:alt"), html.Content(image.Alt))).(html.Node),
			)
		}),
		html.Map(o.Videos, func(video OpenGraphMedia) html.Node {
			return html.Fragment(
				html.Meta(html.PropertyAttr("og:video"), html.Content(video.URL)),
				html.If(video.Width > 0, html.Meta(html.PropertyAttr("og:video:width"), html.Content(fmt.Sprintf("%d", video.Width)))).(html.Node),
				html.If(video.Height > 0, html.Meta(html.PropertyAttr("og:video:height"), html.Content(fmt.Sprintf("%d", video.Height)))).(html.Node),
			)
		}),
		html.Map(o.Audio, func(audio OpenGraphMedia) html.Node {
			return html.Meta(html.PropertyAttr("og:audio"), html.Content(audio.URL))
		}),
	)
}

type Twitter struct {
	Card        string
	Title       string
	Description string
	SiteID      string
	Creator     string
	CreatorID   string
	Images      []string
}

func (t *Twitter) Node() html.Node {
	return html.Fragment(
		html.If(len(t.Card) > 0, html.Meta(html.Name("twitter:card"), html.Content(t.Card))).(html.Node),
		html.If(len(t.Title) > 0, html.Meta(html.Name("twitter:title"), html.Content(t.Title))).(html.Node),
		html.If(len(t.Description) > 0, html.Meta(html.Name("twitter:description"), html.Content(t.Description))).(html.Node),
		html.If(len(t.SiteID) > 0, html.Meta(html.Name("twitter:site:id"), html.Content(t.SiteID))).(html.Node),
		html.If(len(t.Creator) > 0, html.Meta(html.Name("twitter:creator"), html.Content(t.Creator))).(html.Node),
		html.If(len(t.CreatorID) > 0, html.Meta(html.Name("twitter:creator:id"), html.Content(t.CreatorID))).(html.Node),
		html.Map(t.Images, func(image string) html.Node {
			return html.Meta(html.Name("twitter:image"), html.Content(image))
		}),
	)
}

type Metadata struct {
	Base            *url.URL
	Title           string
	Description     string
	Generator       string
	ApplicationName string
	Referrer        string
	Keywords        []string
	Authors         []Author
	Creator         string
	Publisher       string
	Canonical       string
	Alternates      []Alternate
	Robots          *Robots
	OpenGraph       *OpenGraph
	Twitter         *Twitter
}

func (m *Metadata) Node() html.Node {
	return html.Fragment(
		html.If(len(m.Title) > 0, html.Title(html.Text(m.Title))).(html.Node),
		html.If(len(m.Description) > 0, html.Meta(html.Name("description"), html.Content(m.Description))).(html.Node),
		html.If(len(m.ApplicationName) > 0, html.Meta(html.Name("application-name"), html.Content(m.ApplicationName))).(html.Node),
		html.If(len(m.Generator) > 0, html.Meta(html.Name("generator"), html.Content(m.Generator))).(html.Node),
		html.If(len(m.Keywords) > 0, html.Meta(html.Name("keywords"), html.Content(strings.Join(m.Keywords, ",")))).(html.Node),
		html.If(len(m.Referrer) > 0, html.Meta(html.Name("referrer"), html.Content(m.Referrer))).(html.Node),
		html.If(len(m.Creator) > 0, html.Meta(html.Name("creator"), html.Content(m.Creator))).(html.Node),
		html.If(len(m.Publisher) > 0, html.Meta(html.Name("publisher"), html.Content(m.Publisher))).(html.Node),
		html.If(len(m.Canonical) > 0, html.Link(html.Rel("canonical"), html.Href(m.Canonical))).(html.Node),
		html.Map(m.Alternates, func(alternate Alternate) html.Node {
			return html.Link(
				html.Rel("alternate"),
				html.If(len(alternate.Href) > 0, html.Href(alternate.Href)),
				html.If(len(alternate.HrefLang) > 0, html.HrefLang(alternate.HrefLang)),
				html.If(len(alternate.Type) > 0, html.Type(alternate.Type)),
				html.If(len(alternate.Media) > 0, html.Media(alternate.Media)),
				html.If(len(alternate.Title) > 0, html.TitleAttr(alternate.Title)),
			)
		}),
		html.Map(m.Authors, func(author Author) html.Node {
			return html.Fragment(
				html.Meta(html.Name("author"), html.Content(author.Name)),
				html.If(len(author.URL) > 0, html.Link(html.Rel("author"), html.Href(author.URL))).(html.Node),
			)
		}),
		html.IfFn(m.Robots != nil, func() html.Item {
			return m.Robots.Node()
		}).(html.Node),
		html.IfFn(m.OpenGraph != nil, func() html.Item {
			return m.OpenGraph.Node()
		}).(html.Node),
		html.IfFn(m.Twitter != nil, func() html.Item {
			return m.Twitter.Node()
		}).(html.Node),
	)
}
