package appcontext

import (
	"context"
	"log"
	"net/http"
	"sync"

	"budget-app/infrastructure"
	"budget-app/infrastructure/httphandler"
	"budget-app/infrastructure/metrics"
	"budget-app/interface/controller"
	"budget-app/interface/repository"
	"budget-app/interface/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AppContext holds the application lifecycle handles.
type AppContext struct {
	context.Context
	CancelF   context.CancelFunc
	WaitGroup *sync.WaitGroup
}

// NewContext wires all dependencies, starts the HTTP server, and registers
// a graceful-shutdown coordinator. Callers own the returned context:
//
//	ctx, err := appcontext.NewContext()
//	defer ctx.CancelF()
func NewContext() (*AppContext, error) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	if err := infrastructure.Connect(); err != nil {
		cancel()
		return nil, err
	}

	db := infrastructure.DB

	// repositories
	accountRepo := repository.NewAccount(db)
	transactionRepo := repository.NewTransaction(db)
	budgetRepo := repository.NewBudget(db)
	purchaseRepo := repository.NewPurchase(db)
	investmentRepo := repository.NewInvestment(db)

	// services
	accountSvc := service.NewAccount(accountRepo)
	transactionSvc := service.NewTransaction(transactionRepo, accountRepo)
	budgetSvc := service.NewBudget(budgetRepo)
	purchaseSvc := service.NewPurchase(purchaseRepo)
	investmentSvc := service.NewInvestment(investmentRepo)
	dashboardSvc := service.NewDashboard(accountRepo, transactionRepo, budgetRepo, purchaseRepo, investmentRepo)

	// router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	r.Use(metrics.Middleware())
	controller.SetupRoutes(r,
		controller.NewAccount(accountSvc),
		controller.NewTransaction(transactionSvc),
		controller.NewBudget(budgetSvc),
		controller.NewPurchase(purchaseSvc),
		controller.NewInvestment(investmentSvc),
		controller.NewDashboard(dashboardSvc),
		controller.NewImport(transactionSvc),
	)

	// API server
	apiSrv := &http.Server{Addr: ":8080", Handler: r}
	go httphandler.New(apiSrv, ctx, wg).Serve()

	// metrics server
	metricsSrv := metrics.NewServer(":9091")
	go httphandler.New(metricsSrv, ctx, wg).Serve()

	// DB teardown — runs after all HTTP handlers have finished
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		infrastructure.Close()
		log.Println("shutdown complete")
	}()

	return &AppContext{
		Context:   ctx,
		CancelF:   cancel,
		WaitGroup: wg,
	}, nil
}
