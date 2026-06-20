package dal

import (
	"context"
	"database/sql"
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

func TestDeleteTodo_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM todos WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := DeleteTodo(context.Background(), db, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM todos WHERE id = \$1`).
		WithArgs(int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = DeleteTodo(context.Background(), db, 99)
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestUpdateTodoDone_HappyPath(t *testing.T) {
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

	todo, err := UpdateTodoDone(context.Background(), db, 1, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.ID != 1 || todo.Done != true {
		t.Errorf("unexpected todo: %+v", todo)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestUpdateTodoDone_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`UPDATE todos SET done = \$1 WHERE id = \$2 RETURNING id, title, done, created_at`).
		WithArgs(true, int64(99)).
		WillReturnError(sql.ErrNoRows)

	_, err = UpdateTodoDone(context.Background(), db, 99, true)
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestCreateTodo_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery(`INSERT INTO todos \(title\) VALUES \(\$1\) RETURNING id, title, done, created_at`).
		WithArgs("Buy milk").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "done", "created_at"}).
			AddRow(1, "Buy milk", false, now))

	todo, err := CreateTodo(context.Background(), db, "Buy milk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.ID != 1 || todo.Title != "Buy milk" || todo.Done != false {
		t.Errorf("unexpected todo: %+v", todo)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestCreateTodo_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`INSERT INTO todos \(title\) VALUES \(\$1\) RETURNING id, title, done, created_at`).
		WithArgs("Buy milk").
		WillReturnError(errQueryFailed)

	_, err = CreateTodo(context.Background(), db, "Buy milk")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

var errQueryFailed = fmt.Errorf("connection refused")
