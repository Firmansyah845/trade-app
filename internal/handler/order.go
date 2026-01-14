package handler

import (
	"awesomeProjectCr/internal/order"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Msg("received create order request")

	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().
			Err(err).
			Str("remote_addr", r.RemoteAddr).
			Msg("invalid request body for create order")

		asErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info().
		Int64("user_id", req.UserID).
		Str("stock_code", req.StockCode).
		Str("order_type", req.OrderType).
		Msg("processing create order request")

	resp, err := h.orderService.ProcessOrder(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", req.UserID).
			Str("stock_code", req.StockCode).
			Msg("failed to process order in handler")
		asErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	statusCode := http.StatusOK
	if resp.Status == "failed" {
		statusCode = http.StatusBadRequest
		log.Warn().
			Int64("user_id", req.UserID).
			Str("stock_code", req.StockCode).
			Str("message", resp.Message).
			Msg("order creation failed")
	} else if resp.Status == "success" {
		statusCode = http.StatusCreated
		log.Info().
			Str("order_id", resp.OrderID).
			Int64("user_id", req.UserID).
			Str("stock_code", req.StockCode).
			Msg("order created successfully in handler")
	}

	asJsonResponse(w, statusCode, "create order success", resp)
}
