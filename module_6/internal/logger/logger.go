package logger

import (
	"context"
	"log/slog"
	"module6/internal/config"
	"os"
)

var logger *Logger

func GetLogger() *Logger {
	return logger
}

type Logger struct {
	cfg *config.Config
	log *slog.Logger
}

func InitLogger(cfg *config.Config) {
	file, err := os.OpenFile("/var/log/app/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	h := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: convertToSlogLevel(cfg.LoggerLevel),
	})
	logger = &Logger{
		cfg: cfg,
		log: slog.New(h),
	}
}

func convertToSlogLevel(level string) slog.Leveler {
	switch level {
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "debug":
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

type CtxLoggerKey struct{}

func WithAttrs(ctx context.Context, args ...slog.Attr) context.Context {
	ctx = context.WithValue(ctx, CtxLoggerKey{}, MergeAttrs(getAttrs(ctx), args))
	return ctx
}

func MergeAttrs(left []slog.Attr, right []slog.Attr) []slog.Attr {
	return append(left, right...)
}

func getAttrs(ctx context.Context) []slog.Attr {
	attrs := ctx.Value(CtxLoggerKey{})
	if attrs == nil {
		return []slog.Attr{}
	}
	result, ok := attrs.([]slog.Attr)
	if !ok {
		return []slog.Attr{}
	}
	return result
}

func convertAttrsToAny(a []slog.Attr) []any {
	result := make([]any, len(a))
	for i, v := range a {
		result[i] = v
	}
	return result
}

func (s *Logger) Info(ctx context.Context, msg string, args ...slog.Attr) {
	s.log.InfoContext(ctx, msg, convertAttrsToAny(MergeAttrs(getAttrs(ctx), args))...)
}

// Update these methods to correctly pass slog.Attr values
func (l *Logger) Error(msg string, args ...any) {
	// Convert args to interface{} pairs
	attrs := make([]any, len(args))
	for i := 0; i < len(args); i++ {
		if i%2 == 0 && i+1 < len(args) {
			// Create a slog.Attr for each key-value pair
			key, ok := args[i].(string)
			if !ok {
				continue
			}
			attrs[i] = key
			attrs[i+1] = args[i+1]
		}
	}
	l.log.Error(msg, attrs...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	// Convert args to interface{} pairs
	attrs := make([]any, len(args))
	for i := 0; i < len(args); i++ {
		if i%2 == 0 && i+1 < len(args) {
			// Create a slog.Attr for each key-value pair
			key, ok := args[i].(string)
			if !ok {
				continue
			}
			attrs[i] = key
			attrs[i+1] = args[i+1]
		}
	}
	l.log.ErrorContext(ctx, msg, attrs...)
}

func (s *Logger) Panic(ctx context.Context, err error, args slog.Attr) {
	s.log.ErrorContext(ctx, err.Error(), args)
	panic(err)
}

func (s *Logger) Fatal(ctx context.Context, err error, args slog.Attr) {
	s.log.ErrorContext(ctx, err.Error(), args)
	os.Exit(1)
}
