package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox3-go/services/blog/testutils/sqlcgo"
)

func Test_GET_admin_articles(t *testing.T) {
	ctx := context.Background()

	conn := newConn(ctx, &s)
	defer conn.Close(ctx)

	shutdown := e2ehelpers.RunServer(ctx, filePathServerBin, &e2ehelpers.RunServerInput{Envs: defaultEnvs(&s)}, healthCheck(ctx))
	defer shutdown() //nolint:errcheck

	cases := []e2ehelpers.PlaywrightTestCaseForSSR{
		{
			Desc: `ng - GET /admin/articles`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/admin/articles")
					require.NoError(t, err)
					require.Equal(t, http.StatusNotFound, res.Status())

					WriteHTML(t, res)

					require.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])
					requireHeader(t, page, false)
				}
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				MustTeardownDB(ctx, conn)
			},
		},
		{
			Desc: `ok - GET /admin/articles - empty articles, check charset=utf-8, header`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					require.NoError(t, page.Context().AddCookies([]playwright.OptionalCookie{
						{Name: "admin_auth_token", Value: "dummy_admin_token", URL: utils.Ptr("http://localhost:8081")},
					}))

					res, err := page.Goto("http://localhost:8080/admin/articles")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					WriteHTML(t, res)

					require.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])
					requireHeader(t, page, true)

					locsBreadcrumb := page.Locator(`[data-e2e-val="breadcrumb"]`)
					require.Equal(t, 2, Count(t, locsBreadcrumb))
					RequireElementInnerText(t, "トップページ ", locsBreadcrumb.Nth(0))
					RequireElementInnerText(t, "記事一覧", locsBreadcrumb.Nth(1))

					RequireElementNotExists(t, page.Locator(`[data-e2e-val="article"]`))
				}
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				MustTeardownDB(ctx, conn)
			},
		},
		{
			Desc: "ok - GET /admin/articles - check paging",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				MustSetupDB(
					ctx,
					conn,
					InitDBArg{
						Articles: sqlcgo.NewCreateArticlesParamsListAtRandom(
							0,
							time.Unix(1, 1),
							23,
						),
					},
				)
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					{
						require.NoError(t, page.Context().AddCookies([]playwright.OptionalCookie{
							{Name: "admin_auth_token", Value: "dummy_admin_token", URL: utils.Ptr("http://localhost:8081")},
						}))

						res, err := page.Goto("http://localhost:8080/admin/articles")
						require.NoError(t, err)
						require.Equal(t, http.StatusOK, res.Status())

						WriteHTML(t, res)

						loc := page.Locator(`[data-e2e-val="articles"]`)
						RequireElementExists(t, loc)

						loc = loc.Locator(`[data-e2e-val="article"]`)
						count := Count(t, loc)
						require.Equal(t, 10, count)

						{
							locArticle := loc.Nth(0)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事10 (Draft)", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementNotExists(t, locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}
						{
							locArticle := loc.Nth(count - 1)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事17", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-18", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}

						locPrevPage := page.Locator(`[data-e2e-val="prev-page"]`)
						RequireElementNotExists(t, locPrevPage)

						locNextPage := page.Locator(`[data-e2e-val="next-page"]`)
						RequireElementExists(t, locNextPage)

						require.NoError(t, locNextPage.Click())
					}
					{
						err := page.WaitForLoadState(
							playwright.PageWaitForLoadStateOptions{
								State: playwright.LoadStateLoad,
							},
						)
						require.NoError(t, err)

						loc := page.Locator(`[data-e2e-val="articles"]`)
						RequireElementExists(t, loc)

						loc = loc.Locator(`[data-e2e-val="article"]`)
						count := Count(t, loc)
						require.Equal(t, 10, count)

						{
							locArticle := loc.Nth(0)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事16", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-17", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}
						{
							locArticle := loc.Nth(count - 1)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事4", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-05", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}

						locPrevPage := page.Locator(`[data-e2e-val="prev-page"]`)
						RequireElementExists(t, locPrevPage)

						locNextPage := page.Locator(`[data-e2e-val="next-page"]`)
						RequireElementExists(t, locNextPage)

						require.NoError(t, locNextPage.Click())
					}
					{
						err := page.WaitForLoadState(
							playwright.PageWaitForLoadStateOptions{
								State: playwright.LoadStateLoad,
							},
						)
						require.NoError(t, err)

						loc := page.Locator(`[data-e2e-val="articles"]`)
						RequireElementExists(t, loc)

						loc = loc.Locator(`[data-e2e-val="article"]`)
						count := Count(t, loc)
						require.Equal(t, 3, count)

						{
							locArticle := loc.Nth(0)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事3", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-04", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}
						{
							locArticle := loc.Nth(count - 1)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事1", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-02", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}

						locPrevPage := page.Locator(`[data-e2e-val="prev-page"]`)
						RequireElementExists(t, locPrevPage)

						locNextPage := page.Locator(`[data-e2e-val="next-page"]`)
						RequireElementNotExists(t, locNextPage)
					}
				}
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				MustTeardownDB(ctx, conn)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			c.Run(t)
		})
	}
}
