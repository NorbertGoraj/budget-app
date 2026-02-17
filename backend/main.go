package main

import (
	"log"

	"budget-app/db"
	"budget-app/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/accounts", handlers.GetAccounts)
		api.POST("/accounts", handlers.CreateAccount)
		api.PUT("/accounts/:id", handlers.UpdateAccount)
		api.DELETE("/accounts/:id", handlers.DeleteAccount)

		api.GET("/transactions", handlers.GetTransactions)
		api.POST("/transactions", handlers.CreateTransaction)
		api.PUT("/transactions/:id", handlers.UpdateTransaction)
		api.DELETE("/transactions/:id", handlers.DeleteTransaction)

		api.POST("/import/csv", handlers.ImportCSV)

		api.GET("/budgets", handlers.GetBudgets)
		api.POST("/budgets", handlers.CreateBudget)
		api.PUT("/budgets/:id", handlers.UpdateBudget)
		api.DELETE("/budgets/:id", handlers.DeleteBudget)

		api.GET("/purchases", handlers.GetPurchases)
		api.POST("/purchases", handlers.CreatePurchase)
		api.PUT("/purchases/:id", handlers.UpdatePurchase)
		api.DELETE("/purchases/:id", handlers.DeletePurchase)

		api.GET("/investments", handlers.GetInvestments)
		api.POST("/investments", handlers.CreateInvestment)
		api.PUT("/investments/:id", handlers.UpdateInvestment)
		api.DELETE("/investments/:id", handlers.DeleteInvestment)

		api.GET("/dashboard", handlers.GetDashboard)
	}

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
