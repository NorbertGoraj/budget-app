package handlers

import (
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"budget-app/models"

	"github.com/gin-gonic/gin"
)

func ImportCSV(c *gin.Context) {
	accountIDStr := c.PostForm("account_id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid CSV: " + err.Error()})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV must have a header row and at least one data row"})
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

		exists, err := models.TransactionExists(date, description, math.Abs(amount))
		if err != nil {
			skipped++
			continue
		}
		if exists {
			skipped++
			continue
		}

		txType := "income"
		if amount < 0 {
			txType = "expense"
			amount = math.Abs(amount)
		}

		t := models.Transaction{
			AccountID:   accountID,
			Amount:      amount,
			Description: description,
			Category:    "Uncategorized",
			Type:        txType,
			Date:        date,
			Imported:    true,
		}
		if err := models.CreateTransaction(&t); err != nil {
			skipped++
			continue
		}
		imported++
	}

	c.JSON(http.StatusOK, gin.H{
		"imported": imported,
		"skipped":  skipped,
		"message":  fmt.Sprintf("Imported %d transactions, skipped %d", imported, skipped),
	})
}
