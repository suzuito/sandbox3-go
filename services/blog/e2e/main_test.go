package main

import (
	"net/http"
	"testing"

	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestXxx(t *testing.T) {
	cases := []HTTPServerTestCase{
		{
			Desc: "ok - GET /health",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, expected *HTTPServerTestCaseExpected) *http.Request {
				expected.StatusCode = http.StatusOK
				expected.Header.Set("X-Robots-Tag", "noindex")

				return MustHTTPRequest(http.MethodGet, "http://localhost:9003/health", nil)
			},
		},
		{
			Desc: "ok - GET /",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, expected *HTTPServerTestCaseExpected) *http.Request {
				expected.StatusCode = http.StatusOK
				expected.Header.Set("Content-type", "text/html; charset=utf-8")

				return MustHTTPRequest(http.MethodGet, "http://localhost:9003", nil)
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
