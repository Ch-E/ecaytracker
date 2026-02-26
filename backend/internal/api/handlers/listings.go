// Package handlers contains HTTP request handlers for the API.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	appdb "ecaytracker/backend/internal/db"
	"ecaytracker/backend/models"
)

// Listings handles GET /api/listings.
// Returns all active listings as { "data": [...], "error": null }.
func Listings(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		listings, err := appdb.GetListings(c.Request.Context(), pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":  nil,
				"error": err.Error(),
			})
			return
		}

		// Return an empty array rather than null when there are no listings.
		if listings == nil {
			listings = make([]models.Listing, 0)
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  listings,
			"error": nil,
		})
	}
}
