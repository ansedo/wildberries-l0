package handlers

import (
	"context"
	"encoding/json"
	"github.com/ansedo/wildberries-l0/internal/models"
	"net/http"

	"go.uber.org/zap"

	"github.com/ansedo/wildberries-l0/internal/logger"
	"github.com/ansedo/wildberries-l0/internal/storages"
)

func GetOrderByUID(ctx context.Context, db storages.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)

		orderUID := r.PostFormValue("order_uid")
		if orderUID == "" {
			encodeErr(w, models.ErrOrderUIDRequired)
			return
		}

		order, err := db.GetOrderByUID(ctx, orderUID)
		if err != nil {
			encodeErr(w, err)
			return
		}

		resp, err := json.Marshal(order)
		if err != nil {
			logger.FromCtx(ctx).Warn("get order by uid: json marshal", zap.Error(err))
			encodeErr(w, err)
			return
		}

		if err = json.NewEncoder(w).Encode("{\"order\":" + string(resp) + "}"); err != nil {
			logger.FromCtx(ctx).Warn("get order by uid: send response", zap.Error(err))
			encodeErr(w, err)
			return
		}
	}
}

func encodeErr(w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(`{"error":"` + err.Error() + `"}`)
}
