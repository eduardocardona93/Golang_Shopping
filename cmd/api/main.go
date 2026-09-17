package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eduardocardona93/golang_shopping/internal/config"
	"github.com/eduardocardona93/golang_shopping/internal/database"
	"github.com/eduardocardona93/golang_shopping/internal/handler"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
	"github.com/eduardocardona93/golang_shopping/internal/router"
	"github.com/eduardocardona93/golang_shopping/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)

	userSvc := service.NewUserService(userRepo)
	productSvc := service.NewProductService(productRepo)
	inventorySvc := service.NewInventoryService(inventoryRepo, productRepo)
	purchaseSvc := service.NewPurchaseService(db, purchaseRepo, inventoryRepo, productRepo, userRepo)

	handlers := router.Handlers{
		User:      handler.NewUserHandler(userSvc),
		Product:   handler.NewProductHandler(productSvc),
		Inventory: handler.NewInventoryHandler(inventorySvc),
		Purchase:  handler.NewPurchaseHandler(purchaseSvc),
	}

	e := router.New(handlers)

	go func() {
		addr := ":" + cfg.ServerPort
		log.Printf("starting server on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("failed to gracefully shut down server: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	log.Println("server stopped")
}
