package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"uid"`
	TenantID string `json:"tid"`
	Role     string `json:"role"`
}

type Service struct {
	repo   *Repository
	secret []byte
}

func NewService(repo *Repository, secret []byte) *Service {
	return &Service{repo: repo, secret: secret}
}

type LoginInput struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResult struct {
	Token    string    `json:"token"`
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	user, err := s.repo.GetByPhone(ctx, input.Phone)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, err
	}
	token, err := s.createToken(user.ID, user.TenantID, user.Role)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		Token:    token,
		UserID:   user.ID,
		TenantID: user.TenantID,
		Name:     user.Name,
		Role:     user.Role,
	}, nil
}

func (s *Service) createToken(userID, tenantID uuid.UUID, role string) (string, error) {
	claims := &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   userID.String(),
		TenantID: tenantID.String(),
		Role:     role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
