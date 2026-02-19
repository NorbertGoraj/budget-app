package controller

import (
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"budget-app/domain"
	"budget-app/interface/service"

	"github.com/gin-gonic/gin"
)

type ImportController struct {
	svc *service.TransactionService
}

func NewImport(svc *service.TransactionService) *ImportController {
	return &ImportController{svc: svc}
}

func (c *ImportController) ImportCSV(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	accountIDStr := ctx.PostForm("account_id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		errBadRequest(ctx, "account_id is required")
		return
	}

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		errBadRequest(ctx, "file is required")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		errBadRequest(ctx, "invalid CSV: "+err.Error())
		return
	}

	if len(records) < 2 {
		errBadRequest(ctx, "CSV must have a header row and at least one data row")
		return
	}

	var imported, skipped int
	for _, row := range records[1:] {
		if len(row) < 3 {
			skipped++
			continue
		}

		date := strings.TrimSpace(row[0])
		description := strings.TrimSpace(row[1])
		amountStr := strings.TrimSpace(row[2])

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			skipped++
			continue
		}

		exists, err := c.svc.Exists(reqCtx, date, description, math.Abs(amount))
		if err != nil || exists {
			skipped++
			continue
		}

		txType := "income"
		if amount < 0 {
			txType = "expense"
			amount = math.Abs(amount)
		}

		t := domain.Transaction{
			AccountID:   accountID,
			Amount:      amount,
			Description: description,
			Category:    "Uncategorized",
			Type:        txType,
			Date:        date,
			Imported:    true,
		}
		if err := c.svc.Create(reqCtx, &t); err != nil {
			skipped++
			continue
		}
		imported++
	}

	ctx.JSON(http.StatusOK, gin.H{
		"imported": imported,
		"skipped":  skipped,
		"message":  fmt.Sprintf("Imported %d transactions, skipped %d", imported, skipped),
	})
}
