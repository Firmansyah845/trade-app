package order

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/shopspring/decimal"
)

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo Repository
}

func (s *service) ProcessOrder(ctx context.Context, req CreateOrderRequest) (*ProcessOrderResponse, error) {
	log.Info().
		Int64("user_id", req.UserID).
		Str("stock_code", req.StockCode).
		Int64("quantity", req.Quantity).
		Str("price", req.Price.String()).
		Str("order_type", req.OrderType).
		Msg("processing order request")

	// Validasi input
	if err := s.validateOrderRequest(req); err != nil {
		log.Warn().
			Err(err).
			Int64("user_id", req.UserID).
			Str("stock_code", req.StockCode).
			Msg("order validation failed")
		return &ProcessOrderResponse{
			Status:  "failed",
			Message: err.Error(),
		}, nil
	}

	// 1. Cek stock dan validasi sharia
	stock, err := s.repo.GetStockInfo(ctx, req.StockCode)
	if err != nil {
		log.Error().
			Err(err).
			Str("stock_code", req.StockCode).
			Int64("user_id", req.UserID).
			Msg("failed to get stock info")
		return &ProcessOrderResponse{
			Status:  "failed",
			Message: "stock not found or inactive",
		}, nil
	}

	// Validasi saham syariah
	if !stock.IsSharia {
		log.Warn().
			Str("stock_code", req.StockCode).
			Int64("user_id", req.UserID).
			Bool("is_sharia", stock.IsSharia).
			Msg("non-sharia stock order rejected")
		return &ProcessOrderResponse{
			Status:   "failed",
			Message:  "non-sharia stock is not allowed",
			IsSharia: false,
		}, nil
	}

	// Hitung total amount
	totalAmount := req.Price.Mul(decimal.NewFromInt(req.Quantity))

	log.Info().
		Str("total_amount", totalAmount.String()).
		Int64("user_id", req.UserID).
		Str("stock_code", req.StockCode).
		Bool("is_sharia", stock.IsSharia).
		Msg("calculated total order amount")

	// Begin transaction untuk atomicity
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", req.UserID).
			Msg("failed to begin transaction for order processing")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Debug().
				Err(err).
				Int64("user_id", req.UserID).
				Msg("transaction rollback (expected if already committed)")
		}
	}()

	// 2. Cek balance user (hanya untuk order type buy)
	if req.OrderType == "buy" {
		balance, err := s.repo.GetUserBalance(ctx, req.UserID)
		if err != nil {
			log.Error().
				Err(err).
				Int64("user_id", req.UserID).
				Msg("failed to get user balance")
			return &ProcessOrderResponse{
				Status:   "failed",
				Message:  "user not found",
				IsSharia: stock.IsSharia,
			}, nil
		}

		log.Info().
			Int64("user_id", req.UserID).
			Str("current_balance", balance.Balance.String()).
			Str("required_amount", totalAmount.String()).
			Msg("checking user balance")

		if balance.Balance.LessThan(totalAmount) {
			log.Warn().
				Int64("user_id", req.UserID).
				Str("current_balance", balance.Balance.String()).
				Str("required_amount", totalAmount.String()).
				Msg("insufficient balance for order")
			return &ProcessOrderResponse{
				Status:   "failed",
				Message:  "insufficient balance",
				IsSharia: stock.IsSharia,
			}, nil
		}

		// Deduct balance
		err = s.repo.DeductBalance(ctx, tx, req.UserID, totalAmount)
		if err != nil {
			log.Error().
				Err(err).
				Int64("user_id", req.UserID).
				Str("amount", totalAmount.String()).
				Msg("failed to deduct balance")
			return &ProcessOrderResponse{
				Status:   "failed",
				Message:  "failed to process balance deduction",
				IsSharia: stock.IsSharia,
			}, nil
		}
	}

	// 3. Create order
	orderID, err := s.repo.CreateOrder(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", req.UserID).
			Str("stock_code", req.StockCode).
			Msg("failed to create order in database")
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Str("order_id", orderID).
			Int64("user_id", req.UserID).
			Msg("failed to commit order transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info().
		Str("order_id", orderID).
		Int64("user_id", req.UserID).
		Str("stock_code", req.StockCode).
		Str("total_amount", totalAmount.String()).
		Str("order_type", req.OrderType).
		Bool("is_sharia", stock.IsSharia).
		Msg("order processed successfully")

	return &ProcessOrderResponse{
		OrderID:     orderID,
		Status:      "success",
		Message:     "order created successfully",
		TotalAmount: totalAmount,
		IsSharia:    stock.IsSharia,
	}, nil
}

func (s *service) validateOrderRequest(req CreateOrderRequest) error {
	if req.UserID <= 0 {
		return fmt.Errorf("invalid user_id")
	}
	if req.StockCode == "" {
		return fmt.Errorf("stock_code is required")
	}
	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	if req.Price.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("price must be greater than 0")
	}
	if req.OrderType != "buy" && req.OrderType != "sell" {
		return fmt.Errorf("invalid order_type, must be 'buy' or 'sell'")
	}
	return nil
}

type Service interface {
	ProcessOrder(ctx context.Context, req CreateOrderRequest) (*ProcessOrderResponse, error)
}
