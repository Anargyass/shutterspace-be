package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBUrl          string
	JWTSecret      string
	AllowedOrigin  string
	JWTExpiryHours int
	GinMode        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: .env file tidak ditemukan, menggunakan environment variables")
	}

	// SECURITY: JWT_SECRET_KEY wajib dari env — tidak ada hardcoded fallback
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY environment variable wajib diisi. Generate dengan: openssl rand -hex 32")
	}

	dbURL := os.Getenv("DB_DSN")
	if dbURL == "" {
		log.Fatal("FATAL: DB_DSN environment variable wajib diisi")
	}

	expiryHours := 24
	if h := os.Getenv("JWT_EXPIRY_HOURS"); h != "" {
		if v, err := strconv.Atoi(h); err == nil && v > 0 {
			expiryHours = v
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}

	return &Config{
		Port:           port,
		DBUrl:          dbURL,
		JWTSecret:      jwtSecret,
		AllowedOrigin:  allowedOrigin,
		JWTExpiryHours: expiryHours,
		GinMode:        ginMode,
	}
}
