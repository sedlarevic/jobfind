package handler

import "net/http"

func (h *HTTPHandler) SetupRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /jobs/refresh", h.Refresh)
	mux.HandleFunc("GET /jobs/active", h.GetAllActive)

	return mux
}
