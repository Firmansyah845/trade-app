package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func NewDBStore(db *sql.DB) Repository {
	return &repo{db: db}
}

type repo struct {
	db *sql.DB
}

type Repository interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (string, error)
	GetOrder(ctx context.Context, req GetOrderRequest) (*[]Order, error)
	GetStockInfo(ctx context.Context, stockCode string) (*Stock, error)
	GetUserBalance(ctx context.Context, userID int64) (*UserBalance, error)
	DeductBalance(ctx context.Context, tx *sql.Tx, userID int64, amount decimal.Decimal) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

// CreateOrder membuat order baru
func (r *repo) CreateOrder(ctx context.Context, req CreateOrderRequest) (string, error) {
	id := uuid.New()
	totalAmount := req.Price.Mul(decimal.NewFromInt(req.Quantity))

	log.Info().
		Str("order_id", id.String()).
		Int64("user_id", req.UserID).
		Str("stock_code", req.StockCode).
		Int64("quantity", req.Quantity).
		Str("order_type", req.OrderType).
		Msg("creating new order")

	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO orders (id, user_id, stock_code, quantity, price, total_amount, order_type, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`)
	if err != nil {
		log.Error().
			Err(err).
			Str("order_id", id.String()).
			Int64("user_id", req.UserID).
			Msg("failed to prepare insert order statement")
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	_, err = stmt.ExecContext(ctx,
		id.String(),
		req.UserID,
		req.StockCode,
		req.Quantity,
		req.Price,
		totalAmount,
		req.OrderType,
		"pending",
		now,
		now,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("order_id", id.String()).
			Int64("user_id", req.UserID).
			Str("stock_code", req.StockCode).
			Msg("failed to execute insert order")
		return "", fmt.Errorf("failed to insert order: %w", err)
	}

	log.Info().
		Str("order_id", id.String()).
		Int64("user_id", req.UserID).
		Str("stock_code", req.StockCode).
		Str("total_amount", totalAmount.String()).
		Msg("order created successfully")

	return id.String(), nil
}

// GetOrder mengambil daftar order dengan filter
func (r *repo) GetOrder(ctx context.Context, req GetOrderRequest) (*[]Order, error) {
	log.Info().
		Int64("user_id", req.UserID).
		Str("stock_code", req.StockCode).
		Str("order_type", req.OrderType).
		Str("status", req.Status).
		Int("limit", req.Limit).
		Int("offset", req.Offset).
		Msg("fetching orders")

	query := `SELECT id, user_id, stock_code, quantity, price, total_amount, order_type, status, created_at, updated_at 
	          FROM orders 
	          WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	// Filter by user_id
	if req.UserID > 0 {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, req.UserID)
		argIndex++
	}

	// Filter by stock_code
	if req.StockCode != "" {
		query += fmt.Sprintf(" AND stock_code = $%d", argIndex)
		args = append(args, req.StockCode)
		argIndex++
	}

	// Filter by order_type
	if req.OrderType != "" {
		query += fmt.Sprintf(" AND order_type = $%d", argIndex)
		args = append(args, req.OrderType)
		argIndex++
	}

	// Filter by status
	if req.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, req.Status)
		argIndex++
	}

	// Filter by date range
	if req.DateFrom != "" {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, req.DateFrom)
		argIndex++
	}

	if req.DateTo != "" {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, req.DateTo)
		argIndex++
	}

	// Order by created_at DESC and add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", req.UserID).
			Str("query", query).
			Msg("failed to execute query orders")
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.StockCode,
			&order.Quantity,
			&order.Price,
			&order.TotalAmount,
			&order.OrderType,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			log.Error().
				Err(err).
				Int64("user_id", req.UserID).
				Msg("error scanning order row")
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		log.Error().
			Err(err).
			Int64("user_id", req.UserID).
			Msg("error iterating order rows")
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Info().
		Int("total_orders", len(orders)).
		Int64("user_id", req.UserID).
		Msg("orders fetched successfully")

	return &orders, nil
}

// GetStockInfo mengambil informasi saham
func (r *repo) GetStockInfo(ctx context.Context, stockCode string) (*Stock, error) {
	log.Info().
		Str("stock_code", stockCode).
		Msg("fetching stock info")

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT code, name, is_sharia, is_active, price, created_at, updated_at
		FROM stocks 
		WHERE code = $1 AND is_active = true
		LIMIT 1
	`)
	if err != nil {
		log.Error().
			Err(err).
			Str("stock_code", stockCode).
			Msg("failed to prepare select stock statement")
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var stock Stock
	err = stmt.QueryRowContext(ctx, stockCode).Scan(
		&stock.Code,
		&stock.Name,
		&stock.IsSharia,
		&stock.IsActive,
		&stock.Price,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().
				Str("stock_code", stockCode).
				Msg("stock not found")
			return nil, fmt.Errorf("stock not found")
		}
		log.Error().
			Err(err).
			Str("stock_code", stockCode).
			Msg("failed to fetch stock info")
		return nil, fmt.Errorf("failed to fetch stock: %w", err)
	}

	log.Info().
		Str("stock_code", stock.Code).
		Bool("is_sharia", stock.IsSharia).
		Str("price", stock.Price.String()).
		Msg("stock info fetched successfully")

	return &stock, nil
}

// GetUserBalance mengambil balance user dengan row lock
func (r *repo) GetUserBalance(ctx context.Context, userID int64) (*UserBalance, error) {
	log.Info().
		Int64("user_id", userID).
		Msg("fetching user balance with lock")

	stmt, err := r.db.PrepareContext(ctx, `
		SELECT user_id, balance, created_at, updated_at
		FROM user_balances 
		WHERE user_id = $1
		FOR UPDATE
	`)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", userID).
			Msg("failed to prepare select user balance statement")
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var balance UserBalance
	err = stmt.QueryRowContext(ctx, userID).Scan(
		&balance.UserID,
		&balance.Balance,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().
				Int64("user_id", userID).
				Msg("user balance not found")
			return nil, fmt.Errorf("user balance not found")
		}
		log.Error().
			Err(err).
			Int64("user_id", userID).
			Msg("failed to fetch user balance")
		return nil, fmt.Errorf("failed to fetch user balance: %w", err)
	}

	log.Info().
		Int64("user_id", balance.UserID).
		Str("balance", balance.Balance.String()).
		Msg("user balance fetched successfully")

	return &balance, nil
}

// DeductBalance mengurangi balance user dalam transaction
func (r *repo) DeductBalance(ctx context.Context, tx *sql.Tx, userID int64, amount decimal.Decimal) error {
	log.Info().
		Int64("user_id", userID).
		Str("amount", amount.String()).
		Msg("deducting user balance")

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE user_balances 
		SET balance = balance - $1, updated_at = $2
		WHERE user_id = $3 AND balance >= $1
	`)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", userID).
			Str("amount", amount.String()).
			Msg("failed to prepare update balance statement")
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, amount, time.Now(), userID)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", userID).
			Str("amount", amount.String()).
			Msg("failed to execute deduct balance")
		return fmt.Errorf("failed to deduct balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", userID).
			Msg("failed to get rows affected from deduct balance")
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Warn().
			Int64("user_id", userID).
			Str("amount", amount.String()).
			Msg("insufficient balance for deduction")
		return fmt.Errorf("insufficient balance")
	}

	log.Info().
		Int64("user_id", userID).
		Str("amount", amount.String()).
		Msg("balance deducted successfully")

	return nil
}

// BeginTx memulai transaction
func (r *repo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	log.Debug().Msg("beginning database transaction")

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	log.Debug().Msg("transaction started successfully")
	return tx, nil
}
