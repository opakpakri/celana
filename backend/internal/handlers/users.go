package handlers

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"celana-api/internal/middleware"
	"celana-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *pgxpool.Pool
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var user models.User
	err := h.DB.QueryRow(context.Background(), `
		SELECT id, name, email, role, business_name, COALESCE(phone, ''), COALESCE(address, ''), created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil profil"})
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	if req.Name == "" || req.Email == "" || req.BusinessName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama, email, dan nama bisnis wajib diisi"})
	}

	ctx := context.Background()
	var user models.User
	var err error

	if strings.TrimSpace(req.Password) != "" {
		if len(req.Password) < 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
		}
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
		}
		err = h.DB.QueryRow(ctx, `
			UPDATE users
			SET name = $1, email = $2, business_name = $3, phone = $4, address = $5, password_hash = $6, updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
			RETURNING id, name, email, role, business_name, COALESCE(phone, ''), COALESCE(address, ''), created_at, updated_at
		`, req.Name, req.Email, req.BusinessName, req.Phone, req.Address, string(hash), userID).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = h.DB.QueryRow(ctx, `
			UPDATE users
			SET name = $1, email = $2, business_name = $3, phone = $4, address = $5, updated_at = CURRENT_TIMESTAMP
			WHERE id = $6
			RETURNING id, name, email, role, business_name, COALESCE(phone, ''), COALESCE(address, ''), created_at, updated_at
		`, req.Name, req.Email, req.BusinessName, req.Phone, req.Address, userID).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email sudah digunakan akun lain"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui profil"})
	}

	return c.JSON(user)
}

func normalizeRole(r string) string {
	r = strings.ToUpper(strings.TrimSpace(r))
	if r == "ADMIN" {
		return "ADMIN"
	}
	return "USER"
}

func (h *UserHandler) AdminListUsers(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	query := `
		SELECT id, name, email, role, business_name, COALESCE(phone, ''), COALESCE(address, ''), created_at, updated_at
		FROM users
	`
	args := make([]interface{}, 0)
	if search != "" {
		query += ` WHERE name ILIKE $1 OR email ILIKE $1 OR business_name ILIKE $1`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY created_at DESC, id DESC`

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar pengguna"})
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.BusinessName, &u.Phone, &u.Address, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data pengguna"})
		}
		users = append(users, u)
	}

	return c.JSON(users)
}

func (h *UserHandler) AdminCreateUser(c *fiber.Ctx) error {
	var req models.AdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.Role = normalizeRole(req.Role)
	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama, email, dan password (minimal 6 karakter) wajib diisi",
		})
	}
	if req.BusinessName == "" {
		req.BusinessName = "Celana - Cellix Laundry"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
	}

	var user models.User
	err = h.DB.QueryRow(context.Background(), `
		INSERT INTO users (name, email, password_hash, role, business_name, phone, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, email, role, business_name, COALESCE(phone, ''), COALESCE(address, ''), created_at, updated_at
	`, req.Name, req.Email, string(hash), req.Role, req.BusinessName, req.Phone, req.Address).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email sudah terdaftar"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat pengguna baru"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) AdminUpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}

	var req models.AdminUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.Role = normalizeRole(req.Role)
	if req.Name == "" || req.Email == "" || req.BusinessName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama, email, dan nama bisnis wajib diisi"})
	}

	ctx := context.Background()
	var user models.User
	if strings.TrimSpace(req.Password) != "" {
		if len(req.Password) < 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
		}
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses password"})
		}
		err = h.DB.QueryRow(ctx, `
			UPDATE users
			SET name = $1, email = $2, role = $3, business_name = $4, phone = $5, address = $6, password_hash = $7, updated_at = CURRENT_TIMESTAMP
			WHERE id = $8
			RETURNING id, name, email, role, business_name, COALESCE(phone, ''), COALESCE(address, ''), created_at, updated_at
		`, req.Name, req.Email, req.Role, req.BusinessName, req.Phone, req.Address, string(hash), id).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = h.DB.QueryRow(ctx, `
			UPDATE users
			SET name = $1, email = $2, role = $3, business_name = $4, phone = $5, address = $6, updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
			RETURNING id, name, email, role, business_name, COALESCE(phone, ''), COALESCE(address, ''), created_at, updated_at
		`, req.Name, req.Email, req.Role, req.BusinessName, req.Phone, req.Address, id).Scan(
			&user.ID, &user.Name, &user.Email, &user.Role, &user.BusinessName, &user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email sudah digunakan akun lain"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pengguna"})
	}

	return c.JSON(user)
}

func (h *UserHandler) AdminDeleteUser(c *fiber.Ctx) error {
	currentUserID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}
	if id == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak dapat menghapus akun Anda sendiri yang sedang aktif"})
	}

	tag, err := h.DB.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus pengguna"})
	}
	if tag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Pengguna berhasil dihapus"})
}
