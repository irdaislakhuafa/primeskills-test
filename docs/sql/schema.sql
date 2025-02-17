CREATE DATABASE IF NOT EXISTS `todo_app`;
USE `todo_app`;

CREATE TABLE `users` (
 `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
 `name` VARCHAR(255) NOT NULL,
 `email` VARCHAR(255) NOT NULL UNIQUE,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);

CREATE TABLE `todos` (
 `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
 `user_id` BIGINT NOT NULL COMMENT 'refer to users.id',
 `title` VARCHAR(255) NOT NULL,
 `description` TEXT NULL,
 `status` VARCHAR(255) NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);

ALTER TABLE `users`
 ADD COLUMN `password` VARCHAR(255) NOT NULL AFTER `email`;

-- dev:done
UPDATE `todos` SET `description` = '' WHERE `description` IS NULL;
ALTER TABLE `todos` MODIFY COLUMN `description` TEXT NOT NULL;

-- dev:done
CREATE TABLE `todo_histories` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `todo_id` BIGINT NOT NULL COMMENT 'refer to todos.id',
 `message` VARCHAR(255) NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);

-- dev:done
ALTER TABLE `users`
 ADD COLUMN `is_active` TINYINT NOT NULL DEFAULT 0 AFTER `email`;

-- dev:pending
CREATE TABLE `otp` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `user_id` BIGINT NOT NULL COMMENT 'refer to users.id',
 `code` VARCHAR(255) NOT NULL,
 `is_used` TINYINT NOT NULL DEFAULT 0,
 `expirate_at` TIMESTAMP NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);
