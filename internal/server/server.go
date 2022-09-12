package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/ansedo/wildberries-l0/internal/logger"
	"github.com/ansedo/wildberries-l0/internal/router"
	"github.com/ansedo/wildberries-l0/internal/services/shutdowner"
	"github.com/ansedo/wildberries-l0/internal/storages"
)

type Server struct {
	db     storages.Storager
	http   http.Server
	logger *zap.Logger
}

func New(ctx context.Context, db storages.Storager, runAddress string) *Server {
	s := &Server{
		db: db,
		http: http.Server{
			Addr:    runAddress,
			Handler: router.New(ctx, db),
		},
		logger: logger.FromCtx(ctx),
	}
	s.addToShutdowner(ctx)
	return s
}

func (s *Server) Run(ctx context.Context) {
	go s.ListenAndServer(ctx)
}

func (s *Server) ListenAndServer(_ context.Context) {
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal(err.Error())
	}
}

func (s *Server) addToShutdowner(ctx context.Context) {
	if shutdown := shutdowner.FromCtx(ctx); shutdown != nil {
		shutdown.AddCloser(func(ctx context.Context) error {
			if err := s.http.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		})
	}
}
