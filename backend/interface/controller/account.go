package controller

import (
	"net/http"

	"budget-app/domain"
	"budget-app/interface/service"
	"budget-app/utils"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	svc *service.AccountService
}

func NewAccount(svc *service.AccountService) *AccountController {
	return &AccountController{svc: svc}
}

func (c *AccountController) GetAll(ctx *gin.Context) {
	accounts, err := c.svc.GetAll(ctx.Request.Context())
	if err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, utils.EnsureNotNil(accounts))
}

func (c *AccountController) Create(ctx *gin.Context) {
	var a domain.Account
	if err := ctx.ShouldBindJSON(&a); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	if err := c.svc.Create(ctx.Request.Context(), &a); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, a)
}

func (c *AccountController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}
	var a domain.Account
	if err := ctx.ShouldBindJSON(&a); err != nil {
		errBadRequest(ctx, err.Error())
		return
	}
	a.ID = id
	if err := c.svc.Update(ctx.Request.Context(), &a); err != nil {
		errInternal(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, a)
}

func (c *AccountController) Delete(ctx *gin.Context) {
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
