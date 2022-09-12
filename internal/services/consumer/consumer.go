package consumer

import (
	"context"
	"encoding/json"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"

	"github.com/ansedo/wildberries-l0/internal/config"
	"github.com/ansedo/wildberries-l0/internal/logger"
	"github.com/ansedo/wildberries-l0/internal/models"
	"github.com/ansedo/wildberries-l0/internal/services/shutdowner"
	"github.com/ansedo/wildberries-l0/internal/storages"
)

const (
	stanClientID    = "test-client-consumer"
	stanDurableName = "test-durable-name"
)

type Consumer struct {
	sc        stan.Conn
	db        storages.Storager
	logger    *zap.Logger
	ch        chan models.Order
	clusterID string
	subject   string
}

func New(ctx context.Context, db storages.Storager, cfg config.Stan) *Consumer {
	c := &Consumer{
		db:        db,
		logger:    logger.FromCtx(ctx),
		clusterID: cfg.ClusterID,
		subject:   cfg.Subject,
	}

	var err error
	c.sc, err = stan.Connect(c.clusterID, stanClientID)
	if err != nil {
		c.logger.Fatal("consumer constructor: connect to nats streaming", zap.Error(err))
	}

	c.addToShutdowner(ctx)
	return c
}

func (c *Consumer) Run(ctx context.Context) {
	if _, err := c.sc.Subscribe(c.subject, func(msg *stan.Msg) {
		c.addOrder(ctx, msg)
	}, stan.DurableName(stanDurableName)); err != nil {
		c.logger.Fatal("consumer run: subscribe to nats streaming", zap.Error(err))
	}
}

func (c *Consumer) addOrder(ctx context.Context, msg *stan.Msg) {
	var order models.Order
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		c.logger.Warn("consumer add order: unmarshal json", zap.Error(err))
	}
	order.Delivery.OrderUID = order.OrderUID
	c.db.AddOrder(ctx, order)
}

func (c *Consumer) addToShutdowner(ctx context.Context) {
	if shutdown := shutdowner.FromCtx(ctx); shutdown != nil {
		shutdown.AddCloser(func(_ context.Context) error {
			if err := c.sc.Close(); err != nil {
				return err
			}
			return nil
		})
	}
}
