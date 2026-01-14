package handler

import (
	"net/http"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status, err := h.healthCheckService.GetStatus(ctx)
	if err != nil {
		asInternalErrorResponse(w, err)
		return
	}

	asJsonResponse(w, http.StatusOK, "success", status)

}
