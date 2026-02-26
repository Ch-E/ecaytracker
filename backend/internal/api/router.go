package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"ecaytracker/backend/internal/api/handlers"
)

// NewRouter creates and configures the Gin engine with all routes and middleware.
func NewRouter(pool *pgxpool.Pool, frontendURL string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// ── CORS ──
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", frontendURL)
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// ── Routes ──
	r.GET("/health", handlers.Health(pool))

	api := r.Group("/api")
	{
		api.GET("/listings", handlers.Listings(pool))
	}

	return r
}
