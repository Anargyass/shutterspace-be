//go:build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgres://postgres:Zevagion123@localhost:5432/shutterspace?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	files := []string{
		"migrations/001_create_enums.sql",
		"migrations/002_create_users.sql",
		"migrations/003_create_studio_types.sql",
		"migrations/004_create_studios.sql",
		"migrations/005_create_availability_slots.sql",
		"migrations/006_create_bookings.sql",
		"migrations/007_create_payments.sql",
		"migrations/seed/001_seed_types.sql",
		"migrations/seed/002_seed_studios_and_slots.sql",
		"migrations/seed/003_seed_users.sql",
	}

	for _, file := range files {
		fmt.Printf("Executing %s...\n", file)
		content, err := ioutil.ReadFile(filepath.Join(".", file))
		if err != nil {
			log.Fatalf("failed to read file %s: %v", file, err)
		}
		
		if err := db.Exec(string(content)).Error; err != nil {
			log.Fatalf("failed to execute %s: %v", file, err)
		}
	}

	fmt.Println("Migrasi dan Seeding selesai sukses!")
}
