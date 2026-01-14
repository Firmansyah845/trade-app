package order

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	UserID    int64           `json:"user_id"`
	StockCode string          `json:"stock_code"`
	Quantity  int64           `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	OrderType string          `json:"order_type"`
}

type GetOrderRequest struct {
	UserID    int64  `json:"user_id,omitempty"`
	StockCode string `json:"stock_code,omitempty"`
	OrderType string `json:"order_type,omitempty"`
	Status    string `json:"status,omitempty"`
	DateFrom  string `json:"date_from,omitempty"`
	DateTo    string `json:"date_to,omitempty"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

type Order struct {
	ID          string          `json:"id"`
	UserID      int64           `json:"user_id"`
	StockCode   string          `json:"stock_code"`
	Quantity    int64           `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	OrderType   string          `json:"order_type"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Stock struct {
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	IsSharia  bool            `json:"is_sharia"`
	IsActive  bool            `json:"is_active"`
	Price     decimal.Decimal `json:"price"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type UserBalance struct {
	UserID    int64           `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ProcessOrderResponse struct {
	OrderID     string          `json:"order_id"`
	Status      string          `json:"status"`
	Message     string          `json:"message"`
	IsSharia    bool            `json:"is_sharia,omitempty"`
	TotalAmount decimal.Decimal `json:"total_amount,omitempty"`
}
