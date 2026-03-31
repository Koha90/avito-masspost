package main

import (
	"log"

	"github.com/koha90/avito-masspost/internal/config"
	"github.com/koha90/avito-masspost/internal/feed"
)

func main() {
	cfg, err := config.Load(config.Path())
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	listings := []feed.Listing{
		{
			ID:              "beton-m300-msk-001",
			Title:           "Товарный бетон М300 с доставкой",
			Description:     "Поставляем товарный бетон М300 частным лицам. Доставка миксером, помощь с подбором объёма, быстрая отгрузка.",
			Category:        "Ремонт и строительство",
			OperationType:   "Продам",
			Address:         "Московская область",
			Price:           5400,
			ContactPhone:    "+79990000000",
			ManagerName:     "Алексей",
			City:            "Москва",
			Region:          "Московская область",
			Images:          []string{"https://example.com/beton-1.jpg", "https://example.com/beton-2.jpg"},
			ConcreteGrade:   "М300",
			ConcreteClass:   "B22.5",
			Mobility:        "П3",
			FrostResistance: "F200",
			WaterResistance: "W8",
			MinVolumeM3:     1,
			Delivery:        true,
		},
	}

	if err := feed.Write(cfg.Feed.OutputPath, listings); err != nil {
		log.Fatalf("write feed: %v", err)
	}

	log.Printf("feed written to %s", cfg.Feed.OutputPath)
}
