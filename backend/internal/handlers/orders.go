package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"celana-api/internal/middleware"
	"celana-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	DB *pgxpool.Pool
}

func isValidOrderStatus(s string) bool {
	switch s {
	case "PENDING", "PROCESSING", "READY", "COMPLETED", "CANCELLED":
		return true
	}
	return false
}

func normalizePaymentMethod(m string) string {
	m = strings.ToUpper(strings.TrimSpace(m))
	if m == "" {
		return "TUNAI"
	}
	return m
}

func (h *OrderHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	status := strings.ToUpper(strings.TrimSpace(c.Query("status")))
	search := strings.TrimSpace(c.Query("search"))
	ctx := context.Background()

	query := `
		SELECT o.id, o.user_id, o.order_code, o.customer_id, c.name, c.phone,
		       o.status, o.payment_status, o.total_amount, o.paid_amount, COALESCE(o.notes, ''),
		       COALESCE(STRING_AGG(s.name || ' (' || oi.quantity || ')', ', '), '') AS items_summary,
		       o.created_at, o.updated_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id AND c.user_id = o.user_id
		LEFT JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN services s ON s.id = oi.service_id
		WHERE o.user_id = $1
	`
	args := []interface{}{userID}
	argIdx := 2

	if status != "" && status != "ALL" {
		query += fmt.Sprintf(" AND o.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (o.order_code ILIKE $%d OR c.name ILIKE $%d OR c.phone ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += ` GROUP BY o.id, c.name, c.phone ORDER BY o.created_at DESC, o.id DESC`

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar pesanan"})
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.OrderCode, &o.CustomerID, &o.CustomerName, &o.CustomerPhone,
			&o.Status, &o.PaymentStatus, &o.TotalAmount, &o.PaidAmount, &o.Notes,
			&o.ItemsSummary, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data pesanan"})
		}
		orders = append(orders, o)
	}

	return c.JSON(orders)
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	if req.CustomerID <= 0 || len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Pelanggan dan minimal satu item layanan wajib dipilih"})
	}
	if req.PaidAmount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nominal pembayaran tidak boleh negatif"})
	}

	ctx := context.Background()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memulai transaksi"})
	}
	defer tx.Rollback(ctx)

	// Cross-tenant validation: customer must belong to authenticated userID
	var custName, custPhone string
	err = tx.QueryRow(ctx, `
		SELECT name, phone FROM customers WHERE id = $1 AND user_id = $2
	`, req.CustomerID, userID).Scan(&custName, &custPhone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pelanggan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memvalidasi pelanggan"})
	}

	// Cross-tenant validation & price lookup for all services
	type resolvedItem struct {
		ServiceID   int
		ServiceName string
		UnitType    string
		Quantity    float64
		UnitPrice   float64
		Subtotal    float64
	}
	resolvedItems := make([]resolvedItem, 0, len(req.Items))
	var totalAmount float64

	for _, item := range req.Items {
		if item.ServiceID <= 0 || item.Quantity <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Setiap item harus memiliki layanan valid dan jumlah > 0"})
		}
		var sName, sUnit string
		var sPrice float64
		var sActive bool
		err := tx.QueryRow(ctx, `
			SELECT name, unit_type, price, is_active
			FROM services
			WHERE id = $1 AND user_id = $2
		`, item.ServiceID, userID).Scan(&sName, &sUnit, &sPrice, &sActive)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Salah satu layanan tidak ditemukan"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memvalidasi layanan"})
		}
		if !sActive {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Layanan '%s' sedang nonaktif", sName)})
		}

		sub := sPrice * item.Quantity
		totalAmount += sub
		resolvedItems = append(resolvedItems, resolvedItem{
			ServiceID:   item.ServiceID,
			ServiceName: sName,
			UnitType:    sUnit,
			Quantity:    item.Quantity,
			UnitPrice:   sPrice,
			Subtotal:    sub,
		})
	}

	// Cap recorded paid_amount at totalAmount (excess is change given to customer)
	recordedPaid := req.PaidAmount
	if recordedPaid > totalAmount {
		recordedPaid = totalAmount
	}

	paymentStatus := "UNPAID"
	if recordedPaid >= totalAmount && totalAmount > 0 {
		paymentStatus = "PAID"
	} else if recordedPaid > 0 {
		paymentStatus = "PARTIAL"
	}

	// Generate unique order code per tenant: CLX-YYYYMMDD-XXX
	todayPrefix := "CLX-" + time.Now().Format("20060102") + "-"
	var seq int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) + 1 FROM orders
		WHERE user_id = $1 AND order_code LIKE $2
	`, userID, todayPrefix+"%").Scan(&seq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat nomor nota"})
	}
	orderCode := fmt.Sprintf("%s%03d", todayPrefix, seq)

	var order models.Order
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, order_code, customer_id, status, payment_status, total_amount, paid_amount, notes)
		VALUES ($1, $2, $3, 'PENDING', $4, $5, $6, $7)
		RETURNING id, user_id, order_code, customer_id, status, payment_status, total_amount, paid_amount, COALESCE(notes, ''), created_at, updated_at
	`, userID, orderCode, req.CustomerID, paymentStatus, totalAmount, recordedPaid, req.Notes).Scan(
		&order.ID, &order.UserID, &order.OrderCode, &order.CustomerID, &order.Status,
		&order.PaymentStatus, &order.TotalAmount, &order.PaidAmount, &order.Notes,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan pesanan"})
	}
	order.CustomerName = custName
	order.CustomerPhone = custPhone
	order.Items = make([]models.OrderItem, 0, len(resolvedItems))

	for _, it := range resolvedItems {
		var oi models.OrderItem
		err := tx.QueryRow(ctx, `
			INSERT INTO order_items (order_id, service_id, quantity, unit_price, subtotal)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, order_id, service_id, quantity, unit_price, subtotal
		`, order.ID, it.ServiceID, it.Quantity, it.UnitPrice, it.Subtotal).Scan(
			&oi.ID, &oi.OrderID, &oi.ServiceID, &oi.Quantity, &oi.UnitPrice, &oi.Subtotal,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan item pesanan"})
		}
		oi.ServiceName = it.ServiceName
		oi.UnitType = it.UnitType
		order.Items = append(order.Items, oi)
	}

	if recordedPaid > 0 {
		pm := normalizePaymentMethod(req.PaymentMethod)
		desc := fmt.Sprintf("Pembayaran pesanan %s (%s)", order.OrderCode, custName)
		if _, err := tx.Exec(ctx, `
			INSERT INTO cash_flows (user_id, order_id, type, category, payment_method, amount, description)
			VALUES ($1, $2, 'IN', 'ORDER_PAYMENT', $3, $4, $5)
		`, userID, order.ID, pm, recordedPaid, desc); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mencatat arus kas"})
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyelesaikan pesanan"})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *OrderHandler) Detail(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pesanan tidak valid"})
	}
	ctx := context.Background()

	var order models.Order
	err = h.DB.QueryRow(ctx, `
		SELECT o.id, o.user_id, o.order_code, o.customer_id, c.name, c.phone,
		       o.status, o.payment_status, o.total_amount, o.paid_amount, COALESCE(o.notes, ''),
		       o.created_at, o.updated_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id AND c.user_id = o.user_id
		WHERE o.id = $1 AND o.user_id = $2
	`, id, userID).Scan(
		&order.ID, &order.UserID, &order.OrderCode, &order.CustomerID, &order.CustomerName, &order.CustomerPhone,
		&order.Status, &order.PaymentStatus, &order.TotalAmount, &order.PaidAmount, &order.Notes,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pesanan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil detail pesanan"})
	}

	rows, err := h.DB.Query(ctx, `
		SELECT oi.id, oi.order_id, oi.service_id, s.name, s.unit_type, oi.quantity, oi.unit_price, oi.subtotal
		FROM order_items oi
		JOIN services s ON s.id = oi.service_id
		WHERE oi.order_id = $1
		ORDER BY oi.id ASC
	`, order.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil item pesanan"})
	}
	defer rows.Close()

	order.Items = make([]models.OrderItem, 0)
	for rows.Next() {
		var oi models.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ServiceID, &oi.ServiceName, &oi.UnitType, &oi.Quantity, &oi.UnitPrice, &oi.Subtotal); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca item pesanan"})
		}
		order.Items = append(order.Items, oi)
	}

	return c.JSON(order)
}

