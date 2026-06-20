package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"fsa-boilerplate/backend/internal/handlers"
)

func New(db *sql.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
		// Add routes here
	}

	return r
}
