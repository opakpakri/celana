package handlers

import (
	"context"

	"celana-api/internal/middleware"
	"celana-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardHandler struct {
	DB *pgxpool.Pool
}

func (h *DashboardHandler) Summary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	ctx := context.Background()

	var summary models.DashboardSummary

	err := h.DB.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE DATE(created_at AT TIME ZONE 'Asia/Jakarta') = DATE(NOW() AT TIME ZONE 'Asia/Jakarta') AND status != 'CANCELLED'),
			COUNT(*) FILTER (WHERE status = 'PROCESSING'),
			COUNT(*) FILTER (WHERE status = 'READY')
		FROM orders
		WHERE user_id = $1
	`, userID).Scan(&summary.OrdersToday, &summary.ProcessingCount, &summary.ReadyCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil statistik pesanan"})
	}

	err = h.DB.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM cash_flows
		WHERE user_id = $1
		  AND type = 'IN'
		  AND DATE(transaction_date AT TIME ZONE 'Asia/Jakarta') = DATE(NOW() AT TIME ZONE 'Asia/Jakarta')
	`, userID).Scan(&summary.RevenueToday)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil pendapatan hari ini"})
	}

	rows, err := h.DB.Query(ctx, `
		SELECT o.id, o.order_code, c.name, o.status, o.payment_status, o.total_amount, o.updated_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id AND c.user_id = o.user_id
		WHERE o.user_id = $1
		ORDER BY o.updated_at DESC
		LIMIT 5
	`, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil aktivitas terbaru"})
	}
	defer rows.Close()

	summary.RecentActivity = make([]models.DashboardActivity, 0)
	for rows.Next() {
		var act models.DashboardActivity
		if err := rows.Scan(&act.ID, &act.OrderCode, &act.CustomerName, &act.Status, &act.PaymentStatus, &act.TotalAmount, &act.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca aktivitas"})
		}
		summary.RecentActivity = append(summary.RecentActivity, act)
	}

	return c.JSON(summary)
}