func isAllowedTransition(currentStatus, nextStatus string) bool {
	if currentStatus == "COMPLETED" || currentStatus == "CANCELLED" {
		return false
	}
	if nextStatus == "CANCELLED" {
		return true
	}
	switch currentStatus {
	case "PENDING":
		return nextStatus == "PROCESSING"
	case "PROCESSING":
		return nextStatus == "READY"
	case "READY":
		return nextStatus == "COMPLETED"
	default:
		return false
	}
}

func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pesanan tidak valid"})
	}

	var req models.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	if !isValidOrderStatus(req.Status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status pesanan tidak valid"})
	}

	ctx := context.Background()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memulai transaksi"})
	}
	defer tx.Rollback(ctx)

	var currentStatus, orderCode, custName string
	var paidAmount float64
	err = tx.QueryRow(ctx, `
		SELECT o.status, o.order_code, c.name, o.paid_amount
		FROM orders o
		JOIN customers c ON c.id = o.customer_id AND c.user_id = o.user_id
		WHERE o.id = $1 AND o.user_id = $2
		FOR UPDATE OF o
	`, id, userID).Scan(&currentStatus, &orderCode, &custName, &paidAmount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pesanan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pesanan"})
	}

	if !isAllowedTransition(currentStatus, req.Status) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status harus diubah secara berurutan (Antre -> Proses -> Siap Diambil -> Selesai)",
		})
	}

	var order models.Order
	err = tx.QueryRow(ctx, `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, order_code, customer_id, status, payment_status, total_amount, paid_amount, COALESCE(notes, ''), created_at, updated_at
	`, req.Status, id, userID).Scan(
		&order.ID, &order.UserID, &order.OrderCode, &order.CustomerID,
		&order.Status, &order.PaymentStatus, &order.TotalAmount, &order.PaidAmount, &order.Notes,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui status pesanan"})
	}

	if req.Status == "CANCELLED" && paidAmount > 0 {
		desc := fmt.Sprintf("Pembatalan pesanan %s (%s) — pengembalian dana", orderCode, custName)
		if _, err := tx.Exec(ctx, `
			INSERT INTO cash_flows (user_id, order_id, type, category, payment_method, amount, description)
			VALUES ($1, $2, 'OUT', 'ORDER_CANCELLED', 'TUNAI', $3, $4)
		`, userID, id, paidAmount, desc); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mencatat laporan pembatalan ke keuangan"})
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan perubahan status"})
	}

	return c.JSON(order)
}

