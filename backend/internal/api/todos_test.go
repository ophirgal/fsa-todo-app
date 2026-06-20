package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestUpdateTodoDoneHandler_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery(`UPDATE todos SET done = \$1 WHERE id = \$2 RETURNING id, title, done, created_at`).
		WithArgs(true, int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "done", "created_at"}).
			AddRow(1, "Buy groceries", true, now))

	r := gin.New()
	r.PATCH("/api/todos/:id/done", UpdateTodoDone(db))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/api/todos/1/done", strings.NewReader(`{"done":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var todo map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &todo); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if todo["done"] != true {
		t.Errorf("expected done=true, got %v", todo["done"])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestUpdateTodoDoneHandler_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`UPDATE todos SET done = \$1 WHERE id = \$2 RETURNING id, title, done, created_at`).
		WithArgs(true, int64(99)).
		WillReturnError(sql.ErrNoRows)

	r := gin.New()
	r.PATCH("/api/todos/:id/done", UpdateTodoDone(db))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/api/todos/99/done", strings.NewReader(`{"done":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestDeleteTodoHandler_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM todos WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r := gin.New()
	r.DELETE("/api/todos/:id", DeleteTodo(db))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestDeleteTodoHandler_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM todos WHERE id = \$1`).
		WithArgs(int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	r := gin.New()
	r.DELETE("/api/todos/:id", DeleteTodo(db))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/99", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
