package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
)

func TestMain(m *testing.M) {
	err := playwright.Install()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func healthCheck(ctx context.Context) func() error {
	return func() error {
		return CheckHTTPServerHealth(ctx, "http://localhost:8080/health")
	}
}

func TestXxx(t *testing.T) {
	filePathServerBin := os.Getenv("FILE_PATH_SERVER_BIN")
	envs := []string{
		"ENV=loc",
		"SITE_ORIGIN=http://localhost:9003",
		"GOOGLE_TAG_MANAGER_ID=dummy_tag_id",
		"ADMIN_TOKEN=dummy_admin_token",
		"DIR_PATH_HTML_TEMPLATE=../go/internal/web",
		"DIR_PATH_CSS=../go/internal/web/_css",
		"LOGGER_TYPE=json",
	}

	ctx := context.Background()
	shutdown := RunServer(ctx, filePathServerBin, &RunServerInput{Envs: envs}, healthCheck(ctx))
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
		{
			Desc: "ok - GET /",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])

					locHeader := page.Locator(`[data-e2e-val="header"]`)
					e2ehelpers.AssertElementExists(t, locHeader)

					locLinkToAdmin := locHeader.Locator(`[data-e2e-val="link-to-admin"]`)
					e2ehelpers.AssertElementNotExists(t, locLinkToAdmin)

					locFooter := page.Locator(`[data-e2e-val="footer"]`)
					e2ehelpers.AssertElementExists(t, locFooter)
				}
			},
		},
		{
			Desc: "ok - GET / as admin",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					require.NoError(t, page.Context().AddCookies([]playwright.OptionalCookie{
						{Name: "admin_auth_token", Value: "dummy_admin_token", URL: utils.Ptr("http://localhost:8080")},
					}))

					res, err := page.Goto("http://localhost:8080")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])

					locHeader := page.Locator(`[data-e2e-val="header"]`)
					e2ehelpers.AssertElementExists(t, locHeader)

					locLinkToAdmin := locHeader.Locator(`[data-e2e-val="link-to-admin"]`)
					e2ehelpers.AssertElementExists(t, locLinkToAdmin)

					locFooter := page.Locator(`[data-e2e-val="footer"]`)
					e2ehelpers.AssertElementExists(t, locFooter)
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

type RunServerInput struct {
	Args []string
	Envs []string
}

func RunServer(
	ctx context.Context,
	filePathBin string,
	input *RunServerInput,
	healthCheckFunc func() error,
) func() error {
	cmd := exec.CommandContext(
		ctx,
		filePathBin,
		input.Args...,
	)

	cmd.Env = append(
		os.Environ(),
		input.Envs...,
	)

	stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		var exiterr *exec.ExitError
		if !errors.As(err, &exiterr) {
			panic(fmt.Sprintf("%s %s: %s", filePathBin, strings.Join(input.Args, " "), err.Error()))
		}
	}

	if err := healthCheckFunc(); err != nil {
		panic("health check error")
	}

	return func() error {
		defer func() {
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println("@@@@@@@ STDOUT @@@@@@@")
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println(stdout.String())
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			fmt.Println()

			if stderr.Len() > 0 {
				fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
				fmt.Println("@@@@@@@ STDERR @@@@@@@")
				fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
				fmt.Println(stderr.String())
				fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
				fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
				fmt.Println("@@@@@@@@@@@@@@@@@@@@@@")
			}
		}()

		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}

		return nil
	}
}

func CheckHTTPServerHealth(
	ctx context.Context,
	u string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cli := http.DefaultClient

	for {
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			return errors.New("health check is failed")
		default:
			res, err := cli.Get(u)
			if err != nil {
				continue
			}

			res.Body.Close()
			if res.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
}
