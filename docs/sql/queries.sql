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
FROM `users` WHERE
 `id` = ?
 AND `is_deleted` = ?;
