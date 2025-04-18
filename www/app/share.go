package app

import (
	"net/url"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/html"
)

//pacis:page path=/share/{slug}
func SharePage(ctx *pages.PageContext) I {
	slug := ctx.Request().PathValue("slug")

	query := url.Values{}
	query.Add("medium", "cpc")

	if cachedb == nil {
		query.Add("campaign", "promo")
	} else {
		campaign, err := cachedb.Get(ctx, "campaign").Result()
		if err == nil {
			query.Add("campaign", campaign)
		} else {
			query.Add("campaign", "promo")
		}
	}

	switch slug {
	case "reddit":
		query.Add("source", "reddit")
	case "x":
		query.Add("source", "x")
	case "bsky":
		query.Add("source", "bsky")
	default:
		query.Add("source", "unknown")
	}

	return ctx.Redirect("/?" + query.Encode())
}
