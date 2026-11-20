package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	BusinessName string    `json:"business_name"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AdminUserRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	Role         string `json:"role"`
	BusinessName string `json:"business_name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
}

type RegisterRequest struct {
	Name                string `json:"name"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	BusinessName        string `json:"business_name"`
	Phone               string `json:"phone"`
	Address             string `json:"address"`
	SeedDefaultServices bool   `json:"seed_default_services"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	BusinessName string `json:"business_name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Password     string `json:"password,omitempty"`
}

type Customer struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	Notes      string    `json:"notes"`
	OrderCount int       `json:"order_count,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type CustomerRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Notes   string `json:"notes"`
}

type Service struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	UnitType  string    `json:"unit_type"`
	Price     float64   `json:"price"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ServiceRequest struct {
	Name     string  `json:"name"`
	UnitType string  `json:"unit_type"`
	Price    float64 `json:"price"`
}

type ServiceMetrics struct {
	Total   int `json:"total"`
	Active  int `json:"active"`
	PerKg   int `json:"per_kg"`
	PerItem int `json:"per_item"`
	PerM2   int `json:"per_m2"`
}

type Order struct {
	ID            int         `json:"id"`
	UserID        int         `json:"user_id"`
	OrderCode     string      `json:"order_code"`
	CustomerID    int         `json:"customer_id"`
	CustomerName  string      `json:"customer_name,omitempty"`
	CustomerPhone string      `json:"customer_phone,omitempty"`
	Status        string      `json:"status"`
	PaymentStatus string      `json:"payment_status"`
	TotalAmount   float64     `json:"total_amount"`
	PaidAmount    float64     `json:"paid_amount"`
	Notes           string      `json:"notes"`
	ItemsSummary    string      `json:"items_summary,omitempty"`
	ServicesSummary string      `json:"services_summary,omitempty"`
	QtySummary      string      `json:"qty_summary,omitempty"`
	Items           []OrderItem `json:"items,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID          int     `json:"id"`
	OrderID     int     `json:"order_id"`
	ServiceID   int     `json:"service_id"`
	ServiceName string  `json:"service_name,omitempty"`
	UnitType    string  `json:"unit_type,omitempty"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

type CreateOrderItemInput struct {
	ServiceID int     `json:"service_id"`
	Quantity  float64 `json:"quantity"`
}

type CreateOrderRequest struct {
	CustomerID    int                    `json:"customer_id"`
	Items         []CreateOrderItemInput `json:"items"`
	Notes         string                 `json:"notes"`
	PaidAmount    float64                `json:"paid_amount"`
	PaymentMethod string                 `json:"payment_method"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

type PayOrderRequest struct {
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

type CashFlow struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	OrderID         *int      `json:"order_id"`
	OrderCode       string    `json:"order_code,omitempty"`
	Type            string    `json:"type"`
	Category        string    `json:"category"`
	PaymentMethod   string    `json:"payment_method"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
}

type ExpenseRequest struct {
	Category      string  `json:"category"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
}

type TransactionRequest struct {
	Type            string  `json:"type"`
	Category        string  `json:"category"`
	PaymentMethod   string  `json:"payment_method"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
}

type FinanceSummary struct {
	TotalIn    float64 `json:"total_in"`
	TotalOut   float64 `json:"total_out"`
	NetBalance float64 `json:"net_balance"`
}

type DashboardActivity struct {
	ID            int       `json:"id"`
	OrderCode     string    `json:"order_code"`
	CustomerName  string    `json:"customer_name"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	TotalAmount   float64   `json:"total_amount"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DashboardSummary struct {
	OrdersToday     int                 `json:"orders_today"`
	RevenueToday    float64             `json:"revenue_today"`
	ProcessingCount int                 `json:"processing_count"`
	ReadyCount      int                 `json:"ready_count"`
	RecentActivity  []DashboardActivity `json:"recent_activity"`
}
