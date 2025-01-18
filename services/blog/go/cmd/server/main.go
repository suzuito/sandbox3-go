package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/inject"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/web"
)

func RunServer(
	server *http.Server,
	logger *slog.Logger,
) int {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	chGracefulShutdown := make(chan error)
	defer close(chGracefulShutdown)
	go func() {
		// Signalのハンドラー
		// SIGINT,SIGTERMをキャッチした後、ctx.Doneが制御を返す
		<-ctx.Done()
		logger.Info("start graceful shut down")
		ctxSignalHandler, cancel := context.WithTimeout(context.Background(), time.Second*100) // 100秒待ってもserver.Shutdown(ctx)が返ってこない場合、強制的にシャットダウンする
		defer cancel()
		// Graceful Shutdownをスタートする。
		// Graceful Shutdownが成功したら、server.Shutdown(ctxSignalHandler)はnilを返す。
		// Graceful Shutdownが失敗したら、server.Shutdown(ctxSignalHandler)は非nilを返す。
		chGracefulShutdown <- server.Shutdown(ctxSignalHandler)
	}()

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server.ListenAndServe() is failed", "err", err)
			return 1
		}
	}
	if err := <-chGracefulShutdown; err != nil {
		logger.Error("graceful shut down is failed (server.Shutdown(ctx) is failed)", "err", err)
		return 2
	}
	logger.Info("graceful shut down is complete")
	return 0
}

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
