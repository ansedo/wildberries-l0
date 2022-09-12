package shutdowner

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ansedo/wildberries-l0/internal/logger"
)

const gracefulShutdownDelay = 5 * time.Second

var once sync.Once

type Shutdowner struct {
	mu           sync.RWMutex
	logger       *zap.Logger
	callbacks    []func(context.Context) error
	chAllClosed  chan struct{}
	ChShutdowned chan struct{}
}

type ctxShutdowner struct{}

func New(ctx context.Context) *Shutdowner {
	s := &Shutdowner{
		logger:       logger.FromCtx(ctx),
		chAllClosed:  make(chan struct{}),
		ChShutdowned: make(chan struct{}),
	}
	go s.catchSignalsAndShutdown(ctx)
	return s
}
func (s *Shutdowner) AddCloser(fn func(ctx context.Context) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callbacks = append(s.callbacks, fn)
}

func (s *Shutdowner) gracefulShutdown() {
	once.Do(func() {
		s.mu.RLock()
		defer s.mu.RUnlock()

		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownDelay)
		defer cancel()

		errs := make(chan error, len(s.callbacks))
		for len(s.callbacks) > 0 {
			lastIndex := len(s.callbacks) - 1
			errs <- s.callbacks[lastIndex](ctx)
			s.callbacks = s.callbacks[:lastIndex]
		}

		close(s.chAllClosed)
	})
}

func (s *Shutdowner) catchSignalsAndShutdown(ctx context.Context) {
	chStopSignalReceived := make(chan os.Signal, 1)
	signal.Notify(chStopSignalReceived, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-chStopSignalReceived

	go s.gracefulShutdown()

	select {
	case <-s.chAllClosed:
		s.logger.Info("graceful shutdown successfully done")
		close(s.ChShutdowned)
		return
	case <-time.After(2 * gracefulShutdownDelay):
		s.logger.Warn("no response while graceful shutdown: exit with error")
		os.Exit(int(syscall.SIGTERM))
	}
}

func CtxWith(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxShutdowner{}, New(ctx))
}

func FromCtx(ctx context.Context) *Shutdowner {
	if s, ok := ctx.Value(ctxShutdowner{}).(*Shutdowner); ok {
		return s
	}
	logger.FromCtx(ctx).Warn("shutdowner does not exist in context")
	return nil
}
