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
)

type ServiceHandler struct {
	DB *pgxpool.Pool
}

func isValidUnitType(u string) bool {
	return u == "PER_KG" || u == "PER_ITEM" || u == "PER_M2"
}

func (h *ServiceHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	ctx := context.Background()

	rows, err := h.DB.Query(ctx, `
		SELECT id, user_id, name, unit_type, price, is_active, created_at
		FROM services
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
	`, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar layanan"})
	}
	defer rows.Close()

	services := make([]models.Service, 0)
	var metrics models.ServiceMetrics
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.UnitType, &s.Price, &s.IsActive, &s.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data layanan"})
		}
		services = append(services, s)
		metrics.Total++
		if s.IsActive {
			metrics.Active++
		}
		switch s.UnitType {
		case "PER_KG":
			metrics.PerKg++
		case "PER_ITEM":
			metrics.PerItem++
		case "PER_M2":
			metrics.PerM2++
		}
	}

	return c.JSON(fiber.Map{
		"services": services,
		"metrics":  metrics,
	})
}

func (h *ServiceHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req models.ServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.UnitType = strings.TrimSpace(req.UnitType)
	if req.Name == "" || !isValidUnitType(req.UnitType) || req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama layanan, tipe satuan (PER_KG/PER_ITEM/PER_M2), dan harga valid wajib diisi"})
	}

	var s models.Service
	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO services (user_id, name, unit_type, price, is_active)
		VALUES ($1, $2, $3, $4, TRUE)
		RETURNING id, user_id, name, unit_type, price, is_active, created_at
	`, userID, req.Name, req.UnitType, req.Price).Scan(
		&s.ID, &s.UserID, &s.Name, &s.UnitType, &s.Price, &s.IsActive, &s.CreatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah layanan"})
	}

	return c.Status(fiber.StatusCreated).JSON(s)
}

func (h *ServiceHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID layanan tidak valid"})
	}

	var req models.ServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.UnitType = strings.TrimSpace(req.UnitType)
	if req.Name == "" || !isValidUnitType(req.UnitType) || req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data layanan tidak valid"})
	}

	var s models.Service
	err = h.DB.QueryRow(context.Background(), `
		UPDATE services
		SET name = $1, unit_type = $2, price = $3
		WHERE id = $4 AND user_id = $5
		RETURNING id, user_id, name, unit_type, price, is_active, created_at
	`, req.Name, req.UnitType, req.Price, id, userID).Scan(
		&s.ID, &s.UserID, &s.Name, &s.UnitType, &s.Price, &s.IsActive, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Layanan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui layanan"})
	}

	return c.JSON(s)
}

func (h *ServiceHandler) Toggle(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID layanan tidak valid"})
	}

	var s models.Service
	err = h.DB.QueryRow(context.Background(), `
		UPDATE services
		SET is_active = NOT is_active
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, name, unit_type, price, is_active, created_at
	`, id, userID).Scan(
		&s.ID, &s.UserID, &s.Name, &s.UnitType, &s.Price, &s.IsActive, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Layanan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengubah status layanan"})
	}

	return c.JSON(s)
}

func (h *ServiceHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID layanan tidak valid"})
	}

	tag, err := h.DB.Exec(context.Background(), `
		DELETE FROM services
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Layanan tidak dapat dihapus karena sudah digunakan pada pesanan. Gunakan fitur Nonaktifkan.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus layanan"})
	}

	if tag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Layanan tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Layanan berhasil dihapus"})
}
