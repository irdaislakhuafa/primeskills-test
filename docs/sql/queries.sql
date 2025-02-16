-- name: CreateUser :execresult
INSERT INTO `users` (`name`, `email`, `created_at`, `created_by`)
VALUES (?, ?, ?, ?);
