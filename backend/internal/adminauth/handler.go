package adminauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hsmart/app/backend/pkg/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	aid, ok := middleware.GetAdminID(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	profile, err := h.svc.GetProfile(r.Context(), aid)
	if err != nil {
		http.Error(w, `{"error":"failed to get profile"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	aid, ok := middleware.GetAdminID(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	var input UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if (input.Name == nil || strings.TrimSpace(*input.Name) == "") &&
		(input.Email == nil || strings.TrimSpace(*input.Email) == "") &&
		(input.NewPassword == nil || strings.TrimSpace(*input.NewPassword) == "") {
		http.Error(w, `{"error":"name, email, or new_password required"}`, http.StatusBadRequest)
		return
	}
	err := h.svc.UpdateProfile(r.Context(), aid, input)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			http.Error(w, `{"error":"email already in use"}`, http.StatusConflict)
			return
		}
		if errors.Is(err, ErrInvalidPassword) {
			http.Error(w, `{"error":"current password invalid"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	profile, _ := h.svc.GetProfile(r.Context(), aid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	if input.Email == "" || input.Password == "" {
		http.Error(w, `{"error":"email and password required"}`, http.StatusBadRequest)
		return
	}
	result, err := h.svc.Login(r.Context(), input)
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
