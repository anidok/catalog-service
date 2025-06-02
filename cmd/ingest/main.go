package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"os"
	"time"

	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"catalog-service/internal/models"
	"catalog-service/internal/opensearch"
	"catalog-service/internal/repository"
)

var (
	dataFile = flag.String("data-file", "data.jsonl", "Path to the JSONL data file")
)

func main() {
	flag.Parse()
	config.Load()
	logger.Setup(config.LogLevel(), config.LogFormat())
	logger.NonContext.Info("starting db ingestion")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	serviceRepo, err := loadServiceRepo()
	if err != nil {
		logger.NonContext.Errorf(err, "failed to initialize service repository")
		return
	}

	count, err := processFile(ctx, *dataFile, serviceRepo)
	if err != nil {
		logger.NonContext.Errorf(err, "failed to process file: %s", *dataFile)
		return
	}

	logger.NonContext.Infof("data ingestion completed. successfully indexed %d documents.", count)
}

func loadServiceRepo() (repository.ServiceRepository, error) {
	client, err := opensearch.NewClient(config.OpenSearch().Host())
	if err != nil {
		return nil, err
	}
	return repository.NewServiceRepository(client)
}

func processFile(ctx context.Context, filePath string, repo repository.ServiceRepository) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var count int
	for scanner.Scan() {
		if err := processLine(ctx, scanner.Bytes(), count+1, repo); err != nil {
			logger.NonContext.Errorf(err, "Failed to process line %d: %v", count+1, err)
			continue
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return count, err
	}
	return count, nil
}

func processLine(ctx context.Context, line []byte, lineNum int, repo repository.ServiceRepository) error {
	if len(line) == 0 {
		return nil
	}

	var service models.Service
	if err := json.Unmarshal(line, &service); err != nil {
		return err
	}

	now := time.Now().UTC()
	if service.CreatedAt.IsZero() {
		service.CreatedAt = now
	}
	service.UpdatedAt = now

	if err := repo.Create(ctx, &service); err != nil {
		return err
	}

	logger.NonContext.Infof("Successfully indexed service from line %d", lineNum)
	return nil
}
