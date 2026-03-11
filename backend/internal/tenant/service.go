package tenant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrPhoneExists = errors.New("phone already registered")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type RegisterInput struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterResult struct {
	TenantID string `json:"tenant_id"`
	Message  string `json:"message"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*RegisterResult, error) {
	if input.Phone == "" || input.Password == "" {
		return nil, errors.New("phone and password required")
	}
	name := input.Name
	if name == "" {
		name = input.Phone
	}
	exists, err := s.repo.ExistsByPhone(ctx, input.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPhoneExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	tenant, _, err := s.repo.Create(ctx, name, input.Phone, string(hash))
	if err != nil {
		return nil, err
	}
	return &RegisterResult{
		TenantID: tenant.ID.String(),
		Message:  "registered; login with phone and password",
	}, nil
}

type SettingsInput struct {
	Name           string `json:"name"`
	ReceiptFooter  string `json:"receipt_footer"`
	DefaultPayment string `json:"default_payment"`
	WhatsAppNumber string `json:"whatsapp_number"`
}

func (s *Service) GetSettings(ctx context.Context, tenantID uuid.UUID) (*Tenant, error) {
	return s.repo.GetByID(ctx, tenantID)
}

func (s *Service) UpdateSettings(ctx context.Context, tenantID uuid.UUID, input SettingsInput) error {
	return s.repo.UpdateSettings(ctx, tenantID, input.Name, input.ReceiptFooter, input.DefaultPayment, input.WhatsAppNumber)
}
