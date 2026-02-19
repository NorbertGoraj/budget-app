package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func errInternal(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func errBadRequest(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

// parseIDParam reads the :id URL parameter and writes a 400 on failure.
// Callers must return immediately when ok is false.
func parseIDParam(ctx *gin.Context) (int, bool) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		errBadRequest(ctx, "invalid id")
		return 0, false
	}
	return id, true
}
