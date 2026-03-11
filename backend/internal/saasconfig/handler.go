package saasconfig

import (
	"encoding/json"
	"net/http"
	"strings"
)

func originFromRequest(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" {
		scheme = "http"
	}
	if s := r.Header.Get("X-Forwarded-Proto"); s != "" {
		scheme = s
	}
	host := r.Host
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		host = h
	}
	return scheme + "://" + host
}

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

// GetPublic returns app_name and logo_url only. No auth required (untuk welcome/login page).
func (h *Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	c, err := h.repo.Get(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"app_name": "HSmart", "logo_url": ""})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"app_name": c.AppName,
		"logo_url": c.LogoURL,
	})
}

// GetManifest returns PWA manifest JSON with app name and logo. No auth required.
func (h *Handler) GetManifest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	c, err := h.repo.Get(r.Context())
	if err != nil {
		h.writeDefaultManifest(w, r)
		return
	}
	appName := c.AppName
	if appName == "" {
		appName = "HSmart"
	}
	iconSrc := c.LogoURL
	if iconSrc == "" {
		iconSrc = "/favicon.svg"
	}
	origin := originFromRequest(r)
	if !strings.HasPrefix(iconSrc, "http://") && !strings.HasPrefix(iconSrc, "https://") {
		iconSrc = origin + "/" + strings.TrimPrefix(iconSrc, "/")
	}
	iconType := "image/png"
	if strings.HasSuffix(strings.ToLower(iconSrc), ".svg") {
		iconType = "image/svg+xml"
	} else if strings.HasSuffix(strings.ToLower(iconSrc), ".webp") {
		iconType = "image/webp"
	}
	icons := []map[string]interface{}{
		{"src": iconSrc, "sizes": "192x192", "type": iconType, "purpose": "any"},
		{"src": iconSrc, "sizes": "512x512", "type": iconType, "purpose": "any maskable"},
	}
	if iconType == "image/svg+xml" {
		icons = []map[string]interface{}{
			{"src": iconSrc, "sizes": "any", "type": iconType, "purpose": "any maskable"},
		}
	}
	manifest := map[string]interface{}{
		"name":             appName,
		"short_name":       appName,
		"description":     "POS & bisnis untuk UMKM",
		"start_url":        origin + "/",
		"display":          "standalone",
		"theme_color":      "#16a34a",
		"background_color": "#ffffff",
		"orientation":      "portrait",
		"icons":            icons,
	}
	w.Header().Set("Content-Type", "application/manifest+json")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate") // selalu fresh agar logo admin langsung teraplikasi
	json.NewEncoder(w).Encode(manifest)
}

func (h *Handler) writeDefaultManifest(w http.ResponseWriter, r *http.Request) {
	origin := originFromRequest(r)
	iconSrc := origin + "/favicon.svg"
	manifest := map[string]interface{}{
		"name":             "HSmart",
		"short_name":       "HSmart",
		"description":     "POS & bisnis untuk UMKM",
		"start_url":        origin + "/",
		"display":          "standalone",
		"theme_color":      "#16a34a",
		"background_color": "#ffffff",
		"orientation":      "portrait",
		"icons": []map[string]interface{}{
			{"src": iconSrc, "sizes": "any", "type": "image/svg+xml", "purpose": "any maskable"},
		},
	}
	w.Header().Set("Content-Type", "application/manifest+json")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	json.NewEncoder(w).Encode(manifest)
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
