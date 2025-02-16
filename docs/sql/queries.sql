-- name: CreateUser :execresult
INSERT INTO `users` (`name`, `password`, `email`, `created_at`, `created_by`)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateUser :execresult
UPDATE `users` SET
 `name` = ?,
 `updated_at` = ?,
 `updated_by` = ?,
 `is_deleted` = ?
WHERE `id` = ?;

-- name: GetOneUser :one
SELECT
 `id`,
 `name`,
 `email`,
 `created_at`,
 `created_by`,
 `updated_at`,
 `updated_by`,
 `deleted_at`,
 `deleted_by`,
 `is_deleted`
FROM
 `users`
WHERE
 `id` = ?
 AND `is_deleted` = ?;

-- name: ListUser :many
SELECT
 `id`,
 `name`,
 `email`,
 `created_at`,
 `created_by`,
 `updated_at`,
 `updated_by`,
 `deleted_at`,
 `deleted_by`,
 `is_deleted`
FROM
 `users`
WHERE
 (
  `name` LIKE CONCAT("%", ? , "%")
  OR `email` LIKE CONCAT("%", ?, "%")
 )
 AND `is_deleted` = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;

-- name: CreateTodo :execresult
INSERT INTO `todos` (
 `user_id`,
 `title`,
 `description`,
 `status`,
 `created_at`,
 `created_by`
) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateTodo :execresult
UPDATE `todos` SET
 `title` = ?,
 `description` = ?,
 `status` = ?,
 `updated_at` = ?,
 `updated_by` = ?,
 `is_deleted` = ?
WHERE `id` = ?;

-- name: GetOneTodo :one
SELECT
 `id`,
 `user_id`,
 `title`,
 `description`,
 `status`,
 `created_at`,
 `created_by`,
 `updated_at`,
 `updated_by`,
 `deleted_at`,
 `deleted_by`,
 `is_deleted`
FROM
 `todos`
WHERE
 `id` = ?
 AND `is_deleted` = ?;

-- name: ListTodo :many
SELECT
 `id`,
 `user_id`,
 `title`,
 `description`,
 `status`,
 `created_at`,
 `created_by`,
 `updated_at`,
 `updated_by`,
 `deleted_at`,
 `deleted_by`,
 `is_deleted`
FROM
 `todos`
WHERE
  `user_id` = ?
  AND `status` LIKE ?
  AND (
   `title` LIKE CONCAT("%", ?, "%")
   OR `description` LIKE CONCAT("%", ?, "%") 
  );
