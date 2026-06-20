package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func New(db *sql.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", Health)
		api.GET("/todos", ListTodos(db))
		api.PATCH("/todos/:id/done", UpdateTodoDone(db))
		api.DELETE("/todos/:id", DeleteTodo(db))
	}

	return r
}
