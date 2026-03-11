package report

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hsmart/app/backend/pkg/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) DailySummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tenantID := middleware.GetTenantID(r.Context())
	dateStr := r.URL.Query().Get("date")
	date := time.Now().UTC()
	if dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = d
		}
	}
	ds, err := h.svc.DailySummary(r.Context(), tenantID, date)
	if err != nil {
		http.Error(w, `{"error":"summary failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ds)
}

func (h *Handler) ProductRanking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tenantID := middleware.GetTenantID(r.Context())
	dateStr := r.URL.Query().Get("date")
	date := time.Now().UTC()
	if dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = d
		}
	}
	rank, err := h.svc.ProductRanking(r.Context(), tenantID, date, 10)
	if err != nil {
		http.Error(w, `{"error":"ranking failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rank)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tenantID := middleware.GetTenantID(r.Context())
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	dateStr := r.URL.Query().Get("date")
	tzStr := r.URL.Query().Get("tz")
	if tzStr == "" {
		tzStr = "Asia/Jakarta"
	}
	loc := time.UTC
	if l, err := time.LoadLocation(tzStr); err == nil {
		loc = l
	}
	if fromStr != "" && toStr != "" {
		from, err1 := time.ParseInLocation("2006-01-02", fromStr, loc)
		to, err2 := time.ParseInLocation("2006-01-02", toStr, loc)
		if err1 == nil && err2 == nil && !from.After(to) {
			dayStart := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, loc)
			dayEnd := time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, loc)
			dash, err := h.svc.DashboardRange(r.Context(), tenantID, dayStart, dayEnd, tzStr)
			if err != nil {
				http.Error(w, `{"error":"dashboard failed"}`, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(dash)
			return
		}
	}
	date := time.Now().In(loc)
	if dateStr != "" {
		if d, err := time.ParseInLocation("2006-01-02", dateStr, loc); err == nil {
			date = d
		}
	}
	// Hari mulai dan akhir di timezone user, konversi ke UTC untuk query
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc).UTC()
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, loc).UTC()
	dash, err := h.svc.DashboardForRange(r.Context(), tenantID, dayStart, dayEnd, tzStr)
	if err != nil {
		http.Error(w, `{"error":"dashboard failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dash)
}
