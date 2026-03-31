package main

import (
	"context"
	"log"

	"github.com/koha90/avito-masspost/internal/avito"
	"github.com/koha90/avito-masspost/internal/config"
)

func main() {
	cfg, err := config.Load(config.Path())
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	httpClient := avito.NewHTTPClient()
	client := avito.NewClient(cfg.Avito.BaseURL, httpClient)
	auth := avito.NewAuth(client, cfg.Avito.ClientID, cfg.Avito.ClientSecret)

	token, err := auth.Token(context.Background())
	if err != nil {
		log.Fatalf("request token: %v", err)
	}

	log.Printf("token type=%s expires_at=%s", token.TokenType, token.ExpiresAt.Format("2006-01-02 15:04:05"))
}
