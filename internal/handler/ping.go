package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	asJsonResponse(c, http.StatusOK, "pong", nil)
}
