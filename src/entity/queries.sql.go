// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package entity

import (
	"context"
	"database/sql"
	"time"
)

const countTodo = `-- name: CountTodo :one
SELECT
 COUNT(` + "`" + `id` + "`" + `) AS total
FROM
 ` + "`" + `todos` + "`" + `
WHERE
  ` + "`" + `user_id` + "`" + ` = ?
  AND ` + "`" + `status` + "`" + ` LIKE ?
  AND ` + "`" + `is_deleted` + "`" + ` = ?
  AND (
   ` + "`" + `title` + "`" + ` LIKE CONCAT("%", ?, "%")
   OR ` + "`" + `description` + "`" + ` LIKE CONCAT("%", ?, "%") 
  )
`

type CountTodoParams struct {
	UserID    int64       `db:"user_id" json:"user_id"`
	Status    string      `db:"status" json:"status"`
	IsDeleted int8        `db:"is_deleted" json:"is_deleted"`
	CONCAT    interface{} `db:"CONCAT" json:"CONCAT"`
	CONCAT_2  interface{} `db:"CONCAT_2" json:"CONCAT_2"`
}

func (q *Queries) CountTodo(ctx context.Context, arg CountTodoParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTodo,
		arg.UserID,
		arg.Status,
		arg.IsDeleted,
		arg.CONCAT,
		arg.CONCAT_2,
	)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const countTodoHistories = `-- name: CountTodoHistories :one
SELECT
 COUNT(` + "`" + `id` + "`" + `) AS total
FROM
 ` + "`" + `todo_histories` + "`" + `
WHERE
 ` + "`" + `is_deleted` + "`" + ` = ?
 AND ` + "`" + `todo_id` + "`" + ` = ?
`

type CountTodoHistoriesParams struct {
	IsDeleted int8  `db:"is_deleted" json:"is_deleted"`
	TodoID    int64 `db:"todo_id" json:"todo_id"`
}

func (q *Queries) CountTodoHistories(ctx context.Context, arg CountTodoHistoriesParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTodoHistories, arg.IsDeleted, arg.TodoID)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const countUser = `-- name: CountUser :one
SELECT
 COUNT(` + "`" + `id` + "`" + `)
FROM
 ` + "`" + `users` + "`" + `
WHERE
 (
  ` + "`" + `name` + "`" + ` LIKE CONCAT("%", ? , "%")
  OR ` + "`" + `email` + "`" + ` LIKE CONCAT("%", ?, "%")
 )
 AND ` + "`" + `is_deleted` + "`" + ` = ?
`

type CountUserParams struct {
	CONCAT    interface{} `db:"CONCAT" json:"CONCAT"`
	CONCAT_2  interface{} `db:"CONCAT_2" json:"CONCAT_2"`
	IsDeleted int8        `db:"is_deleted" json:"is_deleted"`
}

func (q *Queries) CountUser(ctx context.Context, arg CountUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUser, arg.CONCAT, arg.CONCAT_2, arg.IsDeleted)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createTodo = `-- name: CreateTodo :execresult
INSERT INTO ` + "`" + `todos` + "`" + ` (
 ` + "`" + `user_id` + "`" + `,
 ` + "`" + `title` + "`" + `,
 ` + "`" + `description` + "`" + `,
 ` + "`" + `status` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `
) VALUES (?, ?, ?, ?, ?, ?)
`

type CreateTodoParams struct {
	UserID      int64     `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createTodo,
		arg.UserID,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.CreatedAt,
		arg.CreatedBy,
	)
}

const createTodoHistory = `-- name: CreateTodoHistory :execresult
INSERT INTO ` + "`" + `todo_histories` + "`" + ` (
 ` + "`" + `todo_id` + "`" + `,
 ` + "`" + `message` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `
) VALUES (?, ?, ?, ?)
`

type CreateTodoHistoryParams struct {
	TodoID    int64     `db:"todo_id" json:"todo_id"`
	Message   string    `db:"message" json:"message"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy string    `db:"created_by" json:"created_by"`
}

func (q *Queries) CreateTodoHistory(ctx context.Context, arg CreateTodoHistoryParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createTodoHistory,
		arg.TodoID,
		arg.Message,
		arg.CreatedAt,
		arg.CreatedBy,
	)
}

