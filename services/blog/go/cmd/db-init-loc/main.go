package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/inject"
	"github.com/suzuito/sandbox3-go/services/blog/testutils/sqlcgo"
)

type Environment struct {
	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     uint16 `envconfig:"DB_PORT" required:"true"`
	DBName     string `envconfig:"DB_NAME" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
}

func (t *Environment) DBURI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		t.DBUser, t.DBPassword,
		t.DBHost, t.DBPort,
		t.DBName,
	)
}

func main() {
	var env Environment
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load environment variable: %v\n", err)
		panic(err)
	}

	ctx := context.Background()
	pgxConn, err := inject.NewPgxConn(ctx, &env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create pgx connection: %v\n", err)
		panic(err)
	}

	tx, err := pgxConn.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to begin tx: %v\n", err)
		panic(err)
	}
	defer tx.Rollback(ctx) // nolint:errcheck

	queries := sqlcgo.New(tx)

	articles := sqlcgo.NewCreateArticlesParamsListAtRandom(
		0,
		time.Now(),
		100,
	)
	createdArticles, err := queries.CreateArticles(ctx, articles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create articles: %v\n", err)
		panic(err)
	}
	fmt.Printf("created articles: %d\n", createdArticles)

	tags := sqlcgo.NewCreateTagsParamsListAtRandom(
		0,
		time.Now(),
		100,
	)
	createdTags, err := queries.CreateTags(ctx, tags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create articles: %v\n", err)
		panic(err)
	}
	fmt.Printf("created tags: %d\n", createdTags)

	relArticlesTags := []sqlcgo.CreateRelArticlesTagsParams{}
	tagI := 0
	for _, article := range articles {
		for range 10 {
			tagI = (tagI + 1) % len(tags)
			relArticlesTags = append(
				relArticlesTags,
				sqlcgo.CreateRelArticlesTagsParams{
					ArticleID: article.ID,
					TagID:     tags[tagI].ID,
				},
			)
		}
	}
	createdRelArticlesTags, err := queries.CreateRelArticlesTags(
		ctx, relArticlesTags,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create rel article tags: %v\n", err)
		panic(err)
	}
	fmt.Printf("created rel articles tags: %d\n", createdRelArticlesTags)

	if err := queries.UpsertAllArticleSearchIndices(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create article search indices: %v\n", err)
		panic(err)
	}
	fmt.Println("created search indices")

	if err := tx.Commit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to commit: %v\n", err)
		panic(err)
	}
}
