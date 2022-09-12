package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/ansedo/wildberries-l0/internal/config"
	"github.com/ansedo/wildberries-l0/internal/logger"
	"github.com/ansedo/wildberries-l0/internal/server"
	"github.com/ansedo/wildberries-l0/internal/services/consumer"
	"github.com/ansedo/wildberries-l0/internal/services/producer"
	"github.com/ansedo/wildberries-l0/internal/services/shutdowner"
	"github.com/ansedo/wildberries-l0/internal/storages/postgres"
)

func main() {
	// sowing seeds
	rand.Seed(time.Now().UnixNano())

	// add logger to context
	ctx := logger.CtxWith(context.Background())

	// add shutdowner to context
	ctx = shutdowner.CtxWith(ctx)

	// construct config
	cfg := config.New(ctx)

	// construct database
	db := postgres.New(ctx, cfg.DatabaseURI)

	// construct and run server
	server.New(ctx, db, cfg.RunAddress).Run(ctx)

	// construct and run consumer
	consumer.New(ctx, db, cfg.Stan).Run(ctx)

	// construct and run producer in goroutine
	go producer.New(ctx, cfg.Stan).Run(ctx)

	// waiting for graceful shutdown if it exists in context
	if shutdown := shutdowner.FromCtx(ctx); shutdown != nil {
		<-shutdown.ChShutdowned
	}
}
