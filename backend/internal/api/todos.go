package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"fsa-boilerplate/backend/internal/dal"
)

func ListTodos(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := dal.ListTodos(c.Request.Context(), db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}
