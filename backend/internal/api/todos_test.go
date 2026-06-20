package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestListTodosHandler_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery(`SELECT id, title, done, created_at FROM todos ORDER BY id`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "done", "created_at"}).
			AddRow(1, "Buy groceries", false, now).
			AddRow(2, "Walk the dog", true, now))

	r := gin.New()
	r.GET("/api/todos", ListTodos(db))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var todos []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &todos); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}
	if todos[0]["title"] != "Buy groceries" {
		t.Errorf("unexpected title: %v", todos[0]["title"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestListTodosHandler_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT id, title, done, created_at FROM todos ORDER BY id`).
		WillReturnError(fmt.Errorf("connection refused"))

	r := gin.New()
	r.GET("/api/todos", ListTodos(db))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Errorf("expected 'error' key in response body")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
