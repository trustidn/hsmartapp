package admin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hsmart/app/backend/internal/auth"
	"github.com/hsmart/app/backend/internal/planconfig"
	"github.com/hsmart/app/backend/internal/sales"
	"github.com/hsmart/app/backend/internal/suborder"
	"github.com/hsmart/app/backend/internal/subscription"
	"github.com/hsmart/app/backend/internal/tenant"
	"github.com/hsmart/app/backend/pkg/middleware"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	tenantRepo       *tenant.Repository
	subscriptionRepo *subscription.Repository
	planConfigRepo   *planconfig.Repository
	salesRepo        *sales.Repository
	suborderRepo     *suborder.Repository
	authRepo         *auth.Repository
}

func NewHandler(tenantRepo *tenant.Repository, subscriptionRepo *subscription.Repository, planConfigRepo *planconfig.Repository, salesRepo *sales.Repository, suborderRepo *suborder.Repository, authRepo *auth.Repository) *Handler {
	return &Handler{
		tenantRepo:       tenantRepo,
		subscriptionRepo: subscriptionRepo,
		planConfigRepo:   planConfigRepo,
		salesRepo:        salesRepo,
		suborderRepo:     suborderRepo,
		authRepo:         authRepo,
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
	tenantIDs := make([]uuid.UUID, len(result.Tenants))
	for i, t := range result.Tenants {
		tenantIDs[i] = t.ID
	}
	stats, _ := h.salesRepo.GetStatsByTenantIDs(r.Context(), tenantIDs)
	tenants := make([]map[string]interface{}, len(result.Tenants))
	for i, t := range result.Tenants {
		m := map[string]interface{}{
			"id": t.ID.String(), "name": t.Name, "phone": t.Phone,
			"plan": t.Plan, "status": t.Status, "created_at": t.CreatedAt,
			"transactions": int64(0), "last_transaction_at": nil,
		}
		if s, ok := stats[t.ID]; ok {
			m["transactions"] = s.Count
			m["last_transaction_at"] = s.LastSaleAt
		}
		tenants[i] = m
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tenants": tenants,
		"total":   result.Total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *Handler) DashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	activeTenants, err := h.tenantRepo.CountActive(ctx)
	if err != nil {
		http.Error(w, `{"error":"failed to get stats"}`, http.StatusInternalServerError)
		return
	}
	totalOrders, totalRevenue, err := h.suborderRepo.Stats(ctx)
	if err != nil {
		http.Error(w, `{"error":"failed to get stats"}`, http.StatusInternalServerError)
		return
	}
	mrr, err := h.suborderRepo.MRR(ctx)
	if err != nil {
		mrr = 0
	}
	months := 6
	tenantGrowth, _ := h.tenantRepo.CountByMonth(ctx, months)
	revenueByMonth, _ := h.suborderRepo.RevenueByMonth(ctx, months)
	resp := map[string]interface{}{
		"active_tenants":    activeTenants,
		"total_orders":      totalOrders,
		"total_revenue":     totalRevenue,
		"mrr":               mrr,
		"tenant_growth":     tenantGrowth,
		"revenue_by_month":  revenueByMonth,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
		expiredAt := sub.ExpiredAt
		if expiredAt == nil && sub.Plan != "" {
			cfg, err := h.planConfigRepo.GetByPlan(r.Context(), sub.Plan)
			if err == nil && cfg.DurationDays > 0 && sub.StartedAt != "" {
				tStart, parseErr := time.Parse(time.RFC3339, sub.StartedAt)
				if parseErr != nil {
					tStart, parseErr = time.Parse("2006-01-02 15:04:05.999999-07", sub.StartedAt)
				}
				if parseErr != nil {
					tStart, parseErr = time.Parse("2006-01-02 15:04:05", sub.StartedAt)
				}
				if parseErr == nil {
					tExp := tStart.UTC().AddDate(0, 0, cfg.DurationDays)
					s := tExp.Format(time.RFC3339)
					expiredAt = &s
					_ = h.subscriptionRepo.UpdateLatest(r.Context(), id, "", "", expiredAt)
				}
			}
		}
		resp["subscription"] = map[string]interface{}{
			"plan":       sub.Plan,
			"status":    sub.Status,
			"started_at": sub.StartedAt,
			"expired_at": expiredAt,
		}
	}
	history, _ := h.subscriptionRepo.ListByTenant(r.Context(), id)
	subs := make([]map[string]interface{}, len(history))
	for i, s := range history {
		subs[i] = map[string]interface{}{"plan": s.Plan, "status": s.Status, "started_at": s.StartedAt, "expired_at": s.ExpiredAt}
	}
	resp["subscription_history"] = subs
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

func (h *Handler) ResetTenantPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID          string `json:"id"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.ID == "" || strings.TrimSpace(input.NewPassword) == "" {
		http.Error(w, `{"error":"id and new_password required"}`, http.StatusBadRequest)
		return
	}
	if len(input.NewPassword) < 6 {
		http.Error(w, `{"error":"password minimal 6 karakter"}`, http.StatusBadRequest)
		return
	}
	tenantID, err := uuid.Parse(input.ID)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	// Verify tenant exists
	if _, err := h.tenantRepo.GetByID(r.Context(), tenantID); err != nil {
		http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	if err := h.authRepo.UpdateOwnerPasswordByTenantID(r.Context(), tenantID, string(hash)); err != nil {
		http.Error(w, `{"error":"failed to reset password"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "password reset successfully"})
}

