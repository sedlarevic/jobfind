package handler

import (
	"jobfind/service"
	"net/http"
)

type HTTPHandler struct {
	jobPostingService *service.JobPostingService
	crawlService      *service.CrawlService
}

func NewHTTPHandler(jps *service.JobPostingService, cs *service.CrawlService) *HTTPHandler {
	return &HTTPHandler{
		jobPostingService: jps,
		crawlService:      cs,
	}
}

func (h *HTTPHandler) CrawlAll(w http.ResponseWriter, r *http.Request) {
}
func (h *HTTPHandler) GetAllActive(w http.ResponseWriter, r *http.Request) {
}
func (h *HTTPHandler) MarkInactive(w http.ResponseWriter, r *http.Request) {
}

func (h *HTTPHandler) SetupRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /crawl", h.CrawlAll)
	mux.HandleFunc("GET /jobs/active", h.GetAllActive)
	mux.HandleFunc("PATCH /jobs/{id}/inactive", h.MarkInactive)

	return mux
}
