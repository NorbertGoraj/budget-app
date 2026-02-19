package controller

import "github.com/gin-gonic/gin"

func SetupRoutes(
	r *gin.Engine,
	accounts *AccountController,
	transactions *TransactionController,
	budgets *BudgetController,
	purchases *PurchaseController,
	investments *InvestmentController,
	dashboard *DashboardController,
	imports *ImportController,
) {
	api := r.Group("/api")
	{
		api.GET("/accounts", accounts.GetAll)
		api.POST("/accounts", accounts.Create)
		api.PUT("/accounts/:id", accounts.Update)
		api.DELETE("/accounts/:id", accounts.Delete)

		api.GET("/transactions", transactions.GetAll)
		api.POST("/transactions", transactions.Create)
		api.PUT("/transactions/:id", transactions.Update)
		api.DELETE("/transactions/:id", transactions.Delete)

		api.POST("/import/csv", imports.ImportCSV)

		api.GET("/budgets", budgets.GetAll)
		api.POST("/budgets", budgets.Create)
		api.PUT("/budgets/:id", budgets.Update)
		api.DELETE("/budgets/:id", budgets.Delete)

		api.GET("/purchases", purchases.GetAll)
		api.POST("/purchases", purchases.Create)
		api.PUT("/purchases/:id", purchases.Update)
		api.DELETE("/purchases/:id", purchases.Delete)

		api.GET("/investments", investments.GetAll)
		api.POST("/investments", investments.Create)
		api.PUT("/investments/:id", investments.Update)
		api.DELETE("/investments/:id", investments.Delete)

		api.GET("/dashboard", dashboard.Get)
	}
}
