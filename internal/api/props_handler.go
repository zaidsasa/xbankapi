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

// swagger:route GET /healthz isHealthy
//
// # Check HTTP Server Status
//
// This endpoint returns a HTTP 200 status code when the API server is accepting incoming
// HTTP requests. This status does currently not include checks whether the database connection is working.
//
//	Responses:
//		200: description: API server is ready to accept connections and always return "ok".
//		500: description: API server is not yet ready to accept connections.
func (h *PropsHandler) health(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte(responseOK)); err != nil {
		http.Error(w, responseNotHealthy, http.StatusInternalServerError)
	}
}

// swagger:route GET /readiness isReady
//
// # Check HTTP Server and Database Status
//
// This endpoint returns a HTTP 200 status code when the API server is up running and the environment dependencies
// (e.g. the database) are responsive as well.
//
//	Responses:
//		200: description: API server is ready to accept requests and always return "ok".
//		503: description: API server is not yet ready to accept requests.
func (h *PropsHandler) readiness(w http.ResponseWriter, r *http.Request) {
	if err := h.conn.Ping(r.Context()); err != nil {
		http.Error(w, responseNotHReady, http.StatusInternalServerError)

		return
	}

	if _, err := w.Write([]byte(responseOK)); err != nil {
		http.Error(w, responseNotHReady, http.StatusInternalServerError)
	}
}
