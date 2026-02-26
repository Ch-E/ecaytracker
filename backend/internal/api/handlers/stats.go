// Package handlers contains HTTP request handlers for the API.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	appdb "ecaytracker/backend/internal/db"
)

// Stats handles GET /api/stats.
// Returns pre-computed dashboard statistics as { "data": {...}, "error": null }.
func Stats(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := appdb.GetStats(c.Request.Context(), pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":  nil,
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  stats,
			"error": nil,
		})
	}
}
