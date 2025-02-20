package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/playwright-community/playwright-go"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

type setting struct {
	Port       int
	DBHost     string
	DBPort     uint16
	DBName     string
	DBUser     string
	DBPassword string
}

func (t *setting) DBURI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		t.DBUser, t.DBPassword,
		t.DBHost, t.DBPort,
		t.DBName,
	)
}

func newConn(ctx context.Context, s *setting) *pgx.Conn {
	cfg, err := pgx.ParseConfig(s.DBURI())
	if err != nil {
		panic(err)
	}

	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}

	return conn
}

func defaultEnvs(s *setting) []string {
	return []string{
		"ENV=loc",
		fmt.Sprintf("PORT=%d", s.Port),
		"SITE_ORIGIN=http://localhost:9003",
		"GOOGLE_TAG_MANAGER_ID=dummy_tag_id",
		"ADMIN_TOKEN=hoge",
		"DIR_PATH_HTML_TEMPLATE=../go/internal/web",
		"DIR_PATH_CSS=../go/internal/web/_css",
		"LOGGER_TYPE=json",
		fmt.Sprintf("DB_HOST=%s", s.DBHost),
		fmt.Sprintf("DB_PORT=%d", s.DBPort),
		fmt.Sprintf("DB_NAME=%s", s.DBName),
		fmt.Sprintf("DB_USER=%s", s.DBUser),
		fmt.Sprintf("DB_PASSWORD=%s", s.DBPassword),
	}
}

var filePathServerBin string

func TestMain(m *testing.M) {
	filePathServerBin = os.Getenv("FILE_PATH_SERVER_BIN")

	err := playwright.Install()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func healthCheck(ctx context.Context) func() error {
	return func() error {
		return e2ehelpers.CheckHTTPServerHealth(ctx, "http://localhost:8080/health")
	}
}

var s = setting{
	Port:       8080,
	DBHost:     "blog-db",
	DBPort:     5432,
	DBName:     "blog_test",
	DBUser:     "root",
	DBPassword: "root",
}

func TestBlogService(t *testing.T) {
	ctx := context.Background()

	conn := newConn(ctx, &s)
	defer conn.Close(ctx)

	shutdown := e2ehelpers.RunServer(ctx, filePathServerBin, &e2ehelpers.RunServerInput{Envs: defaultEnvs(&s)}, healthCheck(ctx))
	defer shutdown() //nolint:errcheck

	cases := []e2ehelpers.PlaywrightTestCaseForSSR{
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
