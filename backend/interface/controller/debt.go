package controller

import (
	"net/http"
	"strconv"

	"budget-app/domain"
	"budget-app/interface/service"
	"budget-app/utils"

	"github.com/gin-gonic/gin"
)

type DebtController struct {
	svc *service.DebtService
}

func NewDebt(svc *service.DebtService) *DebtController {
	return &DebtController{svc: svc}
}

func (c *DebtController) GetAll(ctx *gin.Context) {
	debts, err := c.svc.GetAll(ctx.Request.Context())
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.EnsureNotNil(debts))
}

func (c *DebtController) Create(ctx *gin.Context) {
	var d domain.Debt
	if err := ctx.ShouldBindJSON(&d); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	if err := c.svc.Create(ctx.Request.Context(), &d); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, d)
}

func (c *DebtController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	var d domain.Debt
	if err := ctx.ShouldBindJSON(&d); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	d.ID = id
	if err := c.svc.Update(ctx.Request.Context(), &d); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, d)
}

func (c *DebtController) Delete(ctx *gin.Context) {
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

func (c *DebtController) GetPayments(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	payments, err := c.svc.GetPayments(ctx.Request.Context(), id)
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.EnsureNotNil(payments))
}

func (c *DebtController) RecordPayment(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	var p domain.DebtPayment
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	if err := c.svc.RecordPayment(ctx.Request.Context(), id, &p); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, p)
}

func (c *DebtController) DeletePayment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errBadRequest(ctx, "invalid id")
		return
	}
	if err := c.svc.DeletePayment(ctx.Request.Context(), id); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted": true})
}
