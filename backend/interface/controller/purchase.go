package controller

import (
	"net/http"

	"budget-app/domain"
	"budget-app/interface/service"
	"budget-app/utils"

	"github.com/gin-gonic/gin"
)

type PurchaseController struct {
	svc *service.PurchaseService
}

func NewPurchase(svc *service.PurchaseService) *PurchaseController {
	return &PurchaseController{svc: svc}
}

func (c *PurchaseController) GetAll(ctx *gin.Context) {
	purchases, err := c.svc.GetAll(ctx.Request.Context())
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.EnsureNotNil(purchases))
}

func (c *PurchaseController) Create(ctx *gin.Context) {
	var p domain.PlannedPurchase
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	if err := c.svc.Create(ctx.Request.Context(), &p); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, p)
}

func (c *PurchaseController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	var p domain.PlannedPurchase
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	p.ID = id
	if err := c.svc.Update(ctx.Request.Context(), &p); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, p)
}

func (c *PurchaseController) Delete(ctx *gin.Context) {
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
