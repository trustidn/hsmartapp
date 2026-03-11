package adminauth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type adminJWTClaims struct {
	jwt.RegisteredClaims
	AdminID string `json:"aid"`
	Role    string `json:"role"`
}

const RoleSuperadmin = "superadmin"

type Service struct {
	repo   *Repository
	secret []byte
}

func NewService(repo *Repository, secret []byte) *Service {
	return &Service{repo: repo, secret: secret}
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	Token  string    `json:"token"`
	AdminID string    `json:"admin_id"`
	Name   string    `json:"name"`
	Role   string    `json:"role"`
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	admin, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		return nil, err
	}
	token, err := s.createToken(admin.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		Token:   token,
		AdminID: admin.ID.String(),
		Name:    admin.Name,
		Role:    RoleSuperadmin,
	}, nil
}

func (s *Service) createToken(adminID uuid.UUID) (string, error) {
	claims := &adminJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		AdminID: adminID.String(),
		Role:    RoleSuperadmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
