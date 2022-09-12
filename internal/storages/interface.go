package storages

import (
	"context"

	"github.com/ansedo/wildberries-l0/internal/models"
)

// Storager is the common interface implemented by all storages.
type Storager interface {
	AddOrder(ctx context.Context, order models.Order)
	GetStats(ctx context.Context) models.Stats
	GetOrderByUID(ctx context.Context, id string) (models.Order, error)
}
