package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	status, err := h.healthCheckService.GetStatus(ctx)
	if err != nil {
		asInternalErrorResponse(c, err)
		return
	}

	asJsonResponse(c, http.StatusOK, "success", status)
}
