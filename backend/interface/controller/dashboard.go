package controller

import (
	"net/http"

	"budget-app/interface/service"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	svc *service.DashboardService
}

func NewDashboard(svc *service.DashboardService) *DashboardController {
	return &DashboardController{svc: svc}
}

func (c *DashboardController) Get(ctx *gin.Context) {
	resp, err := c.svc.GetDashboard(ctx.Request.Context())
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
