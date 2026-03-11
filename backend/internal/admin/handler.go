package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/hsmart/app/backend/internal/subscription"
	"github.com/hsmart/app/backend/internal/tenant"
)

type Handler struct {
	tenantRepo       *tenant.Repository
	subscriptionRepo *subscription.Repository
}

func NewHandler(tenantRepo *tenant.Repository, subscriptionRepo *subscription.Repository) *Handler {
	return &Handler{
		tenantRepo:       tenantRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (h *Handler) ListTenants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	limit := 20
	if l := q.Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	offset := 0
	if o := q.Get("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}
	search := q.Get("search")

	result, err := h.tenantRepo.List(r.Context(), limit, offset, search)
	if err != nil {
		http.Error(w, `{"error":"failed to list tenants"}`, http.StatusInternalServerError)
		return
	}
	tenants := make([]map[string]interface{}, len(result.Tenants))
	for i, t := range result.Tenants {
		tenants[i] = map[string]interface{}{
			"id": t.ID.String(), "name": t.Name, "phone": t.Phone,
			"plan": t.Plan, "status": t.Status, "created_at": t.CreatedAt,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tenants": tenants,
		"total":   result.Total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *Handler) GetTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"id required"}`, http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	t, err := h.tenantRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
		return
	}
	sub, _ := h.subscriptionRepo.GetByTenant(r.Context(), id)
	resp := map[string]interface{}{
		"id":               t.ID.String(),
		"name":             t.Name,
		"phone":            t.Phone,
		"plan":             t.Plan,
		"status":           t.Status,
		"created_at":       t.CreatedAt,
		"receipt_footer":   t.ReceiptFooter,
		"default_payment":  t.DefaultPayment,
		"whatsapp_number":  t.WhatsAppNumber,
		"logo_url":         t.LogoURL,
	}
	if sub != nil {
		resp["subscription"] = map[string]interface{}{
			"plan":       sub.Plan,
			"status":    sub.Status,
			"started_at": sub.StartedAt,
			"expired_at": sub.ExpiredAt,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) UpdateTenantStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.ID == "" || input.Status == "" {
		http.Error(w, `{"error":"id and status required"}`, http.StatusBadRequest)
		return
	}
	if input.Status != "active" && input.Status != "suspended" && input.Status != "inactive" {
		http.Error(w, `{"error":"invalid status"}`, http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(input.ID)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	if err := h.tenantRepo.UpdateStatus(r.Context(), id, input.Status); err != nil {
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": input.Status})
}
