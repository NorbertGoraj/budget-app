package controller

import (
	"net/http"

	"budget-app/domain"
	"budget-app/interface/service"
	"budget-app/utils"

	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	svc *service.BudgetService
}

func NewBudget(svc *service.BudgetService) *BudgetController {
	return &BudgetController{svc: svc}
}

func (c *BudgetController) GetAll(ctx *gin.Context) {
	budgets, err := c.svc.GetAll(ctx.Request.Context())
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.EnsureNotNil(budgets))
}

func (c *BudgetController) Create(ctx *gin.Context) {
	var b domain.Budget
	if err := ctx.ShouldBindJSON(&b); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	if err := c.svc.Create(ctx.Request.Context(), &b); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, b)
}

func (c *BudgetController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	var b domain.Budget
	if err := ctx.ShouldBindJSON(&b); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	b.ID = id
	if err := c.svc.Update(ctx.Request.Context(), &b); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, b)
}

func (c *BudgetController) Delete(ctx *gin.Context) {
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
