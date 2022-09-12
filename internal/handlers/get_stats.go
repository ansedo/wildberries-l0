package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ansedo/wildberries-l0/internal/storages"
)

func GetStats(ctx context.Context, db storages.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("X-Content-Type-Options", "nosniff")

		stats := db.GetStats(ctx)
		resp, err := json.Marshal(&stats)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("{error:" + err.Error() + "}")
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = fmt.Fprint(w, string(resp)); err != nil {
			json.NewEncoder(w).Encode("{error:" + err.Error() + "}")
			return
		}
	}
}
