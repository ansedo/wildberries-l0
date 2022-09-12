package logger

import (
	"context"
	"log"

	"go.uber.org/zap"
)

type ctxLogger struct{}

func New(_ context.Context) *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func CtxWith(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxLogger{}, New(ctx))
}

func FromCtx(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
		return logger
	}
	logger := zap.L()
	logger.Warn("logger does not exist in context")
	return logger
}
