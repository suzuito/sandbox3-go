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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox3-go/services/blog/testutils/sqlcgo"
)

type setting struct {
	Port       int
	DBHost     string
	DBPort     uint16
	DBName     string
	DBUser     string
	DBPassword string
}

var s = setting{
	Port:       8080,
	DBHost:     "blog-db",
	DBPort:     5432,
	DBName:     "blog_test",
	DBUser:     "root",
	DBPassword: "root",
}

func newConn(ctx context.Context) *pgx.Conn {
	runtimeParams := map[string]string{
		"sslmode": "disable",
	}

	conf := pgx.ConnConfig{
		Config: pgconn.Config{
			Host:          s.DBHost,
			Port:          s.DBPort,
			User:          s.DBUser,
			Password:      s.DBPassword,
			Database:      s.DBName,
			RuntimeParams: runtimeParams,
		},
	}

	conn, err := pgx.ConnectConfig(ctx, &conf)
	if err != nil {
		panic(err)
	}

	return conn
}

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

func TestBlogService(t *testing.T) {
	ctx := context.Background()

	filePathServerBin := os.Getenv("FILE_PATH_SERVER_BIN")
	envs := []string{
		"ENV=loc",
		fmt.Sprintf("PORT=%d", s.Port),
		"SITE_ORIGIN=http://localhost:9003",
		"GOOGLE_TAG_MANAGER_ID=dummy_tag_id",
		"ADMIN_TOKEN=dummy_admin_token",
		"DIR_PATH_HTML_TEMPLATE=../go/internal/web",
		"DIR_PATH_CSS=../go/internal/web/_css",
		"LOGGER_TYPE=json",
		fmt.Sprintf("DB_HOST=%s", s.DBHost),
		fmt.Sprintf("DB_PORT=%d", s.DBPort),
		fmt.Sprintf("DB_NAME=%s", s.DBName),
		fmt.Sprintf("DB_USER=%s", s.DBUser),
		fmt.Sprintf("DB_PASSWORD=%s", s.DBPassword),
	}

	conn := newConn(ctx)
	defer conn.Close(ctx)

	shutdown := RunServer(ctx, filePathServerBin, &RunServerInput{Envs: envs}, healthCheck(ctx))
	defer shutdown() //nolint:errcheck

	queries := sqlcgo.New(conn)

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

					assertHeader(t, page, false)
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

					assertHeader(t, page, true)
				}
			},
		},
		{
			Desc: "ok - GET /articles",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, exe *e2ehelpers.PlaywrightTestCaseForSSRExec) {
				// TODO 続きはここから 2025/01/25
				queries.CreateArticles(ctx, []sqlcgo.CreateArticlesParams{
					{},
				})
				exe.Do = func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
					res, err := page.Goto("http://localhost:8080/articles")
					require.NoError(t, err)

					assert.Equal(t, http.StatusOK, res.Status())
					assert.Equal(t, "text/html; charset=utf-8", res.Headers()["content-type"])

					assertHeader(t, page, false)
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
	printStdoutStderr := func() {
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
	}
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
		printStdoutStderr()
		panic(fmt.Errorf("health check error: %w", err))
	}

	return func() error {
		defer func() {
			printStdoutStderr()
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
