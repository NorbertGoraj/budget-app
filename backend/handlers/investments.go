package handlers

import (
	"net/http"
	"strconv"

	"budget-app/models"

	"github.com/gin-gonic/gin"
)

func GetInvestments(c *gin.Context) {
	investments, err := models.GetAllInvestments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if investments == nil {
		investments = []models.Investment{}
	}
	c.JSON(http.StatusOK, investments)
}

func CreateInvestment(c *gin.Context) {
	var inv models.Investment
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if inv.Status == "" {
		inv.Status = "planned"
	}
	if err := models.CreateInvestment(&inv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, inv)
}

func UpdateInvestment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var inv models.Investment
	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv.ID = id
	if err := models.UpdateInvestment(&inv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func DeleteInvestment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := models.DeleteInvestment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