const createUser = `-- name: CreateUser :execresult
INSERT INTO ` + "`" + `users` + "`" + ` (` + "`" + `name` + "`" + `, ` + "`" + `password` + "`" + `, ` + "`" + `email` + "`" + `, ` + "`" + `created_at` + "`" + `, ` + "`" + `created_by` + "`" + `)
VALUES (?, ?, ?, ?, ?)
`

type CreateUserParams struct {
	Name      string    `db:"name" json:"name"`
	Password  string    `db:"password" json:"password"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy string    `db:"created_by" json:"created_by"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.Name,
		arg.Password,
		arg.Email,
		arg.CreatedAt,
		arg.CreatedBy,
	)
}

const getOneTodo = `-- name: GetOneTodo :one
SELECT
 ` + "`" + `id` + "`" + `,
 ` + "`" + `user_id` + "`" + `,
 ` + "`" + `title` + "`" + `,
 ` + "`" + `description` + "`" + `,
 ` + "`" + `status` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `,
 ` + "`" + `updated_at` + "`" + `,
 ` + "`" + `updated_by` + "`" + `,
 ` + "`" + `deleted_at` + "`" + `,
 ` + "`" + `deleted_by` + "`" + `,
 ` + "`" + `is_deleted` + "`" + `
FROM
 ` + "`" + `todos` + "`" + `
WHERE
 ` + "`" + `id` + "`" + ` = ?
 AND ` + "`" + `is_deleted` + "`" + ` = ?
`

type GetOneTodoParams struct {
	ID        int64 `db:"id" json:"id"`
	IsDeleted int8  `db:"is_deleted" json:"is_deleted"`
}

type GetOneTodoRow struct {
	ID          int64          `db:"id" json:"id"`
	UserID      int64          `db:"user_id" json:"user_id"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Status      string         `db:"status" json:"status"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	CreatedBy   string         `db:"created_by" json:"created_by"`
	UpdatedAt   sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy   sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt   sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy   sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted   int8           `db:"is_deleted" json:"is_deleted"`
}

func (q *Queries) GetOneTodo(ctx context.Context, arg GetOneTodoParams) (GetOneTodoRow, error) {
	row := q.db.QueryRowContext(ctx, getOneTodo, arg.ID, arg.IsDeleted)
	var i GetOneTodoRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
		&i.IsDeleted,
	)
	return i, err
}

const getOneUser = `-- name: GetOneUser :one
SELECT
 ` + "`" + `id` + "`" + `,
 ` + "`" + `name` + "`" + `,
 ` + "`" + `email` + "`" + `,
 ` + "`" + `password` + "`" + `,
 ` + "`" + `is_active` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `,
 ` + "`" + `updated_at` + "`" + `,
 ` + "`" + `updated_by` + "`" + `,
 ` + "`" + `deleted_at` + "`" + `,
 ` + "`" + `deleted_by` + "`" + `,
 ` + "`" + `is_deleted` + "`" + `
FROM
 ` + "`" + `users` + "`" + `
WHERE
 (` + "`" + `id` + "`" + ` = ? OR ` + "`" + `email` + "`" + ` = ?)
 AND ` + "`" + `is_deleted` + "`" + ` = ?
`

type GetOneUserParams struct {
	ID        int64  `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	IsDeleted int8   `db:"is_deleted" json:"is_deleted"`
}

type GetOneUserRow struct {
	ID        int64          `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Email     string         `db:"email" json:"email"`
	Password  string         `db:"password" json:"password"`
	IsActive  int8           `db:"is_active" json:"is_active"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	CreatedBy string         `db:"created_by" json:"created_by"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted int8           `db:"is_deleted" json:"is_deleted"`
}

func (q *Queries) GetOneUser(ctx context.Context, arg GetOneUserParams) (GetOneUserRow, error) {
	row := q.db.QueryRowContext(ctx, getOneUser, arg.ID, arg.Email, arg.IsDeleted)
	var i GetOneUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.DeletedAt,
		&i.DeletedBy,
		&i.IsDeleted,
	)
	return i, err
}

