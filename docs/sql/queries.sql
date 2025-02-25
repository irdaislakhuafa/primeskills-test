-- name: CreateUser :execresult
INSERT INTO `users` (`name`, `password`, `email`, `created_at`, `created_by`)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateUser :execresult
UPDATE `users` SET
 `name` = ?,
 `updated_at` = ?,
 `updated_by` = ?,
 `deleted_at` = ?,
 `deleted_by` = ?,
 `is_deleted` = ?
WHERE `id` = ?;

-- name: ChangePasswordUser :execresult
UPDATE `users` SET
 `password` = ?,
 `updated_at` = ?,
 `updated_by` = ?
WHERE `id` = ?;

-- name: GetOneUser :one
SELECT
 `id`,
 `name`,
 `email`,
 `password`,
 `is_active`,
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
 (`id` = ? OR `email` = ?)
 AND `is_deleted` = ?;

-- name: ListUser :many
SELECT
 `id`,
 `name`,
 `email`,
 `is_active`,
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

-- name: CountUser :one
SELECT
 COUNT(`id`)
FROM
 `users`
WHERE
 (
  `name` LIKE CONCAT("%", ? , "%")
  OR `email` LIKE CONCAT("%", ?, "%")
 )
 AND `is_deleted` = ?;
 
-- name: UpdateActivationUser :execresult
UPDATE `users` SET `is_active` = ? WHERE `id` = ?;


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
 `deleted_at` = ?,
 `deleted_by` = ?,
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
  AND `is_deleted` = ?
  AND (
   `title` LIKE CONCAT("%", ?, "%")
   OR `description` LIKE CONCAT("%", ?, "%") 
  )
ORDER BY `id` DESC
LIMIT ?
OFFSET ?;

-- name: CountTodo :one
SELECT
 COUNT(`id`) AS total
FROM
 `todos`
WHERE
  `user_id` = ?
  AND `status` LIKE ?
  AND `is_deleted` = ?
  AND (
   `title` LIKE CONCAT("%", ?, "%")
   OR `description` LIKE CONCAT("%", ?, "%") 
  );

-- name: CreateTodoHistory :execresult
INSERT INTO `todo_histories` (
 `todo_id`,
 `message`,
 `created_at`,
 `created_by`
) VALUES (?, ?, ?, ?);

-- name: ListTodoHistories :many
SELECT
 `id`,
 `todo_id`,
 `message`,
 `created_at`,
 `created_by`,
 `updated_at`,
 `updated_by`,
 `deleted_at`,
 `deleted_by`,
 `is_deleted`
FROM
 `todo_histories`
WHERE
 `todo_id` = ?
 AND `is_deleted` = ?
ORDER BY `id` DESC
LIMIT ?
OFFSET ?;

-- name: CountTodoHistories :one
SELECT
 COUNT(`id`) AS total
FROM
 `todo_histories`
WHERE
 `is_deleted` = ?
 AND `todo_id` = ?;

-- name: CreateOTP :execresult
INSERT INTO `otp` (
 `user_id`,
 `code`,
 `expirate_at`,
 `created_at`,
 `created_by`
) VALUES (?, ?, ?, ?, ?);

-- name: GetOneOTP :one
SELECT
 `id`,
 `user_id`,
 `code`,
 `expirate_at`,
 `is_used`,
 `created_at`,
 `created_by`,
 `updated_at`,
 `updated_by`,
 `deleted_at`,
 `deleted_by`,
 `is_deleted`
FROM
 `otp`
WHERE
 (`id` LIKE ? OR `user_id` LIKE ?)
 AND `is_deleted` = ?
 AND `expirate_at` > ?
 AND `is_used` = ?;

-- name: UpdateOTP :execresult
UPDATE `otp` SET
 `code` = ?,
 `is_used` = ?,
 `updated_at` = ?,
 `updated_by` = ?,
 `deleted_at` = ?,
 `deleted_by` = ?,
 `is_deleted` = ?
WHERE `id` = ?;
