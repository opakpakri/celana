package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	CORSOrigin  string
}

func Load() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/celana_db?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "celana-dev-secret-key-change-me"
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "*"
	}

	return Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		CORSOrigin:  corsOrigin,
	}
}

func ConnectDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	candidates := []string{
		filepath.Join("migrations", "001_init_multi_tenant.sql"),
		filepath.Join("backend", "migrations", "001_init_multi_tenant.sql"),
	}

	var sqlBytes []byte
	var err error
	for _, p := range candidates {
		sqlBytes, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	return nil
}

func EnsureAdminUser(ctx context.Context, pool *pgxpool.Pool) error {
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE role = 'ADMIN'`).Scan(&count); err != nil {
		return fmt.Errorf("failed to check admin user: %w", err)
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash default admin password: %w", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (name, email, password_hash, business_name, phone, address, role)
		VALUES ($1, $2, $3, $4, $5, $6, 'ADMIN')
		ON CONFLICT (email) DO UPDATE SET role = 'ADMIN'
	`, "Administrator", "admin@celana.id", string(hash), "Celana Pusat (Admin)", "081200000000", "Kantor Pusat")
	if err != nil {
		return fmt.Errorf("failed to seed default admin: %w", err)
	}
	return nil
}
