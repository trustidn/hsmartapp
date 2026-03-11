package saasconfig

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	c, err := h.repo.Get(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"app_name":             c.AppName,
		"logo_url":             c.LogoURL,
		"admin_contact":        c.AdminContact,
		"bank_name":            c.BankName,
		"bank_account_number":  c.BankAccountNumber,
		"bank_account_name":    c.BankAccountName,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		AppName             *string `json:"app_name"`
		LogoURL             *string `json:"logo_url"`
		AdminContact        *string `json:"admin_contact"`
		BankName            *string `json:"bank_name"`
		BankAccountNumber   *string `json:"bank_account_number"`
		BankAccountName     *string `json:"bank_account_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if err := h.repo.Update(r.Context(), input.AppName, input.LogoURL, input.AdminContact, input.BankName, input.BankAccountNumber, input.BankAccountName); err != nil {
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"ok": "updated"})
}
