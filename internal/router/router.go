package router

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ansedo/wildberries-l0/internal/handlers"
	"github.com/ansedo/wildberries-l0/internal/storages"
)

func New(ctx context.Context, db storages.Storager) chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Compress(5),
	)

	r.Get("/", handlers.Index(ctx))
	r.Get("/getStats", handlers.GetStats(ctx, db))
	r.Post("/getOrderByUID", handlers.GetOrderByUID(ctx, db))

	return r
}
