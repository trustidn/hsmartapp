package suborder

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hsmart/app/backend/internal/planconfig"
	"github.com/hsmart/app/backend/pkg/middleware"
)

type Handler struct {
	svc          *Service
	planConfigRepo *planconfig.Repository
}

func NewHandler(svc *Service, planConfigRepo *planconfig.Repository) *Handler {
	return &Handler{svc: svc, planConfigRepo: planConfigRepo}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	tenantID := middleware.GetTenantID(r.Context())
	var input struct {
		PlanSlug    string `json:"plan_slug"`
		PaymentNote string `json:"payment_note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.PlanSlug == "" {
		http.Error(w, `{"error":"plan_slug required"}`, http.StatusBadRequest)
		return
	}
	cfg, err := h.planConfigRepo.GetByPlan(r.Context(), input.PlanSlug)
	if err != nil || !cfg.IsActive {
		http.Error(w, `{"error":"plan tidak valid atau tidak aktif"}`, http.StatusBadRequest)
		return
	}
	amount := cfg.PriceRupiah
	if amount < 0 {
		amount = 0
	}
	order, err := h.svc.Create(r.Context(), tenantID, input.PlanSlug, amount, input.PaymentNote)
	if err != nil {
		http.Error(w, `{"error":"gagal membuat order"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) SetPaymentProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	tenantID := middleware.GetTenantID(r.Context())
	var input struct {
		OrderID       string `json:"order_id"`
		PaymentProofURL string `json:"payment_proof_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.OrderID == "" || input.PaymentProofURL == "" {
		http.Error(w, `{"error":"order_id and payment_proof_url required"}`, http.StatusBadRequest)
		return
	}
	orderID, err := uuid.Parse(input.OrderID)
	if err != nil {
		http.Error(w, `{"error":"invalid order_id"}`, http.StatusBadRequest)
		return
	}
	if err := h.svc.SetPaymentProof(r.Context(), orderID, tenantID, input.PaymentProofURL); err != nil {
		http.Error(w, `{"error":"gagal mengupload bukti atau order tidak valid"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "updated"})
}

func (h *Handler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	tenantID := middleware.GetTenantID(r.Context())
	orders, err := h.svc.ListByTenant(r.Context(), tenantID)
	if err != nil {
		http.Error(w, `{"error":"gagal memuat orders"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"orders": orders})
}
