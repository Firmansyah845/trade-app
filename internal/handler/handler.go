package handler

import (
	"awesomeProjectCr/internal/health_checks"
	"awesomeProjectCr/internal/order"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	healthCheckService health_checks.Service
	orderService       order.Service
}

func NewHandler(db *sql.DB) *Handler {

	healthCheckService := health_checks.NewService(db)

	orderService := order.NewService(order.NewDBStore(db))
	return &Handler{
		healthCheckService: *healthCheckService,
		orderService:       orderService,
	}
}

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"contents"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func asJsonResponse(c *gin.Context, httpStatus int, message string, data any) {
	c.JSON(httpStatus, Response{
		Message:    message,
		Data:       data,
		StatusCode: httpStatus,
	})
}

func asErrorResponse(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, ErrorResponse{
		Message: message,
		Code:    httpStatus,
	})
}

func asInternalErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
	})
}