func (h *Handler) UpdateTenantSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID        string  `json:"id"`
		Plan      string  `json:"plan"`
		ExpiredAt *string `json:"expired_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.ID == "" {
		http.Error(w, `{"error":"id required"}`, http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(input.ID)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	if _, err := h.tenantRepo.GetByID(r.Context(), id); err != nil {
		http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
		return
	}

	if input.ExpiredAt != nil {
		if err := h.subscriptionRepo.UpdateLatest(r.Context(), id, "", "", input.ExpiredAt); err != nil {
			http.Error(w, `{"error":"gagal mengubah tanggal kadaluarsa"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "updated"})
		return
	}

	if input.Plan == "" {
		http.Error(w, `{"error":"plan required untuk menambah langganan"}`, http.StatusBadRequest)
		return
	}
	validPlans := map[string]bool{"free": true, "premium": true, "premium_1m": true, "premium_3m": true, "premium_6m": true, "premium_1y": true, "platinum": true}
	if !validPlans[input.Plan] {
		http.Error(w, `{"error":"invalid plan"}`, http.StatusBadRequest)
		return
	}
	cfg, err := h.planConfigRepo.GetByPlan(r.Context(), input.Plan)
	if err != nil || !cfg.IsActive {
		http.Error(w, `{"error":"plan tidak aktif"}`, http.StatusBadRequest)
		return
	}
	if err := h.addSubscriptionWithAccumulation(r.Context(), id, input.Plan); err != nil {
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "updated"})
}

