package main

import (
	"net/http"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestXxx(t *testing.T) {
	cases := []PlaywrightTestCaseForSSR{
		{
			Desc: "ok - GET /health",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:9003/health")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "noindex", res.Headers()["x-robots-tag"])
				}
			},
		},
		{
			Desc: "ok - GET /",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:9003")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])
				}
			},
		},
		/*
			{
				Desc: "ok - sitemap.xml",
				Setup: func(t *testing.T, testID e2ehelpers.TestID, expected *HTTPServerTestCaseExpected) *http.Request {
					expected.StatusCode = http.StatusOK
					expected.Header.Set("X-Robots-Tag", "noindex")
					expected.Header.Set("Content-Type", "application/xml")

					return MustHTTPRequest(http.MethodGet, "http://localhost:9003/sitemap.xml", nil)
				},
			},
			{
				Desc: "ok - rss",
				Setup: func(t *testing.T, testID e2ehelpers.TestID, expected *HTTPServerTestCaseExpected) *http.Request {
					expected.StatusCode = http.StatusOK
					expected.Header.Set("X-Robots-Tag", "noindex")
					expected.Header.Set("Content-Type", "application/rss+xml")

					return MustHTTPRequest(http.MethodGet, "http://localhost:9003/sitemap.xml", nil)
				},
			},
		*/
	}
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			c.Run(t)
		})
	}
}
