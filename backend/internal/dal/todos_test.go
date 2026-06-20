package dal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListTodos_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at"}).
		AddRow(1, "Buy groceries", false, now).
		AddRow(2, "Walk the dog", true, now)

	mock.ExpectQuery(`SELECT id, title, done, created_at FROM todos ORDER BY id`).
		WillReturnRows(rows)

	todos, err := ListTodos(context.Background(), db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}
	if todos[0].Title != "Buy groceries" || todos[0].Done != false {
		t.Errorf("unexpected first todo: %+v", todos[0])
	}
	if todos[1].Title != "Walk the dog" || todos[1].Done != true {
		t.Errorf("unexpected second todo: %+v", todos[1])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestListTodos_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT id, title, done, created_at FROM todos ORDER BY id`).
		WillReturnError(errQueryFailed)

	_, err = ListTodos(context.Background(), db)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

var errQueryFailed = fmt.Errorf("connection refused")
