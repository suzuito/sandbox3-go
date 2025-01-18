package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/inject"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/web"
)

func main() {
	var env inject.Environment
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load environment variable: %v\n", err)
		os.Exit(1)
	}

	logger := inject.NewLogger(&env)

	w, err := web.New(&env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to new web: %v\n", err)
		os.Exit(1)
	}

	e := gin.New()
	e.Use(gin.Recovery())
	w.SetEngine(e)

	server := &http.Server{
		Addr:    ":8080",
		Handler: e.Handler(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	os.Exit(utils.RunHTTPServerWithGracefulShutdown(
		ctx,
		server,
		logger,
	))
}
