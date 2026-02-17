package handlers

import (
	"net/http"
	"strconv"

	"budget-app/models"

	"github.com/gin-gonic/gin"
)

func GetPurchases(c *gin.Context) {
	purchases, err := models.GetAllPurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if purchases == nil {
		purchases = []models.PlannedPurchase{}
	}
	c.JSON(http.StatusOK, purchases)
}

func CreatePurchase(c *gin.Context) {
	var p models.PlannedPurchase
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if p.Status == "" {
		p.Status = "planned"
	}
	if err := models.CreatePurchase(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func UpdatePurchase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var p models.PlannedPurchase
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	if err := models.UpdatePurchase(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func DeletePurchase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := models.DeletePurchase(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
