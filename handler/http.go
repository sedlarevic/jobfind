package handler

import (
	"jobfind/service"
	"log"
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

func (h *HTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	result, _ := h.crawlService.RefreshJobs(r.Context())

	log.Printf("Result: %v", result)

}

func (h *HTTPHandler) GetAllActive(w http.ResponseWriter, r *http.Request) {
	result, _ := h.jobPostingService.GetAllActive(r.Context())

	log.Printf("Result: %v", result)
}
