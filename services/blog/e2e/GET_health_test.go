package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func Test_GET_health(t *testing.T) {
	ctx := context.Background()

	conn := newConn(ctx, &s)
	defer conn.Close(ctx)

	shutdown := e2ehelpers.RunServer(ctx, filePathServerBin, &e2ehelpers.RunServerInput{Envs: defaultEnvs(&s)}, healthCheck(ctx))
	defer shutdown() //nolint:errcheck

	cases := []e2ehelpers.PlaywrightTestCaseForSSR{
		{
			Desc: "ok - GET /health",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/health")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "noindex", res.Headers()["x-robots-tag"])
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
