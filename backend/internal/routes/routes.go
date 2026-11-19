package routes

import (
	"celana-api/internal/handlers"
	"celana-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(app *fiber.App, db *pgxpool.Pool, jwtSecret string) {
	authH := &handlers.AuthHandler{DB: db, JWTSecret: jwtSecret}
	dashH := &handlers.DashboardHandler{DB: db}
	orderH := &handlers.OrderHandler{DB: db}
	serviceH := &handlers.ServiceHandler{DB: db}
	customerH := &handlers.CustomerHandler{DB: db}
	financeH := &handlers.FinanceHandler{DB: db}
	userH := &handlers.UserHandler{DB: db}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "app": "Celana API"})
	})

	api := app.Group("/api")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authH.Register)
	auth.Post("/login", authH.Login)

	// Protected tenant routes
	protected := api.Group("", middleware.AuthTenant(jwtSecret))

	// Dashboard
	protected.Get("/dashboard/summary", dashH.Summary)

	// Orders
	protected.Get("/orders", orderH.List)
	protected.Post("/orders", orderH.Create)
	protected.Get("/orders/:id", orderH.Detail)
	protected.Patch("/orders/:id/status", orderH.UpdateStatus)
	protected.Post("/orders/:id/pay", orderH.Pay)

	// Services
	protected.Get("/services", serviceH.List)
	protected.Post("/services", serviceH.Create)
	protected.Put("/services/:id", serviceH.Update)
	protected.Patch("/services/:id/toggle", serviceH.Toggle)
	protected.Delete("/services/:id", serviceH.Delete)

	// Customers
	protected.Get("/customers", customerH.List)
	protected.Get("/customers/:id", customerH.Detail)
	protected.Post("/customers", customerH.Create)
	protected.Put("/customers/:id", customerH.Update)

	// Finance
	protected.Get("/finance/summary", financeH.Summary)
	protected.Get("/finance/cashflows", financeH.CashFlows)
	protected.Post("/finance/expenses", financeH.CreateExpense)
	protected.Post("/finance/transactions", financeH.CreateTransaction)
	protected.Delete("/finance/expenses/:id", financeH.DeleteExpense)

	// Users / Settings
	protected.Get("/users/profile", userH.GetProfile)
	protected.Put("/users/profile", userH.UpdateProfile)

	// Admin-only User Management
	admin := protected.Group("/admin", middleware.RequireAdmin(db))
	admin.Get("/users", userH.AdminListUsers)
	admin.Post("/users", userH.AdminCreateUser)
	admin.Put("/users/:id", userH.AdminUpdateUser)
	admin.Delete("/users/:id", userH.AdminDeleteUser)
}
