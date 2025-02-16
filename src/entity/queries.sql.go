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

const getOneUser = `-- name: GetOneUser :one
SELECT
 ` + "`" + `id` + "`" + `,
 ` + "`" + `name` + "`" + `,
 ` + "`" + `email` + "`" + `,
 ` + "`" + `created_at` + "`" + `,
 ` + "`" + `created_by` + "`" + `,
 ` + "`" + `updated_at` + "`" + `,
 ` + "`" + `updated_by` + "`" + `,
 ` + "`" + `deleted_at` + "`" + `,
 ` + "`" + `deleted_by` + "`" + `,
 ` + "`" + `is_deleted` + "`" + `
FROM ` + "`" + `users` + "`" + ` WHERE
 ` + "`" + `id` + "`" + ` = ?
 AND ` + "`" + `is_deleted` + "`" + ` = ?
`

type GetOneUserParams struct {
	ID        int64 `db:"id" json:"id"`
	IsDeleted int8  `db:"is_deleted" json:"is_deleted"`
}

type GetOneUserRow struct {
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
}

func (q *Queries) GetOneUser(ctx context.Context, arg GetOneUserParams) (GetOneUserRow, error) {
	row := q.db.QueryRowContext(ctx, getOneUser, arg.ID, arg.IsDeleted)
	var i GetOneUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
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

const updateUser = `-- name: UpdateUser :execresult
UPDATE ` + "`" + `users` + "`" + ` SET
 ` + "`" + `name` + "`" + ` = ?,
 ` + "`" + `updated_at` + "`" + ` = ?,
 ` + "`" + `updated_by` + "`" + ` = ?,
 ` + "`" + `is_deleted` + "`" + ` = ?
WHERE ` + "`" + `id` + "`" + ` = ?
`

type UpdateUserParams struct {
	Name      string         `db:"name" json:"name"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
	IsDeleted int8           `db:"is_deleted" json:"is_deleted"`
	ID        int64          `db:"id" json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUser,
		arg.Name,
		arg.UpdatedAt,
		arg.UpdatedBy,
		arg.IsDeleted,
		arg.ID,
	)
}
