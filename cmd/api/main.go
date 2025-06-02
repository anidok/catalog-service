package main

import (
	"catalog-service/internal/api"
	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	config.Load()
	logger.Setup(config.LogLevel(), config.LogFormat())

	client, err := opensearch.NewClient(config.OpenSearch().Host())
	if err != nil {
		logger.NonContext.Error("failed to create opensearch client: %v", err)
	}
	repo, err := repository.NewServiceRepository(client)
	if err != nil {
		logger.NonContext.Error("failed to create service repository: %v", err)
	}

	r := api.NewRouter(repo)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port()),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.NonContext.Error("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.NonContext.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.NonContext.Error("Server forced to shutdown: %v", err)
	} else {
		logger.NonContext.Info("Server exited gracefully")
	}
}
