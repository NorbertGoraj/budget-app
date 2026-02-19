package controller

import (
	"net/http"

	"budget-app/domain"
	"budget-app/interface/service"
	"budget-app/utils"

	"github.com/gin-gonic/gin"
)

type InvestmentController struct {
	svc *service.InvestmentService
}

func NewInvestment(svc *service.InvestmentService) *InvestmentController {
	return &InvestmentController{svc: svc}
}

func (c *InvestmentController) GetAll(ctx *gin.Context) {
	investments, err := c.svc.GetAll(ctx.Request.Context())
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.EnsureNotNil(investments))
}

func (c *InvestmentController) Create(ctx *gin.Context) {
	var inv domain.Investment
	if err := ctx.ShouldBindJSON(&inv); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	if err := c.svc.Create(ctx.Request.Context(), &inv); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, inv)
}

func (c *InvestmentController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	var inv domain.Investment
	if err := ctx.ShouldBindJSON(&inv); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	inv.ID = id
	if err := c.svc.Update(ctx.Request.Context(), &inv); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, inv)
}

func (c *InvestmentController) Delete(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted": true})
}
