package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthTenant validates the Bearer JWT and injects userID exclusively from token claims.
func AuthTenant(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token autentikasi tidak ditemukan",
			})
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token autentikasi tidak valid",
			})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Metode tanda tangan token tidak valid")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Sesi telah berakhir atau token tidak valid",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Klaim token tidak valid",
			})
		}

		rawUserID, exists := claims["user_id"]
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Identitas tenant tidak ditemukan pada token",
			})
		}

		var userID int
		switch v := rawUserID.(type) {
		case float64:
			userID = int(v)
		case int:
			userID = v
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Format identitas tenant tidak valid",
			})
		}

		if userID <= 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Identitas tenant tidak valid",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}

// GetUserID extracts the authenticated tenant's userID from Fiber locals.
func GetUserID(c *fiber.Ctx) int {
	id, _ := c.Locals("userID").(int)
	return id
}

// RequireAdmin verifies in the database that the authenticated user has role ADMIN.
func RequireAdmin(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		var role string
		err := db.QueryRow(context.Background(), `SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
		if err != nil || role != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya Admin yang dapat mengakses fitur ini.",
			})
		}
		return c.Next()
	}
}
