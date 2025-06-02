package main

import (
	"catalog-service/internal/config"
	"catalog-service/internal/logger"
	"catalog-service/internal/migrate"
	"flag"
	"log"

	"github.com/opensearch-project/opensearch-go/v2"
)

var (
	schemaDir = flag.String("schema-dir", "migrations", "Directory containing schema files")
)

func main() {
	flag.Parse()
	config.Load()
	logger.Setup(config.LogLevel(), config.LogFormat())
	logger.NonContext.Info("starting db migrations")

	config := opensearch.Config{
		Addresses: config.OpenSearch().Host(),
	}

	client, err := opensearch.NewClient(config)
	if err != nil {
		logger.NonContext.Errorf(err, "failed to create openSearch client:")
		panic("failed to create openSearch client")
	}

	migrator := migrate.New(client)

	log.Printf("Starting migrations from directory: %s", *schemaDir)
	if err := migrator.Run(*schemaDir); err != nil {
		logger.NonContext.Errorf(err, "failed to run migrations")
		panic("failed to run migrations")
	}

	log.Println("migrations completed successfully")
}
