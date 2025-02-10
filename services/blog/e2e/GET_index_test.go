package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
)

func Test_GET_index(t *testing.T) {
	ctx := context.Background()

	conn := newConn(ctx, &s)
	defer conn.Close(ctx)

	shutdown := e2ehelpers.RunServer(ctx, filePathServerBin, &e2ehelpers.RunServerInput{Envs: defaultEnvs(&s)}, healthCheck(ctx))
	defer shutdown() //nolint:errcheck

	cases := []e2ehelpers.PlaywrightTestCaseForSSR{
		{
			Desc: "ok - GET /",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])

					requireHeader(t, page, false)
				}
			},
		},
		{
			Desc: "ok - GET / as admin",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					require.NoError(t, page.Context().AddCookies([]playwright.OptionalCookie{
						{Name: "admin_auth_token", Value: "dummy_admin_token", URL: utils.Ptr("http://localhost:8081")},
					}))

					res, err := page.Goto("http://localhost:8080")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])

					requireHeader(t, page, true)
				}
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			c.Run(t)
		})
	}
}
