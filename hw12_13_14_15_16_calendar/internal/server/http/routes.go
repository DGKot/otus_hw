package internalhttp

import (
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.homeHandler)

	return h.LoggingMiddleware(mux)
}

func (h *Handler) homeHandler(_ http.ResponseWriter, _ *http.Request) {
	h.logger.Info("Welcome to calendar")
}
