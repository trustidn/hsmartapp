package adminauth

import (
	"context"
	"errors"
	"strings"
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

type Profile struct {
	AdminID string `json:"admin_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
}

func (s *Service) GetProfile(ctx context.Context, adminID uuid.UUID) (*Profile, error) {
	admin, err := s.repo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return &Profile{
		AdminID: admin.ID.String(),
		Name:    admin.Name,
		Email:   admin.Email,
		Role:    RoleSuperadmin,
	}, nil
}

type UpdateProfileInput struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	CurrentPassword *string `json:"current_password"`
	NewPassword     *string `json:"new_password"`
}

var (
	ErrEmailTaken      = errors.New("email already in use")
	ErrInvalidPassword = errors.New("current password invalid")
)

func (s *Service) UpdateProfile(ctx context.Context, adminID uuid.UUID, input UpdateProfileInput) error {
	admin, err := s.repo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}
	var name, email, passwordHash *string
	if input.Name != nil {
		n := strings.TrimSpace(*input.Name)
		if n != "" {
			name = &n
		}
	}
	if input.Email != nil {
		e := strings.TrimSpace(strings.ToLower(*input.Email))
		if e != "" {
			existing, _ := s.repo.GetByEmail(ctx, e)
			if existing != nil && existing.ID != adminID {
				return ErrEmailTaken
			}
			email = &e
		}
	}
	if input.NewPassword != nil {
		p := strings.TrimSpace(*input.NewPassword)
		if len(p) < 6 {
			return errors.New("password minimal 6 karakter")
		}
		if input.CurrentPassword == nil || *input.CurrentPassword == "" {
			return ErrInvalidPassword
		}
		if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(*input.CurrentPassword)); err != nil {
			return ErrInvalidPassword
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		h := string(hash)
		passwordHash = &h
	}
	return s.repo.UpdateProfile(ctx, adminID, name, email, passwordHash)
}