const listTodo = `-- name: ListTodo :many
SELECT
 ` + "`" + `id` + "`" + `,
 ` + "`" + `user_id` + "`" + `,
 ` + "`" + `title` + "`" + `,
 ` + "`" + `description` + "`" + `,
 ` + "`" + `status` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `,
 ` + "`" + `updated_at` + "`" + `,
 ` + "`" + `updated_by` + "`" + `,
 ` + "`" + `deleted_at` + "`" + `,
 ` + "`" + `deleted_by` + "`" + `,
 ` + "`" + `is_deleted` + "`" + `
FROM
 ` + "`" + `todos` + "`" + `
WHERE
  ` + "`" + `user_id` + "`" + ` = ?
  AND ` + "`" + `status` + "`" + ` LIKE ?
  AND ` + "`" + `is_deleted` + "`" + ` = ?
  AND (
   ` + "`" + `title` + "`" + ` LIKE CONCAT("%", ?, "%")
   OR ` + "`" + `description` + "`" + ` LIKE CONCAT("%", ?, "%") 
  )
ORDER BY ` + "`" + `id` + "`" + ` DESC
LIMIT ?
OFFSET ?
`

type ListTodoParams struct {
	UserID    int64       `db:"user_id" json:"user_id"`
	Status    string      `db:"status" json:"status"`
	IsDeleted int8        `db:"is_deleted" json:"is_deleted"`
	CONCAT    interface{} `db:"CONCAT" json:"CONCAT"`
	CONCAT_2  interface{} `db:"CONCAT_2" json:"CONCAT_2"`
	Limit     int32       `db:"limit" json:"limit"`
	Offset    int32       `db:"offset" json:"offset"`
}

type ListTodoRow struct {
	ID          int64          `db:"id" json:"id"`
	UserID      int64          `db:"user_id" json:"user_id"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Status      string         `db:"status" json:"status"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	CreatedBy   string         `db:"created_by" json:"created_by"`
	UpdatedAt   sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy   sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt   sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy   sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted   int8           `db:"is_deleted" json:"is_deleted"`
}

