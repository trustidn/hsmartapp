package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	dsn := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/hsmart?sslmode=disable")
	email := getEnv("ADMIN_EMAIL", "admin@hsmart.app")
	password := getEnv("ADMIN_PASSWORD", "admin123")
	name := getEnv("ADMIN_NAME", "Super Admin")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash: %v", err)
	}

	id := uuid.New()
	_, err = pool.Exec(ctx, `
		INSERT INTO superadmins (id, email, password_hash, name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO NOTHING
	`, id, email, string(hash), name)
	if err != nil {
		log.Fatalf("insert: %v", err)
	}

	fmt.Printf("Superadmin seeded: %s (%s)\n", email, name)
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
