package logger

import (
	"log/slog"
	"os"
	"strings"
)

// DefaultLogger is a default logger.
var DefaultProdLogger = slog.New(
	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "level" {
				return slog.Attr{}
			}
			if a.Key == slog.SourceKey {
				str := a.Value.String()
				split := strings.Split(str, "/")
				if len(split) > 2 {
					a.Value = slog.StringValue(
						strings.Join(split[len(split)-2:], "/"),
					)
					a.Value = slog.StringValue(
						strings.ReplaceAll(a.Value.String(), "}", ""),
					)
				}
			}

			return a
		}}),
)
