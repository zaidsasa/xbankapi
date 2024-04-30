package api

import (
	"net/http"

	"github.com/zaidsasa/xbankapi/internal/storage"
)

const (
	routeHealth    = "GET /healthz"
	routeReadiness = "GET /readiness"

	responseOK         = "OK"
	responseNotHealthy = "not healthy"
	responseNotHReady  = "not ready"
)

type PropsHandler struct {
	conn storage.DBConnection
}

// NewPropsHandler returns a new PropsHandler.
func NewPropsHandler(conn storage.DBConnection) *PropsHandler {
	return &PropsHandler{
		conn: conn,
	}
}

// Register routes.
func (h *PropsHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(routeHealth, h.health)
	mux.HandleFunc(routeReadiness, h.readiness)
}

func (h *PropsHandler) health(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte(responseOK)); err != nil {
		http.Error(w, responseNotHealthy, http.StatusInternalServerError)
	}
}

func (h *PropsHandler) readiness(w http.ResponseWriter, r *http.Request) {
	if err := h.conn.Ping(r.Context()); err != nil {
		http.Error(w, responseNotHReady, http.StatusInternalServerError)

		return
	}

	if _, err := w.Write([]byte(responseOK)); err != nil {
		http.Error(w, responseNotHReady, http.StatusInternalServerError)
	}
}
