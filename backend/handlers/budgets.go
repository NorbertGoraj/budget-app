package handlers

import (
	"net/http"
	"strconv"

	"budget-app/models"

	"github.com/gin-gonic/gin"
)

func GetBudgets(c *gin.Context) {
	budgets, err := models.GetAllBudgets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if budgets == nil {
		budgets = []models.Budget{}
	}
	c.JSON(http.StatusOK, budgets)
}

func CreateBudget(c *gin.Context) {
	var b models.Budget
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.CreateBudget(&b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func UpdateBudget(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var b models.Budget
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.ID = id
	if err := models.UpdateBudget(&b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func DeleteBudget(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := models.DeleteBudget(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
