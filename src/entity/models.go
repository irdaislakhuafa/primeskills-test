// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package entity

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID int64 `db:"id" json:"id"`
	// refer to users.id
	UserID      int64          `db:"user_id" json:"user_id"`
	Title       string         `db:"title" json:"title"`
	Status      string         `db:"status" json:"status"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	CreatedBy   string         `db:"created_by" json:"created_by"`
	UpdatedAt   sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy   sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt   sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy   sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted   int8           `db:"is_deleted" json:"is_deleted"`
	Description string         `db:"description" json:"description"`
}

type TodoHistory struct {
	ID int64 `db:"id" json:"id"`
	// refer to todos.id
	TodoID    int64          `db:"todo_id" json:"todo_id"`
	Message   string         `db:"message" json:"message"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	CreatedBy string         `db:"created_by" json:"created_by"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted int8           `db:"is_deleted" json:"is_deleted"`
}

type User struct {
	ID        int64          `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Email     string         `db:"email" json:"email"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	CreatedBy string         `db:"created_by" json:"created_by"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted int8           `db:"is_deleted" json:"is_deleted"`
	Password  string         `db:"password" json:"password"`
	IsActive  int8           `db:"is_active" json:"is_active"`
}
