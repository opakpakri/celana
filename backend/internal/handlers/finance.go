package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"celana-api/internal/middleware"
	"celana-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FinanceHandler struct {
	DB *pgxpool.Pool
}

func (h *FinanceHandler) Summary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var summary models.FinanceSummary

	err := h.DB.QueryRow(context.Background(), `
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE type = 'IN'), 0),
			COALESCE(SUM(amount) FILTER (WHERE type = 'OUT'), 0)
		FROM cash_flows
		WHERE user_id = $1
	`, userID).Scan(&summary.TotalIn, &summary.TotalOut)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil ringkasan keuangan"})
	}

	summary.NetBalance = summary.TotalIn - summary.TotalOut
	return c.JSON(summary)
}

func (h *FinanceHandler) CashFlows(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	flowType := strings.ToUpper(strings.TrimSpace(c.Query("type")))
	search := strings.TrimSpace(c.Query("search"))
	startDate := strings.TrimSpace(c.Query("start_date"))
	endDate := strings.TrimSpace(c.Query("end_date"))

	query := `
		SELECT cf.id, cf.user_id, cf.order_id, COALESCE(o.order_code, ''), cf.type,
		       cf.category, cf.payment_method, cf.amount, COALESCE(cf.description, ''), cf.transaction_date
		FROM cash_flows cf
		LEFT JOIN orders o ON o.id = cf.order_id AND o.user_id = cf.user_id
		WHERE cf.user_id = $1
	`
	args := []interface{}{userID}
	argIdx := 2

	if flowType == "IN" || flowType == "OUT" {
		query += fmt.Sprintf(" AND cf.type = $%d", argIdx)
		args = append(args, flowType)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (cf.category ILIKE $%d OR COALESCE(cf.description, '') ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if startDate != "" {
		query += fmt.Sprintf(" AND DATE(cf.transaction_date AT TIME ZONE 'Asia/Jakarta') >= $%d", argIdx)
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += fmt.Sprintf(" AND DATE(cf.transaction_date AT TIME ZONE 'Asia/Jakarta') <= $%d", argIdx)
		args = append(args, endDate)
		argIdx++
	}

	query += ` ORDER BY cf.transaction_date DESC, cf.id DESC`

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil riwayat arus kas"})
	}
	defer rows.Close()

	flows := make([]models.CashFlow, 0)
	for rows.Next() {
		var cf models.CashFlow
		if err := rows.Scan(
			&cf.ID, &cf.UserID, &cf.OrderID, &cf.OrderCode, &cf.Type,
			&cf.Category, &cf.PaymentMethod, &cf.Amount, &cf.Description, &cf.TransactionDate,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca arus kas"})
		}
		flows = append(flows, cf)
	}

	return c.JSON(flows)
}

func (h *FinanceHandler) CreateExpense(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req models.ExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Category = strings.ToUpper(strings.TrimSpace(req.Category))
	if req.Category == "" {
		req.Category = "OPERATIONAL"
	}
	req.PaymentMethod = normalizePaymentMethod(req.PaymentMethod)
	req.Description = strings.TrimSpace(req.Description)
	if req.Amount <= 0 || req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nominal (> 0) dan keterangan pengeluaran wajib diisi"})
	}

	var cf models.CashFlow
	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO cash_flows (user_id, order_id, type, category, payment_method, amount, description)
		VALUES ($1, NULL, 'OUT', $2, $3, $4, $5)
		RETURNING id, user_id, order_id, type, category, payment_method, amount, COALESCE(description, ''), transaction_date
	`, userID, req.Category, req.PaymentMethod, req.Amount, req.Description).Scan(
		&cf.ID, &cf.UserID, &cf.OrderID, &cf.Type, &cf.Category, &cf.PaymentMethod, &cf.Amount, &cf.Description, &cf.TransactionDate,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mencatat pengeluaran"})
	}

	return c.Status(fiber.StatusCreated).JSON(cf)
}

func (h *FinanceHandler) CreateTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req models.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	req.Type = strings.ToUpper(strings.TrimSpace(req.Type))
	if req.Type != "IN" && req.Type != "OUT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tipe transaksi harus IN atau OUT"})
	}

	req.Category = strings.TrimSpace(req.Category)
	if req.Category == "" {
		if req.Type == "IN" {
			req.Category = "Pemasukan Lainnya"
		} else {
			req.Category = "Operasional"
		}
	}
	req.PaymentMethod = normalizePaymentMethod(req.PaymentMethod)
	req.Description = strings.TrimSpace(req.Description)
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nominal harus lebih dari 0"})
	}

	txDate := time.Now()
	if strings.TrimSpace(req.TransactionDate) != "" {
		parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(req.TransactionDate), time.Local)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal tidak valid (YYYY-MM-DD)"})
		}
		now := time.Now()
		txDate = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
	}

	var cf models.CashFlow
	err := h.DB.QueryRow(context.Background(), `
		INSERT INTO cash_flows (user_id, order_id, type, category, payment_method, amount, description, transaction_date)
		VALUES ($1, NULL, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, order_id, type, category, payment_method, amount, COALESCE(description, ''), transaction_date
	`, userID, req.Type, req.Category, req.PaymentMethod, req.Amount, req.Description, txDate).Scan(
		&cf.ID, &cf.UserID, &cf.OrderID, &cf.Type, &cf.Category, &cf.PaymentMethod, &cf.Amount, &cf.Description, &cf.TransactionDate,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan transaksi keuangan"})
	}

	return c.Status(fiber.StatusCreated).JSON(cf)
}

func (h *FinanceHandler) DeleteExpense(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID transaksi tidak valid"})
	}

	tag, err := h.DB.Exec(context.Background(), `
		DELETE FROM cash_flows
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus catatan keuangan"})
	}

	if tag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Catatan keuangan tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"message": "Catatan keuangan berhasil dihapus"})
}
