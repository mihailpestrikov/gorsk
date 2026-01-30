package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/zlog"
)

// New creates new http server
func New(cfg *config.Config, log *zlog.Logger, r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

// Health returns an OK message
func Health(cfg *config.Config, log *zlog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
