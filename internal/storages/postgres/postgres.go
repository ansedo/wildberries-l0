package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/ansedo/wildberries-l0/internal/logger"
	"github.com/ansedo/wildberries-l0/internal/models"
	"github.com/ansedo/wildberries-l0/internal/services/shutdowner"
	"github.com/ansedo/wildberries-l0/internal/storages"
)

const (
	queryTimeout              = time.Second
	ordersFlushThresholdCount = 5
)

var (
	//go:embed sql/migrations/*.sql
	migrationsFs     embed.FS
	migrationsFsName = "sql/migrations"

	//go:embed sql/queries/*.sql
	queriesFS     embed.FS
	queriesFsName = "sql/queries"
)

func getQueryFromFile(filename string) (string, error) {
	query, err := queriesFS.ReadFile(fmt.Sprintf("%s/%s", queriesFsName, filename))
	if err != nil {
		return "", err
	}
	return string(query), nil
}

type Storage struct {
	db      *sqlx.DB
	cache   *models.Cache
	logger  *zap.Logger
	queries struct {
		migrations string

		selectOrders string
		selectItems  string

		selectOrderByID          string
		selectItemsByTrackNumber string

		insertOrder    string
		insertDelivery string
		insertItem     string
		insertPayment  string
	}
}

var _ storages.Storager = (*Storage)(nil)

func New(ctx context.Context, DatabaseURI string) storages.Storager {
	s := &Storage{
		logger: logger.FromCtx(ctx),
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	var err error
	if s.db, err = sqlx.ConnectContext(ctx, "pgx", DatabaseURI); err != nil {
		s.logger.Fatal("storage constructor: connect to postgres db", zap.Error(err))
	}

	ctx, cancel = context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err = s.db.PingContext(ctx); err != nil {
		s.logger.Fatal("storage constructor: verify connection to db", zap.Error(err))
	}
	if err = s.migrate(ctx); err != nil {
		s.logger.Fatal("storage constructor: applies migrations", zap.Error(err))
	}
	if err = s.setQueries(ctx); err != nil {
		s.logger.Fatal("storage constructor: set queries", zap.Error(err))
	}

	s.restoreCache(ctx)
	s.addToShutdowner(ctx)
	return s
}

func (s *Storage) AddOrder(ctx context.Context, order models.Order) {
	s.cache.AddOrder(order)
	if s.cache.NewOrderCount() > ordersFlushThresholdCount {
		s.flushCache(ctx)
	}
}

func (s *Storage) AddOrders(ctx context.Context, orders []models.Order) (err error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Warn("storage add order: transaction begin", zap.Error(err))
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			if err = tx.Rollback(); err != nil {
				s.logger.Warn("storage add order: transaction rollback", zap.Error(err))
			}
			s.logger.Fatal("storage add order: transaction panic", zap.Any("panic", p))
		case err != nil:
			s.logger.Warn("storage add order: transaction error", zap.Error(err))
			if err = tx.Rollback(); err != nil {
				s.logger.Warn("storage add order: transaction rollback", zap.Error(err))
			}
		default:
			if err = tx.Commit(); err != nil {
				s.logger.Fatal("storage add order: transaction commit", zap.Error(err))
			}
		}
	}()

	ctx, cancel = context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if _, err = tx.NamedExecContext(ctx, s.queries.insertOrder, orders); err != nil {
		return
	}

	ctx, cancel = context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if _, err = tx.NamedExecContext(ctx, s.queries.insertDelivery, orders); err != nil {
		return
	}

	ctx, cancel = context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if _, err = tx.NamedExecContext(ctx, s.queries.insertPayment, orders); err != nil {
		return
	}

	var items []models.Item
	for _, order := range orders {
		items = append(items, order.Items...)
	}
	ctx, cancel = context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	if _, err = s.db.NamedExecContext(ctx, s.queries.insertItem, items); err != nil {
		return
	}

	return
}