func (h *OrderHandler) Pay(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pesanan tidak valid"})
	}

	var req models.PayOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nominal pembayaran harus lebih dari 0"})
	}

	ctx := context.Background()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memulai transaksi"})
	}
	defer tx.Rollback(ctx)

	var orderCode, custName string
	var totalAmount, paidAmount float64
	err = tx.QueryRow(ctx, `
		SELECT o.order_code, c.name, o.total_amount, o.paid_amount
		FROM orders o
		JOIN customers c ON c.id = o.customer_id AND c.user_id = o.user_id
		WHERE o.id = $1 AND o.user_id = $2
		FOR UPDATE OF o
	`, id, userID).Scan(&orderCode, &custName, &totalAmount, &paidAmount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pesanan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pesanan"})
	}

	remaining := totalAmount - paidAmount
	if remaining <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Pesanan sudah lunas"})
	}

	appliedAmount := req.Amount
	if appliedAmount > remaining {
		appliedAmount = remaining
	}

	newPaid := paidAmount + appliedAmount
	paymentStatus := "PARTIAL"
	if newPaid >= totalAmount {
		paymentStatus = "PAID"
	}

	var updated models.Order
	err = tx.QueryRow(ctx, `
		UPDATE orders
		SET paid_amount = $1, payment_status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, order_code, customer_id, status, payment_status, total_amount, paid_amount, COALESCE(notes, ''), created_at, updated_at
	`, newPaid, paymentStatus, id, userID).Scan(
		&updated.ID, &updated.UserID, &updated.OrderCode, &updated.CustomerID,
		&updated.Status, &updated.PaymentStatus, &updated.TotalAmount, &updated.PaidAmount,
		&updated.Notes, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pembayaran pesanan"})
	}

	pm := normalizePaymentMethod(req.PaymentMethod)
	desc := fmt.Sprintf("Pelunasan/Cicilan pesanan %s (%s)", orderCode, custName)
	if _, err := tx.Exec(ctx, `
		INSERT INTO cash_flows (user_id, order_id, type, category, payment_method, amount, description)
		VALUES ($1, $2, 'IN', 'ORDER_PAYMENT', $3, $4, $5)
	`, userID, id, pm, appliedAmount, desc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mencatat pembayaran ke arus kas"})
	}

	if err := tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan pembayaran"})
	}

	return c.JSON(fiber.Map{
		"order":          updated,
		"applied_amount": appliedAmount,
		"change_amount":  req.Amount - appliedAmount,
	})
}
