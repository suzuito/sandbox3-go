package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/infra/rdb/repositories"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/inject"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/usecases"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/web"
)

func main() {
	var env inject.Environment
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load environment variable: %v\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := inject.NewLogger(&env)
	pgxConn, err := inject.NewPgxConn(ctx, &env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create pgx connection: %v\n", err)
		os.Exit(1)
	}

	repo := repositories.NewImpl(pgxConn)
	uc := usecases.NewImpl(repo)

	w, err := web.New(&env, logger, uc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to new web: %v\n", err)
		os.Exit(1)
	}

	e := gin.New()
	e.Use(gin.Recovery())
	w.SetEngine(e)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Port),
		Handler: e.Handler(),
	}

	os.Exit(utils.RunHTTPServerWithGracefulShutdown(
		ctx,
		server,
		logger,
	))
}
