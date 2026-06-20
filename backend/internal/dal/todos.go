package dal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

func ListTodos(ctx context.Context, db *sql.DB) ([]Todo, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, title, done, created_at FROM todos ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return todos, nil
}

func UpdateTodoDone(ctx context.Context, db *sql.DB, id int64, done bool) (*Todo, error) {
	var t Todo
	err := db.QueryRowContext(ctx,
		`UPDATE todos SET done = $1 WHERE id = $2 RETURNING id, title, done, created_at`,
		done, id,
	).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("update todo done: %w", err)
	}
	return &t, nil
}

func DeleteTodo(ctx context.Context, db *sql.DB, id int64) error {
	result, err := db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
