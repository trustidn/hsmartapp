package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Max sizes
const (
	MaxLogoSize       = 2 << 20  // 2MB
	MaxPaymentProofSize = 5 << 20 // 5MB
)

type Handler struct {
	uploadDir string
	baseURL   string
}

func NewHandler(uploadDir, baseURL string) *Handler {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	return &Handler{uploadDir: uploadDir, baseURL: baseURL}
}

func (h *Handler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"file required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()
	if header.Size > MaxLogoSize {
		http.Error(w, `{"error":"file too large (max 2MB)"}`, http.StatusBadRequest)
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".svg" && ext != ".webp" {
		ext = ".png"
	}
	name := uuid.New().String() + ext
	dir := filepath.Join(h.uploadDir, "logos")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[upload] mkdir logos: %v", err)
		http.Error(w, `{"error":"upload failed"}`, http.StatusInternalServerError)
		return
	}
	path := filepath.Join(dir, name)
	dst, err := os.Create(path)
	if err != nil {
		log.Printf("[upload] create file: %v", err)
		http.Error(w, `{"error":"upload failed"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(path)
		log.Printf("[upload] copy: %v", err)
		http.Error(w, `{"error":"upload failed"}`, http.StatusInternalServerError)
		return
	}
	url := h.buildURL("/uploads/logos/" + name)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"url":"%s"}`, url)))
}

func (h *Handler) UploadPaymentProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"file required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()
	if header.Size > MaxPaymentProofSize {
		http.Error(w, `{"error":"file too large (max 5MB)"}`, http.StatusBadRequest)
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" && ext != ".pdf" {
		ext = ".jpg"
	}
	name := uuid.New().String() + ext
	dir := filepath.Join(h.uploadDir, "payment-proofs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[upload] mkdir payment-proofs: %v", err)
		http.Error(w, `{"error":"upload failed"}`, http.StatusInternalServerError)
		return
	}
	path := filepath.Join(dir, name)
	dst, err := os.Create(path)
	if err != nil {
		log.Printf("[upload] create file: %v", err)
		http.Error(w, `{"error":"upload failed"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(path)
		log.Printf("[upload] copy: %v", err)
		http.Error(w, `{"error":"upload failed"}`, http.StatusInternalServerError)
		return
	}
	url := h.buildURL("/uploads/payment-proofs/" + name)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"url":"%s"}`, url)))
}

func (h *Handler) buildURL(path string) string {
	if h.baseURL != "" {
		return strings.TrimSuffix(h.baseURL, "/") + path
	}
	return path
}

// ServeUploads serves static files from uploadDir at /uploads/
func ServeUploads(uploadDir string) http.Handler {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	abs, _ := filepath.Abs(uploadDir)
	// StripPrefix handles /uploads/ and /uploads (redirect) - use /uploads for broader match
	return http.StripPrefix("/uploads", http.FileServer(http.Dir(filepath.Join(abs))))
}
