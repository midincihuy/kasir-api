package handlers

import (
	"kasir-api/services"
	"net/http"
	"encoding/json"

)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// HandleReportToday - Get /api/report/hari-ini
func (h *ReportHandler) HandleReportToday(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ReportToday(w, r)		
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReportHandler) ReportToday(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.GetReportToday()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// HandleReport - Get /api/report?start_date=2026-01-01&end_date=2026-02-31
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Report(w, r)		
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}
func (h *ReportHandler) Report(w http.ResponseWriter, r *http.Request) {
	start_date := r.URL.Query().Get("start_date")
	end_date := r.URL.Query().Get("end_date")
	reports, err := h.service.GetReport(start_date, end_date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}