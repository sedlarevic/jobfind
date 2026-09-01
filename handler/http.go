package handler

import (
	"encoding/json"
	"jobfind/service"
	"log/slog"
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
	result, err := h.crawlService.RefreshJobs(r.Context())

	if err != nil {
		slog.Error("refresh jobs request failed", "method", r.Method, "path", r.URL.Path, "error", err)

		http.Error(w, "failed to refresh jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("failed to encode refresh response", "error", err)
		return
	}

	slog.Info(
		"refresh jobs request completed",
		"crawled", result.Crawled,
		"new", result.New,
		"updated", result.Updated,
		"deactivated", result.Deactivated)
}

func (h *HTTPHandler) GetAllActive(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobPostingService.GetAllActive(r.Context())

	if err != nil {
		slog.Error("get all active jobs request failed", "method", r.Method, "path", r.URL.Path, "error", err)

		http.Error(w, "failed to get active jobs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		slog.Error("failed to encode get all active jobs response", "error", err)
		return
	}

	slog.Info("get active jobs request completed", "jobs", len(jobs))
}
