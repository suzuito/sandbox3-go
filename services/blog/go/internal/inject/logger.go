package inject

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"github.com/suzuito/sandbox2-common-go/libs/clog"
)

func NewLogger(env *Environment) *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(env.LogLevel)); err != nil {
		fmt.Printf("use LogLevel 'DEBUG' because cannot parse LOG_LEVEL: %s", env.LogLevel)
		level = slog.LevelDebug
	}

	var slogHandler slog.Handler
	switch env.LoggerType {
	case "devslog":
		slogHandler = devslog.NewHandler(os.Stdout, &devslog.Options{
			HandlerOptions: &slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		})
	case "json":
		slogHandler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.LevelKey {
						a.Key = "severity"
					}
					return a
				},
			},
		)
	default:
		slogHandler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		)
	}

	slogCustomHandler := clog.CustomHandler{
		Handler: slogHandler,
	}
	return slog.New(&slogCustomHandler)
}
