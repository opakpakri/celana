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

type CustomerHandler struct {
	DB *pgxpool.Pool
}

func (h *CustomerHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	search := strings.TrimSpace(c.Query("search"))
	ctx := context.Background()

	query := `
		SELECT c.id, c.user_id, c.name, c.phone, COALESCE(c.address, ''), COALESCE(c.notes, ''),
		       COUNT(o.id) AS order_count, c.created_at
		FROM customers c
		LEFT JOIN orders o ON o.customer_id = c.id AND o.user_id = c.user_id
		WHERE c.user_id = $1
	`
	args := []interface{}{userID}
	if search != "" {
		query += ` AND (c.name ILIKE $2 OR c.phone ILIKE $2)`
		args = append(args, "%"+search+"%")
	}
	query += ` GROUP BY c.id ORDER BY c.created_at DESC, c.id DESC`

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar pelanggan"})
	}
	defer rows.Close()

	customers := make([]models.Customer, 0)
	for rows.Next() {
		var cust models.Customer
		if err := rows.Scan(&cust.ID, &cust.UserID, &cust.Name, &cust.Phone, &cust.Address, &cust.Notes, &cust.OrderCount, &cust.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data pelanggan"})
		}
		customers = append(customers, cust)
	}

	return c.JSON(customers)
}

func (h *CustomerHandler) Detail(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pelanggan tidak valid"})
	}
	ctx := context.Background()

	var cust models.Customer
	err = h.DB.QueryRow(ctx, `
		SELECT id, user_id, name, phone, COALESCE(address, ''), COALESCE(notes, ''), created_at
		FROM customers
		WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&cust.ID, &cust.UserID, &cust.Name, &cust.Phone, &cust.Address, &cust.Notes, &cust.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pelanggan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil detail pelanggan"})
	}

	rows, err := h.DB.Query(ctx, `
		SELECT id, order_code, status, payment_status, total_amount, paid_amount, created_at
		FROM orders
		WHERE customer_id = $1 AND user_id = $2
		ORDER BY created_at DESC
	`, id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil riwayat pesanan pelanggan"})
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.OrderCode, &o.Status, &o.PaymentStatus, &o.TotalAmount, &o.PaidAmount, &o.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca riwayat pesanan"})
		}
		orders = append(orders, o)
	}

	return c.JSON(fiber.Map{
		"customer": cust,
		"orders":   orders,
	})
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req models.CustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Name == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama dan nomor telepon wajib diisi"})
	}

	var cust models.Customer
	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO customers (user_id, name, phone, address, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, phone, COALESCE(address, ''), COALESCE(notes, ''), created_at
	`, userID, req.Name, req.Phone, req.Address, req.Notes).Scan(
		&cust.ID, &cust.UserID, &cust.Name, &cust.Phone, &cust.Address, &cust.Notes, &cust.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Nomor telepon pelanggan sudah terdaftar"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah pelanggan"})
	}

	return c.Status(fiber.StatusCreated).JSON(cust)
}

func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pelanggan tidak valid"})
	}

	var req models.CustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Name == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama dan nomor telepon wajib diisi"})
	}

	var cust models.Customer
	err = h.DB.QueryRow(context.Background(), `
		UPDATE customers
		SET name = $1, phone = $2, address = $3, notes = $4
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, name, phone, COALESCE(address, ''), COALESCE(notes, ''), created_at
	`, req.Name, req.Phone, req.Address, req.Notes, id, userID).Scan(
		&cust.ID, &cust.UserID, &cust.Name, &cust.Phone, &cust.Address, &cust.Notes, &cust.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pelanggan tidak ditemukan"})
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Nomor telepon sudah digunakan pelanggan lain"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pelanggan"})
	}

	return c.JSON(cust)
}
