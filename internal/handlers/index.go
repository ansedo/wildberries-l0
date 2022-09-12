package handlers

import (
	"context"
	"net/http"
)

func Index(_ context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index.html")
	}
}