func (h *Handler) RevokeTenantSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.ID == "" {
		http.Error(w, `{"error":"id required"}`, http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(input.ID)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	if _, err := h.tenantRepo.GetByID(r.Context(), id); err != nil {
		http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
		return
	}
	deleted, err := h.subscriptionRepo.DeleteLatestPaid(r.Context(), id)
	if err != nil {
		log.Printf("[admin] RevokeTenantSubscription DeleteLatestPaid error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "gagal mencabut langganan", "detail": err.Error()})
		return
	}
	if !deleted {
		http.Error(w, `{"error":"tidak ada langganan berbayar untuk dicabut"}`, http.StatusBadRequest)
		return
	}
	cur, _ := h.subscriptionRepo.GetByTenant(r.Context(), id)
	newPlan := "free"
	if cur != nil && cur.Plan != "" {
		newPlan = cur.Plan
	}
	if err := h.tenantRepo.UpdatePlan(r.Context(), id, newPlan); err != nil {
		http.Error(w, `{"error":"gagal memperbarui plan"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "revoked"})
}

// addSubscriptionWithAccumulation inserts new subscription row; masa terakumulasi dengan langganan sebelumnya.
func (h *Handler) addSubscriptionWithAccumulation(ctx context.Context, tenantID uuid.UUID, plan string) error {
	cfg, err := h.planConfigRepo.GetByPlan(ctx, plan)
	if err != nil {
		return err
	}
	var expiredAt *string
	if cfg.DurationDays > 0 {
		base := time.Now().UTC()
		cur, _ := h.subscriptionRepo.GetByTenant(ctx, tenantID)
		if cur != nil && cur.ExpiredAt != nil && *cur.ExpiredAt != "" {
			var tExp time.Time
			if t, e := time.Parse(time.RFC3339, *cur.ExpiredAt); e == nil {
				tExp = t
			} else if t, e := time.Parse("2006-01-02 15:04:05.999999-07", *cur.ExpiredAt); e == nil {
				tExp = t
			} else if t, e := time.Parse("2006-01-02 15:04:05", *cur.ExpiredAt); e == nil {
				tExp = t
			}
			if !tExp.IsZero() && tExp.After(time.Now().UTC()) {
				base = tExp.UTC()
			}
		}
		tNew := base.AddDate(0, 0, cfg.DurationDays)
		s := tNew.Format(time.RFC3339)
		expiredAt = &s
	} else {
		empty := ""
		expiredAt = &empty
	}
	if err := h.subscriptionRepo.Create(ctx, tenantID, plan, "active", expiredAt); err != nil {
		return err
	}
	return h.tenantRepo.UpdatePlan(ctx, tenantID, plan)
}

func (h *Handler) ListPlanConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	list, err := h.planConfigRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to list plans"}`, http.StatusInternalServerError)
		return
	}
	plans := make([]map[string]interface{}, len(list))
	for i, p := range list {
		plans[i] = map[string]interface{}{
			"plan_slug":       p.PlanSlug,
			"name":            p.Name,
			"duration_months": p.DurationMonths,
			"duration_days":   p.DurationDays,
			"max_products":    p.MaxProducts,
			"report_days":     p.ReportDays,
			"price_rupiah":    p.PriceRupiah,
			"sort_order":      p.SortOrder,
			"is_active":       p.IsActive,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"plans": plans})
}

func (h *Handler) UpdatePlanConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		PlanSlug     string  `json:"plan_slug"`
		SortOrder    *int    `json:"sort_order"`
		DurationDays *int    `json:"duration_days"`
		MaxProducts  *int    `json:"max_products"`
		ReportDays   *int    `json:"report_days"`
		PriceRupiah  *int64  `json:"price_rupiah"`
		IsActive     *bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.PlanSlug == "" {
		http.Error(w, `{"error":"plan_slug required"}`, http.StatusBadRequest)
		return
	}
	if input.MaxProducts != nil && *input.MaxProducts < -1 {
		http.Error(w, `{"error":"max_products must be >= -1"}`, http.StatusBadRequest)
		return
	}
	if input.ReportDays != nil && *input.ReportDays < -1 {
		http.Error(w, `{"error":"report_days must be >= -1"}`, http.StatusBadRequest)
		return
	}
	if input.DurationDays != nil && *input.DurationDays < -1 {
		http.Error(w, `{"error":"duration_days must be >= -1"}`, http.StatusBadRequest)
		return
	}
	if input.PriceRupiah != nil && *input.PriceRupiah < 0 {
		http.Error(w, `{"error":"price_rupiah must be >= 0"}`, http.StatusBadRequest)
		return
	}
	if err := h.planConfigRepo.Update(r.Context(), input.PlanSlug, input.MaxProducts, input.ReportDays, input.DurationDays, input.SortOrder, input.PriceRupiah, input.IsActive); err != nil {
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "updated"})
}

