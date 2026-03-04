//go:build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	predator "github.com/Hyoshii-Farm/nursery/feature/report/predator/repositories"
)

func main() {
	dsn := "host=ep-rapid-glitter-a1apdechv.ap-southeast-1.aws.neon.tech user=dashboard_owner password=npg_pt0zihCV9Mgd dbname=dashboard port=5432 sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	repo := predator.GetRepository(db)

	date := "2026-01-01"
	predatorIDs := []uint{1, 2}

	variants, err := repo.GetVariantSummary(date, predatorIDs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Variants for IDs 1, 2: %+v\n", variants)
}