func (s *Storage) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	if err := s.db.SelectContext(ctx, &orders, s.queries.selectOrders); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Order{}, models.ErrOrderUIDNotExist
		}
		s.logger.Fatal("storage get orders: query", zap.Error(err))
	}

	ordersMap := make(map[string]*models.Order)
	for i := range orders {
		ordersMap[orders[i].TrackNumber] = &orders[i]
	}

	ctx, cancel = context.WithTimeout(ctx, 2*queryTimeout)
	defer cancel()
	rows, err := s.db.QueryxContext(ctx, s.queries.selectItems)
	if err != nil {
		s.logger.Fatal("storage get items: query", zap.Error(err))
	}
	for rows.Next() {
		var item models.Item
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		if p, ok := ordersMap[item.TrackNumber]; ok {
			(*p).Items = append((*p).Items, item)
		}
	}
	return orders, nil
}

func (s *Storage) GetOrderByUID(_ context.Context, orderUID string) (models.Order, error) {
	if order := s.cache.GetByOrderUID(orderUID); !order.IsEmpty() {
		return order, nil
	}
	return models.Order{}, models.ErrOrderUIDNotExist
}

func (s *Storage) GetStats(_ context.Context) models.Stats {
	return models.Stats{
		OrderCount:    s.cache.GetOrderCount(),
		ItemCount:     s.cache.GetItemCount(),
		LastOrderUIDs: s.cache.GetLastOrderUIDs(),
	}
}

func (s *Storage) flushCache(ctx context.Context) {
	if s.cache.NewOrderCount() == 0 {
		return
	}
	if err := s.AddOrders(ctx, s.cache.GetNewOrders()); err != nil {
		s.logger.Warn("flush cache: add orders", zap.Error(err))
		return
	}
	s.cache.DeleteNewOrderUIDs()
}

func (s *Storage) restoreCache(ctx context.Context) {
	s.cache = models.NewCache(ctx)
	orderUIDTrackNumberMap := map[string]string{}
	ctx, cancel := context.WithTimeout(ctx, 5*queryTimeout)
	defer cancel()
	rows, err := s.db.QueryxContext(ctx, s.queries.selectOrders)
	if err != nil {
		s.logger.Warn("storage get orders: query", zap.Error(err))
		return
	}
	for rows.Next() {
		var order models.Order
		if err = rows.StructScan(&order); err != nil {
			s.logger.Warn("storage get orders: next rows query", zap.Error(err))
			return
		}
		order.InDB = true
		s.cache.AddOrder(order)
		orderUIDTrackNumberMap[order.TrackNumber] = order.OrderUID
	}

	ctx, cancel = context.WithTimeout(ctx, 5*queryTimeout)
	defer cancel()
	rows, err = s.db.QueryxContext(ctx, s.queries.selectItems)
	if err != nil {
		s.logger.Warn("storage get items: query", zap.Error(err))
		return
	}
	for rows.Next() {
		var item models.Item
		if err = rows.StructScan(&item); err != nil {
			s.logger.Warn("storage get items: next rows query", zap.Error(err))
			return
		}
		s.cache.AddItem(orderUIDTrackNumberMap[item.TrackNumber], item)
	}
	return
}

func (s *Storage) setQueries(_ context.Context) error {
	files, err := queriesFS.ReadDir(queriesFsName)
	if err != nil {
		return err
	}

	for _, file := range files {
		query, err := getQueryFromFile(file.Name())
		if err != nil {
			return err
		}

		switch file.Name() {
		case "insert_order.sql":
			s.queries.insertOrder = query
		case "insert_delivery.sql":
			s.queries.insertDelivery = query
		case "insert_payment.sql":
			s.queries.insertPayment = query
		case "insert_item.sql":
			s.queries.insertItem = query

		case "select_orders.sql":
			s.queries.selectOrders = query
		case "select_items.sql":
			s.queries.selectItems = query
		case "select_order_by_uid.sql":
			s.queries.selectOrderByID = query
		case "select_items_by_track_number.sql":
			s.queries.selectItemsByTrackNumber = query
		}
	}
	return err
}

func (s *Storage) migrate(_ context.Context) (err error) {
	goose.SetBaseFS(migrationsFs)
	if err = goose.SetDialect("postgres"); err != nil {
		return
	}
	err = goose.Up(s.db.DB, migrationsFsName)
	return
}

func (s *Storage) addToShutdowner(ctx context.Context) {
	if shutdown := shutdowner.FromCtx(ctx); shutdown != nil {
		shutdown.AddCloser(func(_ context.Context) error {
			s.flushCache(context.TODO())
			s.db.Close()
			return nil
		})
	}
}
