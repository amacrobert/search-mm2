package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"search-mm2/backend/internal/api"
	"search-mm2/backend/internal/config"
	"search-mm2/backend/internal/database"
	"search-mm2/backend/internal/scraper"
	"search-mm2/backend/migrations"
)

func main() {
	cfg := config.Load()

	log.Println("running database migrations...")
	if err := database.RunMigrations(cfg.DatabaseURL, migrations.FS); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	db := database.NewQueries(pool)
	scraperService := scraper.NewService(db)

	router := api.NewRouter(cfg, db, scraperService)

	go func() {
		ticker := time.NewTicker(cfg.ScrapeInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				scraperService.RunAll(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		server.Shutdown(shutdownCtx)
	}()

	log.Printf("server listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
