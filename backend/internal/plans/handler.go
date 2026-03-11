package plans

import (
	"encoding/json"
	"net/http"

	"github.com/hsmart/app/backend/internal/planconfig"
)

type Handler struct {
	repo *planconfig.Repository
}

func NewHandler(repo *planconfig.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	list, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed"}`, http.StatusInternalServerError)
		return
	}
	active := make([]map[string]interface{}, 0)
	for _, p := range list {
		if !p.IsActive {
			continue
		}
		active = append(active, map[string]interface{}{
			"plan_slug":       p.PlanSlug,
			"name":            p.Name,
			"duration_days":   p.DurationDays,
			"price_rupiah":    p.PriceRupiah,
			"max_products":    p.MaxProducts,
			"report_days":     p.ReportDays,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"plans": active})
}
