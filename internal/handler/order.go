package handler

import (
	"awesomeProjectCr/internal/order"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()

	log.Info().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("remote_addr", c.ClientIP()).
		Msg("received create order request")

	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().
			Err(err).
			Str("remote_addr", c.ClientIP()).
			Msg("invalid request body for create order")

		asErrorResponse(c, http.StatusBadRequest, err.Error())
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
		asErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	asJsonResponse(c, statusCode, "create order success", resp)
}
