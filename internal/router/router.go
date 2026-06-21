package router

import (
	"github.com/gin-gonic/gin"

	"shutterspace/internal/handler"
	"shutterspace/internal/middleware"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	Studio       *handler.StudioHandler
	Availability *handler.AvailabilityHandler
	Booking      *handler.BookingHandler
	Payment      *handler.PaymentHandler
}

func Setup(r *gin.Engine, h Handlers, jwtSecret, allowedOrigin string) {
	// Global middlewares
	r.Use(middleware.CORSMiddleware(allowedOrigin))

	// Rate limit global: 20 req/detik, burst 50
	// Rate limit ketat untuk auth: 5 req/detik, burst 10
	authLimiter    := middleware.RateLimiter(5, 10)
	defaultLimiter := middleware.RateLimiter(20, 50)

	api := r.Group("/api/v1")
	api.Use(defaultLimiter)

	// ─── Auth ───────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Use(authLimiter)
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/logout", middleware.RequireAuth(jwtSecret), h.Auth.Logout)
		auth.GET("/me", middleware.RequireAuth(jwtSecret), h.Auth.Me)
	}

	// ─── Studios (Public) ───────────────────────────────────────
	studios := api.Group("/studios")
	{
		studios.GET("", h.Studio.ListStudios)
		studios.GET("/types", h.Studio.GetStudioTypes)
		studios.GET("/:id", h.Studio.GetStudio)
		studios.GET("/:id/availability", h.Availability.GetAvailability)
	}

	// ─── Bookings (Requires Auth) ───────────────────────────────
	bookings := api.Group("/bookings")
	bookings.Use(middleware.RequireAuth(jwtSecret))
	{
		bookings.POST("", h.Booking.CreateBooking)
		bookings.GET("", h.Booking.GetMyBookings)
		bookings.GET("/:id", h.Booking.GetBookingByID)
		bookings.PATCH("/:id/cancel", h.Booking.CancelBooking)
	}

	// ─── Payments (Requires Auth) ───────────────────────────────
	payments := api.Group("/payments")
	payments.Use(middleware.RequireAuth(jwtSecret))
	{
		payments.POST("/:bookingId/pay", h.Payment.ProcessPayment)
	}

	// ─── Admin (Requires Auth + studio_admin role) ──────────────
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAuth(jwtSecret))
	admin.Use(middleware.RequireRole("studio_admin"))
	{
		admin.GET("/studios", h.Studio.GetMyStudios)
		admin.GET("/studios/:studioId/bookings", h.Booking.GetStudioBookings)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "shutterspace-api"})
	})
}
