package producer

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"

	"github.com/ansedo/wildberries-l0/internal/config"
	"github.com/ansedo/wildberries-l0/internal/logger"
	"github.com/ansedo/wildberries-l0/internal/models"
	"github.com/ansedo/wildberries-l0/internal/services/shutdowner"
)

const stanClientID = "test-client-producer"

type Producer struct {
	sc        stan.Conn
	logger    *zap.Logger
	clusterID string
	subject   string
	chCancel  chan struct{}
}

func New(ctx context.Context, cfg config.Stan) *Producer {
	p := &Producer{
		logger:    logger.FromCtx(ctx),
		clusterID: cfg.ClusterID,
		subject:   cfg.Subject,
		chCancel:  make(chan struct{}),
	}

	var err error
	if p.sc, err = stan.Connect(p.clusterID, stanClientID); err != nil {
		p.logger.Fatal("producer constructor: connect to nats streaming subsystem")
	}

	p.addToShutdowner(ctx)
	return p
}

func (p *Producer) Run(ctx context.Context) {
	for {
		select {
		case <-p.chCancel:
			return
		default:
			if err := p.sc.Publish(p.subject, getRandData(ctx)); err != nil {
				p.logger.Fatal("producer publish: message #", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}
}

func (p *Producer) addToShutdowner(ctx context.Context) {
	if shutdown := shutdowner.FromCtx(ctx); shutdown != nil {
		shutdown.AddCloser(func(_ context.Context) error {
			p.chCancel <- struct{}{}
			if err := p.sc.Close(); err != nil {
				return err
			}
			return nil
		})
	}
}

func getRandData(_ context.Context) []byte {
	orderUID := getRandString(19)
	trackNumber := getRandString(14)
	order := models.Order{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             getRandString(4),
		Locale:            getRandString(2),
		InternalSignature: getRandString(getRandInt(0, 19)),
		CustomerId:        getRandString(getRandInt(0, 19)),
		DeliveryService:   getRandString(getRandInt(0, 19)),
		Shardkey:          strconv.Itoa(getRandInt(0, 1e9-1)),
		SmId:              getRandInt(1, 1e8-1),
		DateCreated:       time.Now(),
		OofShard:          strconv.Itoa(getRandInt(0, 1e9-1)),
		Delivery: models.Delivery{
			Name: getRandString(
				getRandInt(5, 15)) + " " +
				getRandString(getRandInt(10, 20),
				),
			Phone: "+" + strconv.Itoa(getRandInt(1e9, 1e14-1)),
			Zip:   strconv.Itoa(getRandInt(1e6, 1e10-1)),
			City: getRandString(
				getRandInt(3, 10)) + " " +
				getRandString(getRandInt(5, 15),
				),
			Address: getRandString(
				getRandInt(5, 12)) + " " +
				getRandString(getRandInt(4, 10)) +
				strconv.Itoa(getRandInt(1, 100)),
			Region: getRandString(getRandInt(5, 15)),
			Email: getRandString(
				getRandInt(5, 15)) + "@" +
				getRandString(getRandInt(2, 10)) + ".com",
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestId:    getRandString(getRandInt(0, 19)),
			Currency:     getRandString(3),
			Provider:     getRandString(getRandInt(5, 10)),
			Amount:       getRandInt(1, 1e6),
			PaymentDt:    getRandInt(1e8, 1e9-1),
			Bank:         getRandString(getRandInt(5, 10)),
			DeliveryCost: getRandInt(100, 1e5),
			GoodsTotal:   getRandInt(1, 1e4),
			CustomFee:    getRandInt(0, 20),
		},
	}
	for i := 0; i < getRandInt(1, 30); i++ {
		price := getRandInt(100, 1e6)
		sale := getRandInt(0, 100)
		order.Items = append(order.Items, models.Item{
			ChrtId:      getRandInt(1e7, 1e9),
			TrackNumber: trackNumber,
			Price:       price,
			Rid:         getRandString(21),
			Name:        getRandString(getRandInt(5, 20)),
			Sale:        sale,
			Size:        strconv.Itoa(getRandInt(0, 10)),
			TotalPrice:  price * (100 - sale) / 100,
			NmId:        getRandInt(1e6, 1e7-1),
			Brand: getRandString(
				getRandInt(5, 10)) + " " +
				getRandString(getRandInt(2, 8),
				),
			Status: getRandInt(100, 500),
		})
	}
	bytes, _ := json.Marshal(&order)
	return bytes
}

func getRandString(n int) string {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"[rand.Intn(61)]
	}
	return string(bytes)
}

func getRandInt(min int, max int) int {
	return rand.Intn(max-min) + min
}
