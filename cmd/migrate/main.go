package main

import (
	"log"

	"github.com/koha90/avito-masspost/internal/config"
	"github.com/koha90/avito-masspost/internal/migrator"
)

func main() {
	cfg, err := config.Load(config.Path())
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := migrator.MigratePostgres(cfg.Database.DSN(), cfg.Migration.Path); err != nil {
		log.Fatalf("migrate postgres: %v", err)
	}

	log.Printf("migrations applied from %s", cfg.Migration.Path)
}
