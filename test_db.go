//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

"shutterspace/internal/repository"
)

func main() {
	dsn := "postgres://postgres:Zevagion123@localhost:5432/shutterspace?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	repo := repository.NewStudioRepository(db)
	studios, total, err := repo.FindAll(context.Background(), repository.StudioFilter{})
	if err != nil {
		log.Fatalf("FindAll error: %v", err)
	}

	fmt.Printf("Success! Total studios: %d\n", total)
	for _, s := range studios {
		fmt.Printf("- %s\n", s.Name)
	}
}
