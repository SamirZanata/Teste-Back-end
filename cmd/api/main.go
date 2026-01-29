package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/back-end/quote-api/internal/client"
	"github.com/back-end/quote-api/internal/config"
	"github.com/back-end/quote-api/internal/handler"
	"github.com/back-end/quote-api/internal/repository"
	"github.com/back-end/quote-api/internal/service"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatalf("conectar ao banco: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping banco: %v", err)
	}

	quoteRepo := repository.NewPostgresQuoteRepository(pool)
	if err := quoteRepo.EnsureSchema(ctx); err != nil {
		log.Fatalf("criar schema: %v", err)
	}

	frClient := client.NewFreteRapidoClient(
		cfg.FreteRapido.BaseURL,
		cfg.FreteRapido.Token,
		cfg.FreteRapido.PlatformCode,
		cfg.FreteRapido.ShipperCNPJ,
		cfg.FreteRapido.DispatcherCEP,
	)

	quoteSvc := service.NewQuoteService(quoteRepo, frClient)
	metricsSvc := service.NewMetricsService(quoteRepo)

	quoteH := handler.NewQuoteHandler(quoteSvc)
	metricsH := handler.NewMetricsHandler(metricsSvc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.POST("/quote", quoteH.CreateQuote)
	r.GET("/metrics", metricsH.GetMetrics)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("API ouvindo em :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("servidor encerrado")
}
