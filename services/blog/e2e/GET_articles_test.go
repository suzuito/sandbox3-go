package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox3-go/services/blog/testutils/sqlcgo"
)

func Test_GET_articles(t *testing.T) {
	ctx := context.Background()

	conn := newConn(ctx, &s)
	defer conn.Close(ctx)

	shutdown := e2ehelpers.RunServer(ctx, filePathServerBin, &e2ehelpers.RunServerInput{Envs: defaultEnvs(&s)}, healthCheck(ctx))
	defer shutdown() //nolint:errcheck

	cases := []e2ehelpers.PlaywrightTestCaseForSSR{
		{
			Desc: `ok - GET /articles - empty articles, check charset=utf-8, header`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/articles")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					WriteHTML(t, res)

					require.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])
					requireHeader(t, page, false)

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
			Desc: `ok - GET /articles - check Article representation in HTML
article without tags, article with tags, deleted article is ignored, deleted tag is ignored
			`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				MustSetupDB(
					ctx,
					conn,
					InitDBArg{
						Articles: sqlcgo.CreateArticlesParamsList{
							{
								ID:          uuid.MustParse("f3f17eb5-c083-4df4-a022-cc3ddbed6826"),
								Title:       "published with no tags",
								PublishedAt: sqlcgo.NewPgTypeFromTime(time.Unix(2, 1)),
								CreatedAt:   sqlcgo.NewPgTypeFromTime(time.Unix(2, 1)),
								UpdatedAt:   sqlcgo.NewPgTypeFromTime(time.Unix(2, 1)),
								DeletedAt:   sqlcgo.NewNilPgType(),
							},
							{
								ID:          uuid.MustParse("a8fc997b-65c3-4d62-b289-0bec281a5b1f"),
								Title:       "published with tags",
								PublishedAt: sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								CreatedAt:   sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								UpdatedAt:   sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								DeletedAt:   sqlcgo.NewNilPgType(),
							},
							{
								ID:          uuid.MustParse("272d11db-9e15-4ae1-a0d7-1ef91b80d855"),
								Title:       "deleted",
								PublishedAt: sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								CreatedAt:   sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								UpdatedAt:   sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								DeletedAt:   sqlcgo.NewPgTypeFromTime(time.Unix(3, 1)),
							},
						},
						Tags: sqlcgo.CreateTagsParamsList{
							{
								ID:        uuid.MustParse("05d29618-e41e-4d05-adf6-09a810d2be97"),
								Name:      "タグ1",
								CreatedAt: sqlcgo.NewPgTypeFromTime(time.Unix(2, 1)),
								UpdatedAt: sqlcgo.NewPgTypeFromTime(time.Unix(2, 1)),
								DeletedAt: sqlcgo.NewNilPgType(),
							},
							{
								ID:        uuid.MustParse("70401edb-975b-4f2e-b091-c34b72a1a38c"),
								Name:      "タグ2",
								CreatedAt: sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								UpdatedAt: sqlcgo.NewPgTypeFromTime(time.Unix(1, 1)),
								DeletedAt: sqlcgo.NewNilPgType(),
							},
							{
								ID:        uuid.MustParse("7b668c5c-9259-4a83-ae42-64f6d04eab69"),
								Name:      "タグ3",
								CreatedAt: sqlcgo.NewPgTypeFromTime(time.Unix(3, 1)),
								UpdatedAt: sqlcgo.NewPgTypeFromTime(time.Unix(3, 1)),
								DeletedAt: sqlcgo.NewPgTypeFromTime(time.Unix(4, 1)),
							},
						},
						RelArticlesTags: sqlcgo.CreateRelArticlesTagsParamsList{
							{
								ArticleID: uuid.MustParse("a8fc997b-65c3-4d62-b289-0bec281a5b1f"),
								TagID:     uuid.MustParse("05d29618-e41e-4d05-adf6-09a810d2be97"),
							},
							{
								ArticleID: uuid.MustParse("a8fc997b-65c3-4d62-b289-0bec281a5b1f"),
								TagID:     uuid.MustParse("70401edb-975b-4f2e-b091-c34b72a1a38c"),
							},
							{
								ArticleID: uuid.MustParse("a8fc997b-65c3-4d62-b289-0bec281a5b1f"),
								TagID:     uuid.MustParse("7b668c5c-9259-4a83-ae42-64f6d04eab69"),
							},
						},
					},
				)
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/articles")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					WriteHTML(t, res)

					locArticles := page.Locator(`[data-e2e-val="articles"]`)
					e2ehelpers.AssertElementExists(t, locArticles)
					locsArticle := RequireElementsCount(t, 2, locArticles.Locator(`[data-e2e-val="article"]`))
					{
						locArticle := locsArticle[0]
						RequireElementInnerText(t, "published with no tags", locArticle.Locator(`[data-e2e-val="article-title"]`))
						RequireElementInnerText(t, "1970-01-01", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						RequireElementHasAttribute(t, "/articles/f3f17eb5-c083-4df4-a022-cc3ddbed6826", locArticle.Locator(`[data-e2e-val="article-link"]`), "href")
						e2ehelpers.AssertElementNotExists(t, locArticle.Locator(`[data-e2e-val="tags"]`))
					}
					{
						locArticle := locsArticle[1]
						RequireElementInnerText(t, "published with tags", locArticle.Locator(`[data-e2e-val="article-title"]`))
						RequireElementInnerText(t, "1970-01-01", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						RequireElementHasAttribute(t, "/articles/a8fc997b-65c3-4d62-b289-0bec281a5b1f", locArticle.Locator(`[data-e2e-val="article-link"]`), "href")

						locTags := locArticle.Locator(`[data-e2e-val="tags"]`)
						e2ehelpers.AssertElementExists(t, locTags)

						locsTag := RequireElementsCount(t, 2, locTags.Locator(`[data-e2e-val="tag"]`))
						{
							locTag := locsTag[0]
							RequireElementInnerText(t, "タグ1", locTag)
							RequireElementHasAttribute(t, `/articles?tag=%e3%82%bf%e3%82%b01`, locTag.Locator(`[data-e2e-val="tag-link"]`), "href")
						}
						{
							locTag := locsTag[1]
							RequireElementInnerText(t, "タグ2", locTag)
							RequireElementHasAttribute(t, `/articles?tag=%e3%82%bf%e3%82%b02`, locTag.Locator(`[data-e2e-val="tag-link"]`), "href")
						}
					}
				}
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				MustTeardownDB(ctx, conn)
			},
		},
		{
			Desc: "ok - GET /articles - check paging",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				MustSetupDB(
					ctx,
					conn,
					InitDBArg{
						Articles: sqlcgo.NewCreateArticlesParamsListAtRandom(
							0,
							time.Unix(1, 1),
							33,
						),
					},
				)
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					{
						res, err := page.Goto("http://localhost:8080/articles")
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
							RequireElementInnerText(t, "テスト記事32", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-02-02", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}
						{
							locArticle := loc.Nth(count - 1)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事23", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-24", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
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
							RequireElementInnerText(t, "テスト記事22", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-23", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}
						{
							locArticle := loc.Nth(count - 1)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事13", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-14", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
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
						require.Equal(t, 10, count)

						{
							locArticle := loc.Nth(0)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事12", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-13", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}
						{
							locArticle := loc.Nth(count - 1)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事3", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-04", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
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
							RequireElementInnerText(t, "テスト記事2", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-03", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
						}
						{
							locArticle := loc.Nth(count - 1)
							RequireElementExists(t, locArticle)
							RequireElementInnerText(t, "テスト記事0", locArticle.Locator(`[data-e2e-val="article-title"]`))
							RequireElementInnerText(t, "1970-01-01", locArticle.Locator(`[data-e2e-val="article-published-at"]`))
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
		{
			Desc: `ok - GET /articles - find by tag`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				MustSetupDB(
					ctx,
					conn,
					InitDBArg{
						Articles: sqlcgo.CreateArticlesParamsList{
							{
								ID:          uuid.MustParse("fe91cdde-8dba-44be-8712-b5f6785caa46"),
								Title:       "記事1",
								PublishedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								CreatedAt:   sqlcgo.NewPgTypeFromTime(time.Now()),
								UpdatedAt:   sqlcgo.NewPgTypeFromTime(time.Now()),
								DeletedAt:   sqlcgo.NewNilPgType(),
							},
							{
								ID:          uuid.MustParse("4655a3bb-b3d9-4dc0-9760-d3b736905e0b"),
								Title:       "記事2",
								PublishedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								CreatedAt:   sqlcgo.NewPgTypeFromTime(time.Now()),
								UpdatedAt:   sqlcgo.NewPgTypeFromTime(time.Now()),
								DeletedAt:   sqlcgo.NewNilPgType(),
							},
						},
						Tags: sqlcgo.CreateTagsParamsList{
							{
								ID:        uuid.MustParse("545e447d-f4e6-4fbb-b428-e8133034720f"),
								Name:      "タグ1",
								CreatedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								UpdatedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								DeletedAt: sqlcgo.NewNilPgType(),
							},
							{
								ID:        uuid.MustParse("570acc2a-e967-46cb-a719-1b85f85b32b6"),
								Name:      "タグ2",
								CreatedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								UpdatedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								DeletedAt: sqlcgo.NewNilPgType(),
							},
						},
						RelArticlesTags: sqlcgo.CreateRelArticlesTagsParamsList{
							{
								ArticleID: uuid.MustParse("fe91cdde-8dba-44be-8712-b5f6785caa46"),
								TagID:     uuid.MustParse("545e447d-f4e6-4fbb-b428-e8133034720f"),
							},
							{
								ArticleID: uuid.MustParse("4655a3bb-b3d9-4dc0-9760-d3b736905e0b"),
								TagID:     uuid.MustParse("570acc2a-e967-46cb-a719-1b85f85b32b6"),
							},
						},
					},
				)
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/articles?tag=タグ1")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					WriteHTML(t, res)

					locsArticle := page.Locator(`[data-e2e-val="article"]`)
					count := Count(t, locsArticle)
					require.Equal(t, 1, count)
					{
						locArticle := locsArticle.Nth(0)
						RequireElementInnerText(t, "記事1", locArticle.Locator(`[data-e2e-val="article-title"]`))
					}
				}
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				MustTeardownDB(ctx, conn)
			},
		},
		{
			Desc: `ok - GET /articles - find by tag that doesn't exist`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				MustSetupDB(
					ctx,
					conn,
					InitDBArg{
						Articles: sqlcgo.CreateArticlesParamsList{
							{
								ID:          uuid.MustParse("fe91cdde-8dba-44be-8712-b5f6785caa46"),
								Title:       "記事1",
								PublishedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								CreatedAt:   sqlcgo.NewPgTypeFromTime(time.Now()),
								UpdatedAt:   sqlcgo.NewPgTypeFromTime(time.Now()),
								DeletedAt:   sqlcgo.NewNilPgType(),
							},
						},
						Tags: sqlcgo.CreateTagsParamsList{
							{
								ID:        uuid.MustParse("545e447d-f4e6-4fbb-b428-e8133034720f"),
								Name:      "タグ1",
								CreatedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								UpdatedAt: sqlcgo.NewPgTypeFromTime(time.Now()),
								DeletedAt: sqlcgo.NewNilPgType(),
							},
						},
						RelArticlesTags: sqlcgo.CreateRelArticlesTagsParamsList{
							{
								ArticleID: uuid.MustParse("fe91cdde-8dba-44be-8712-b5f6785caa46"),
								TagID:     uuid.MustParse("545e447d-f4e6-4fbb-b428-e8133034720f"),
							},
						},
					},
				)
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/articles?tag=タグxx")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					WriteHTML(t, res)

					locsArticle := page.Locator(`[data-e2e-val="article"]`)
					RequireElementNotExists(t, locsArticle)
				}
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				MustTeardownDB(ctx, conn)
			},
		},
		{
			Desc: `ok - GET /articles - find by time range`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				tm := time.Unix(100, 0)
				MustSetupDB(
					ctx,
					conn,
					InitDBArg{
						Articles: sqlcgo.CreateArticlesParamsList{
							{
								ID:          uuid.MustParse("fe91cdde-8dba-44be-8712-b5f6785caa46"),
								Title:       "記事1",
								PublishedAt: sqlcgo.NewPgTypeFromTime(tm),
								CreatedAt:   sqlcgo.NewPgTypeFromTime(tm),
								UpdatedAt:   sqlcgo.NewPgTypeFromTime(tm),
								DeletedAt:   sqlcgo.NewNilPgType(),
							},
						},
					},
				)
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/articles?since=99")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					locsArticle := page.Locator(`[data-e2e-val="article"]`)
					count := Count(t, locsArticle)
					require.Equal(t, 1, count)
					{
						locArticle := locsArticle.Nth(0)
						RequireElementInnerText(t, "記事1", locArticle.Locator(`[data-e2e-val="article-title"]`))
					}

					res, err = page.Goto("http://localhost:8080/articles?since=100")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					locsArticle = page.Locator(`[data-e2e-val="article"]`)
					count = Count(t, locsArticle)
					require.Equal(t, 1, count)
					{
						locArticle := locsArticle.Nth(0)
						RequireElementInnerText(t, "記事1", locArticle.Locator(`[data-e2e-val="article-title"]`))
					}

					res, err = page.Goto("http://localhost:8080/articles?since=101")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					locsArticle = page.Locator(`[data-e2e-val="article"]`)
					RequireElementNotExists(t, locsArticle)

					res, err = page.Goto("http://localhost:8080/articles?until=101")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					locsArticle = page.Locator(`[data-e2e-val="article"]`)
					count = Count(t, locsArticle)
					require.Equal(t, 1, count)
					{
						locArticle := locsArticle.Nth(0)
						RequireElementInnerText(t, "記事1", locArticle.Locator(`[data-e2e-val="article-title"]`))
					}

					res, err = page.Goto("http://localhost:8080/articles?until=100")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					locsArticle = page.Locator(`[data-e2e-val="article"]`)
					count = Count(t, locsArticle)
					require.Equal(t, 1, count)
					{
						locArticle := locsArticle.Nth(0)
						RequireElementInnerText(t, "記事1", locArticle.Locator(`[data-e2e-val="article-title"]`))
					}

					res, err = page.Goto("http://localhost:8080/articles?until=99")
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.Status())

					locsArticle = page.Locator(`[data-e2e-val="article"]`)
					RequireElementNotExists(t, locsArticle)
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
