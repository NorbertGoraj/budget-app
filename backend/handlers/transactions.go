package handlers

import (
	"net/http"
	"strconv"

	"budget-app/models"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {
	month := c.Query("month")
	accountID := c.Query("account_id")
	category := c.Query("category")

	txns, err := models.GetTransactions(month, accountID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if txns == nil {
		txns = []models.Transaction{}
	}
	c.JSON(http.StatusOK, txns)
}

func CreateTransaction(c *gin.Context) {
	var t models.Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.CreateTransaction(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func UpdateTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var t models.Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = id
	if err := models.UpdateTransaction(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func DeleteTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := models.DeleteTransaction(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