func (q *Queries) ListTodo(ctx context.Context, arg ListTodoParams) ([]ListTodoRow, error) {
	rows, err := q.db.QueryContext(ctx, listTodo,
		arg.UserID,
		arg.Status,
		arg.IsDeleted,
		arg.CONCAT,
		arg.CONCAT_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTodoRow
	for rows.Next() {
		var i ListTodoRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
			&i.DeletedAt,
			&i.DeletedBy,
			&i.IsDeleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTodoHistories = `-- name: ListTodoHistories :many
SELECT
 ` + "`" + `id` + "`" + `,
 ` + "`" + `todo_id` + "`" + `,
 ` + "`" + `message` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `,
 ` + "`" + `updated_at` + "`" + `,
 ` + "`" + `updated_by` + "`" + `,
 ` + "`" + `deleted_at` + "`" + `,
 ` + "`" + `deleted_by` + "`" + `,
 ` + "`" + `is_deleted` + "`" + `
FROM
 ` + "`" + `todo_histories` + "`" + `
WHERE
 ` + "`" + `todo_id` + "`" + ` = ?
 AND ` + "`" + `is_deleted` + "`" + ` = ?
ORDER BY ` + "`" + `id` + "`" + ` DESC
LIMIT ?
OFFSET ?
`

type ListTodoHistoriesParams struct {
	TodoID    int64 `db:"todo_id" json:"todo_id"`
	IsDeleted int8  `db:"is_deleted" json:"is_deleted"`
	Limit     int32 `db:"limit" json:"limit"`
	Offset    int32 `db:"offset" json:"offset"`
}

func (q *Queries) ListTodoHistories(ctx context.Context, arg ListTodoHistoriesParams) ([]TodoHistory, error) {
	rows, err := q.db.QueryContext(ctx, listTodoHistories,
		arg.TodoID,
		arg.IsDeleted,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TodoHistory
	for rows.Next() {
		var i TodoHistory
		if err := rows.Scan(
			&i.ID,
			&i.TodoID,
			&i.Message,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
			&i.DeletedAt,
			&i.DeletedBy,
			&i.IsDeleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUser = `-- name: ListUser :many
SELECT
 ` + "`" + `id` + "`" + `,
 ` + "`" + `name` + "`" + `,
 ` + "`" + `email` + "`" + `,
 ` + "`" + `is_active` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `,
 ` + "`" + `updated_at` + "`" + `,
 ` + "`" + `updated_by` + "`" + `,
 ` + "`" + `deleted_at` + "`" + `,
 ` + "`" + `deleted_by` + "`" + `,
 ` + "`" + `is_deleted` + "`" + `
FROM
 ` + "`" + `users` + "`" + `
WHERE
 (
  ` + "`" + `name` + "`" + ` LIKE CONCAT("%", ? , "%")
  OR ` + "`" + `email` + "`" + ` LIKE CONCAT("%", ?, "%")
 )
 AND ` + "`" + `is_deleted` + "`" + ` = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?
`

type ListUserParams struct {
	CONCAT    interface{} `db:"CONCAT" json:"CONCAT"`
	CONCAT_2  interface{} `db:"CONCAT_2" json:"CONCAT_2"`
	IsDeleted int8        `db:"is_deleted" json:"is_deleted"`
	Limit     int32       `db:"limit" json:"limit"`
	Offset    int32       `db:"offset" json:"offset"`
}

type ListUserRow struct {
	ID        int64          `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Email     string         `db:"email" json:"email"`
	IsActive  int8           `db:"is_active" json:"is_active"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	CreatedBy string         `db:"created_by" json:"created_by"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted int8           `db:"is_deleted" json:"is_deleted"`
}

func (q *Queries) ListUser(ctx context.Context, arg ListUserParams) ([]ListUserRow, error) {
	rows, err := q.db.QueryContext(ctx, listUser,
		arg.CONCAT,
		arg.CONCAT_2,
		arg.IsDeleted,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserRow
	for rows.Next() {
		var i ListUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.IsActive,
			&i.CreatedAt,
			&i.CreatedBy,
			&i.UpdatedAt,
			&i.UpdatedBy,
			&i.DeletedAt,
			&i.DeletedBy,
			&i.IsDeleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateActivationUser = `-- name: UpdateActivationUser :execresult
UPDATE ` + "`" + `users` + "`" + ` SET ` + "`" + `is_active` + "`" + ` = ? WHERE ` + "`" + `id` + "`" + ` = ?
`

type UpdateActivationUserParams struct {
	IsActive int8  `db:"is_active" json:"is_active"`
	ID       int64 `db:"id" json:"id"`
}

func (q *Queries) UpdateActivationUser(ctx context.Context, arg UpdateActivationUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateActivationUser, arg.IsActive, arg.ID)
}

const updateTodo = `-- name: UpdateTodo :execresult
UPDATE ` + "`" + `todos` + "`" + ` SET
 ` + "`" + `title` + "`" + ` = ?,
 ` + "`" + `description` + "`" + ` = ?,
 ` + "`" + `status` + "`" + ` = ?,
 ` + "`" + `updated_at` + "`" + ` = ?,
 ` + "`" + `updated_by` + "`" + ` = ?,
 ` + "`" + `deleted_at` + "`" + ` = ?,
 ` + "`" + `deleted_by` + "`" + ` = ?,
 ` + "`" + `is_deleted` + "`" + ` = ?
WHERE ` + "`" + `id` + "`" + ` = ?
`

type UpdateTodoParams struct {
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Status      string         `db:"status" json:"status"`
	UpdatedAt   sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy   sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt   sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy   sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted   int8           `db:"is_deleted" json:"is_deleted"`
	ID          int64          `db:"id" json:"id"`
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateTodo,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.UpdatedAt,
		arg.UpdatedBy,
		arg.DeletedAt,
		arg.DeletedBy,
		arg.IsDeleted,
		arg.ID,
	)
}

const updateUser = `-- name: UpdateUser :execresult
UPDATE ` + "`" + `users` + "`" + ` SET
 ` + "`" + `name` + "`" + ` = ?,
 ` + "`" + `updated_at` + "`" + ` = ?,
 ` + "`" + `updated_by` + "`" + ` = ?,
 ` + "`" + `deleted_at` + "`" + ` = ?,
 ` + "`" + `deleted_by` + "`" + ` = ?,
 ` + "`" + `is_deleted` + "`" + ` = ?
WHERE ` + "`" + `id` + "`" + ` = ?
`

type UpdateUserParams struct {
	Name      string         `db:"name" json:"name"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at" json:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by" json:"deleted_by"`
	IsDeleted int8           `db:"is_deleted" json:"is_deleted"`
	ID        int64          `db:"id" json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.Name,
		arg.UpdatedAt,
		arg.UpdatedBy,
		arg.DeletedAt,
		arg.DeletedBy,
		arg.IsDeleted,
		arg.ID,
	)
}