func (h *Handler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	planSlug := r.URL.Query().Get("plan_slug")
	if planSlug == "" {
		http.Error(w, `{"error":"plan_slug required"}`, http.StatusBadRequest)
		return
	}
	if err := h.planConfigRepo.Delete(r.Context(), planSlug); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "deleted"})
}

func (h *Handler) RestorePlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		PlanSlug string `json:"plan_slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.PlanSlug == "" {
		http.Error(w, `{"error":"plan_slug required"}`, http.StatusBadRequest)
		return
	}
	if err := h.planConfigRepo.Restore(r.Context(), input.PlanSlug); err != nil {
		http.Error(w, `{"error":"restore failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "restored"})
}

func (h *Handler) ListSubscriptionOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	limit := 50
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
	status := q.Get("status")
	orders, total, err := h.suborderRepo.ListAll(r.Context(), limit, offset, status)
	if err != nil {
		log.Printf("[admin] ListSubscriptionOrders: %v", err)
		http.Error(w, `{"error":"failed to list orders"}`, http.StatusInternalServerError)
		return
	}
	tenantIDs := make(map[uuid.UUID]bool)
	for _, o := range orders {
		tenantIDs[o.TenantID] = true
	}
	tenantNames := make(map[uuid.UUID]string)
	for id := range tenantIDs {
		if t, err := h.tenantRepo.GetByID(r.Context(), id); err == nil {
			tenantNames[id] = t.Name
		}
	}
	out := make([]map[string]interface{}, len(orders))
	for i, o := range orders {
		out[i] = map[string]interface{}{
			"id":                o.ID.String(),
			"tenant_id":         o.TenantID.String(),
			"tenant_name":       tenantNames[o.TenantID],
			"plan_slug":         o.PlanSlug,
			"amount_rupiah":     o.AmountRupiah,
			"status":            o.Status,
			"payment_note":      o.PaymentNote,
			"payment_proof_url": o.PaymentProofURL,
			"created_at":        o.CreatedAt,
			"approved_at":       o.ApprovedAt,
			"rejection_reason":  o.RejectionReason,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"orders": out, "total": total, "limit": limit, "offset": offset})
}

func (h *Handler) ApproveSubscriptionOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		OrderID string `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.OrderID == "" {
		http.Error(w, `{"error":"order_id required"}`, http.StatusBadRequest)
		return
	}
	orderID, err := uuid.Parse(input.OrderID)
	if err != nil {
		http.Error(w, `{"error":"invalid order_id"}`, http.StatusBadRequest)
		return
	}
	order, err := h.suborderRepo.GetByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, `{"error":"order not found"}`, http.StatusNotFound)
		return
	}
	if order.Status != "pending" && order.Status != "paid" {
		http.Error(w, `{"error":"order sudah diproses"}`, http.StatusBadRequest)
		return
	}
	adminID, ok := middleware.GetAdminID(r.Context())
	if !ok {
		adminID = uuid.Nil
	}
	if err := h.addSubscriptionWithAccumulation(r.Context(), order.TenantID, order.PlanSlug); err != nil {
		http.Error(w, `{"error":"gagal menambah subscription"}`, http.StatusInternalServerError)
		return
	}
	if err := h.suborderRepo.Approve(r.Context(), orderID, adminID); err != nil {
		http.Error(w, `{"error":"approve failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "approved"})
}

func (h *Handler) RejectSubscriptionOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		OrderID string `json:"order_id"`
		Reason  string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.OrderID == "" {
		http.Error(w, `{"error":"order_id required"}`, http.StatusBadRequest)
		return
	}
	orderID, err := uuid.Parse(input.OrderID)
	if err != nil {
		http.Error(w, `{"error":"invalid order_id"}`, http.StatusBadRequest)
		return
	}
	if err := h.suborderRepo.Reject(r.Context(), orderID, input.Reason); err != nil {
		http.Error(w, `{"error":"reject failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "rejected"})
}
