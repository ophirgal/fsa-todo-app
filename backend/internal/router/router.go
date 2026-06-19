package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"fsa-boilerplate/backend/internal/handlers"
	"fsa-boilerplate/backend/internal/middleware"
)

func New(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
		// Add routes here
	}

	return r
}
