package handler

import "net/http"

func (h *HTTPHandler) SetupRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /crawl", h.CrawlAll)
	mux.HandleFunc("GET /jobs/active", h.GetAllActive)
	mux.HandleFunc("PATCH /jobs/{id}/inactive", h.MarkInactive)

	return mux
}
