package inject

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

type CustomHandler struct {
	slog.Handler
}

func (t *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			terr, ok := a.Value.Any().(terrors.TraceableError)
			if ok {
				traceInfos := []string{}
				for _, st := range terr.StackTrace() {
					traceInfos = append(traceInfos, fmt.Sprintf("%s:%d", st.Filename, st.Line))
				}
				r.AddAttrs(slog.Attr{Key: "traceInfos", Value: slog.AnyValue(traceInfos)})
			}
			return false
		}
		return true
	})
	return t.Handler.Handle(ctx, r)
}

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

	slogCustomHandler := CustomHandler{
		Handler: slogHandler,
	}
	return slog.New(&slogCustomHandler)
}
