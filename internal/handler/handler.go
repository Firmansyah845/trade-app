package handler

import (
	"awesomeProjectCr/internal/health_checks"
	"awesomeProjectCr/internal/order"
	"database/sql"
	"encoding/json"
	"net/http"
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

func asJsonResponse(w http.ResponseWriter, httpStatus int, message string, data any) {
	response := Response{
		Message:    message,
		Data:       data,
		StatusCode: httpStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	_ = json.NewEncoder(w).Encode(response)
}

func asErrorResponse(w http.ResponseWriter, httpStatus int, message string) {
	response := ErrorResponse{
		Message: message,
		Code:    httpStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	_ = json.NewEncoder(w).Encode(response)
}

func asInternalErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}
