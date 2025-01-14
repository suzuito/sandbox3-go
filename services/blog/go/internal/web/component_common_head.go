package web

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

const SiteName = "otiuzu pages"

type siteMetaData struct {
	OGP       ogpData
	Canonical string
	LDJSON    []ldJSONData
}

type ogpData struct {
	Title       string
	Description string
	Locale      string
	Type        string
	URL         string
	SiteName    string
	Image       string
}

type ldJSONData struct {
	Context          string            `json:"@context"`
	Type             string            `json:"@type"`
	Headline         string            `json:"headline,omitempty"`
	Description      string            `json:"description,omitempty"`
	MainEntityOfPage string            `json:"mainEntityOfPage,omitempty"`
	DatePublished    string            `json:"datePublished,omitempty"`
	Author           *ldJSONDataAuthor `json:"author,omitempty"`
	ItemListElement  []ldJSONItem      `json:"itemListElement,omitempty"`
}

type ldJSONItem struct {
	Type     string `json:"@type"`
	Position int    `json:"position"`
	Name     string `json:"name"`
	Item     string `json:"item,omitempty"`
}

type ldJSONDataAuthor struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

func newSiteMetaDataFromContext(
	ctx *gin.Context,
	siteOrigin url.URL,
	title string,
	description string,
	t string,
	image string,
) *siteMetaData {
	return &siteMetaData{
		OGP: ogpData{
			Title:       title,
			Description: description,
			Locale:      "ja_JP",
			Type:        t,
			URL:         newPageURLFromContext(ctx, siteOrigin),
			Image:       image,
			SiteName:    SiteName,
		},
		Canonical: newPageURLFromContext(ctx, siteOrigin),
	}
}

func newPageURLFromContext(ctx *gin.Context, siteOrigin url.URL) string {
	if ctx.Request != nil {
		siteOrigin.Path = ctx.Request.URL.Path
	}
	return siteOrigin.String()
}

func newPageURL(siteOrigin url.URL, path string) string {
	siteOrigin.Path = path
	return siteOrigin.String()
}

type componentCommonHead struct {
	Title              string
	Meta               *siteMetaData
	GoogleTagManagerID string
}
