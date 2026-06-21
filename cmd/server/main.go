package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"shutterspace/config"
	"shutterspace/internal/handler"
	"shutterspace/internal/repository"
	"shutterspace/internal/router"
	"shutterspace/internal/service"
)

func main() {
	// 1. Load konfigurasi dari .env
	cfg := config.Load()

	// 2. Setup Gin mode
	gin.SetMode(cfg.GinMode)

	// 3. Koneksi ke PostgreSQL via GORM
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("FATAL: Gagal koneksi ke database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("FATAL: Gagal mendapatkan SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// Connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	log.Println("✅ Koneksi database berhasil")

	// 4. Inisialisasi semua repository
	userRepo        := repository.NewUserRepository(db)
	studioRepo      := repository.NewStudioRepository(db)
	studioTypeRepo  := repository.NewStudioTypeRepository(db)
	slotRepo        := repository.NewSlotRepository(db)
	bookingRepo     := repository.NewBookingRepository(db)
	paymentRepo     := repository.NewPaymentRepository(db)

	// 5. Inisialisasi semua service
	authSvc         := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	studioSvc       := service.NewStudioService(studioRepo, studioTypeRepo)
	availabilitySvc := service.NewAvailabilityService(slotRepo, bookingRepo)
	bookingSvc      := service.NewBookingService(bookingRepo, studioRepo, slotRepo, paymentRepo, db)
	paymentSvc      := service.NewPaymentService(paymentRepo, bookingRepo)

	// 6. Inisialisasi semua handler
	handlers := router.Handlers{
		Auth:         handler.NewAuthHandler(authSvc),
		Studio:       handler.NewStudioHandler(studioSvc),
		Availability: handler.NewAvailabilityHandler(availabilitySvc),
		Booking:      handler.NewBookingHandler(bookingSvc),
		Payment:      handler.NewPaymentHandler(paymentSvc),
	}

	// 7. Setup Gin engine dan routes
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	router.Setup(r, handlers, cfg.JWTSecret, cfg.AllowedOrigin)

	// 8. Jalankan server
	addr := ":" + cfg.Port
	log.Printf("🚀 Shutterspace API berjalan di http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("FATAL: Server gagal berjalan: %v", err)
	}
}
