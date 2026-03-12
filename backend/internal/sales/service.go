package sales

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateSaleRequest struct {
	Total         int64              `json:"total"`
	PaymentMethod string             `json:"payment_method"`
	Items         []CreateItemRequest `json:"items"`
}

type CreateItemRequest struct {
	ProductID   string `json:"product_id"`   // optional for custom items
	ProductName string `json:"product_name"` // required when product_id empty
	Qty         int    `json:"qty"`
	Price       int64  `json:"price"`
	Subtotal    int64  `json:"subtotal"`
}

func (s *Service) Create(ctx context.Context, tenantID uuid.UUID, req CreateSaleRequest) (*Sale, error) {
	if req.PaymentMethod == "" {
		req.PaymentMethod = "cash"
	}
	items := make([]CreateSaleItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		var productID *uuid.UUID
		if it.ProductID != "" {
			pid, err := uuid.Parse(it.ProductID)
			if err == nil {
				productID = &pid
			}
		}
		if productID == nil && it.ProductName == "" {
			it.ProductName = "Item"
		}
		items = append(items, CreateSaleItemInput{
			ProductID:   productID,
			ProductName: it.ProductName,
			Qty:         it.Qty,
			Price:       it.Price,
			Subtotal:    it.Subtotal,
		})
	}
	return s.repo.Create(ctx, tenantID, CreateSaleInput{
		Total:         req.Total,
		PaymentMethod: req.PaymentMethod,
		Items:         items,
	})
}

func (s *Service) Get(ctx context.Context, id, tenantID uuid.UUID) (*Sale, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, from, to *time.Time, limit, offset int) ([]Sale, error) {
	return s.repo.ListByTenant(ctx, tenantID, from, to, limit, offset)
}
