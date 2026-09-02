package handler

import (
	"encoding/json"
	"io"
	"jobfind/service"
	"log/slog"
	"net/http"
	"os"
)

type HTTPHandler struct {
	jobPostingService *service.JobPostingService
	crawlService      *service.CrawlService
	cvService         *service.CVService
}

func NewHTTPHandler(jps *service.JobPostingService, cs *service.CrawlService, cvs *service.CVService) *HTTPHandler {
	return &HTTPHandler{
		jobPostingService: jps,
		crawlService:      cs,
		cvService:         cvs,
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

func (h *HTTPHandler) CrawlCompany(w http.ResponseWriter, r *http.Request) {
	company := r.PathValue("company")

	result, err := h.crawlService.CrawlCompany(r.Context(), company)

	if err != nil {
		// TODO: Error could also be that request is sent for company that doesn't exist. That should be logged.
		slog.Error("company crawl failed", "error", err)
		http.Error(w, "company crawl failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("failed to encode crawled company jobs", "error", err)
		return
	}

	slog.Info("company crawl completed", "company", company, "jobs", len(result.JobPostings))
}

func (h *HTTPHandler) ExtractCV(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	// NOTE: REQUEST FILE
	requestFile, _, err := r.FormFile("cv")

	if err != nil {
		slog.Error("failed to capture uploaded cv", "error", err)
		http.Error(w, "failed to process cv", http.StatusBadRequest)
		return
	}
	defer requestFile.Close()

	// NOTE: TEMPORARY FILE
	tempFile, err := os.CreateTemp("", "jobfind-cv-*.pdf")

	if err != nil {
		slog.Error("failed to create temporary file", "error", err)
		http.Error(w, "failed to process cv", http.StatusInternalServerError)
		return
	}

	tempPath := tempFile.Name()

	defer os.Remove(tempPath)
	defer tempFile.Close()

	// NOTE: COPYING REQUEST FILE TO TEMP PATH
	if _, err := io.Copy(tempFile, requestFile); err != nil {
		slog.Error("failed to copy request file to temporary path", "error", err)
		http.Error(w, "failed to process cv", http.StatusInternalServerError)
		return
	}

	if err := tempFile.Close(); err != nil {
		slog.Error("failed to close temporary file handle", "error", err)
		http.Error(w, "failed to process cv", http.StatusInternalServerError)
		return
	}

	// NOTE: EXTRACTING TEXT FROM PDF
	extractedText, err := h.cvService.ExtractText(tempPath)
	if err != nil {
		slog.Error("failed to extract cv text", "error", err)
		http.Error(w, "failed to process cv", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"text": extractedText,
	}); err != nil {
		slog.Error("failed to encode cv response", "error", err)
		return
	}
	slog.Info("cv extract success", "chars", len(extractedText))
}
