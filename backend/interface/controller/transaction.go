package controller

import (
	"net/http"

	"budget-app/domain"
	"budget-app/interface/service"
	"budget-app/utils"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	svc *service.TransactionService
}

func NewTransaction(svc *service.TransactionService) *TransactionController {
	return &TransactionController{svc: svc}
}

func (c *TransactionController) GetAll(ctx *gin.Context) {
	f := domain.TransactionFilter{
		Month:     ctx.Query("month"),
		AccountID: ctx.Query("account_id"),
		Category:  ctx.Query("category"),
	}
	txns, err := c.svc.GetAll(ctx.Request.Context(), f)
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.EnsureNotNil(txns))
}

func (c *TransactionController) Create(ctx *gin.Context) {
	var t domain.Transaction
	if err := ctx.ShouldBindJSON(&t); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	if err := c.svc.Create(ctx.Request.Context(), &t); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, t)
}

func (c *TransactionController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	var t domain.Transaction
	if err := ctx.ShouldBindJSON(&t); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	t.ID = id
	if err := c.svc.Update(ctx.Request.Context(), &t); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, t)
}

func (c *TransactionController) Delete(ctx *gin.Context) {
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
