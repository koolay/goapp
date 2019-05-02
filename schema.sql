CREATE DATABASE IF NOT EXISTS `goapp`;

CREATE TABLE IF NOT EXISTS `users` (
    `id` char(20) NOT NULL, 
    `account` varchar(255)  NOT NULL,
    `password` varchar(1024)  NOT NULL,
    `display_name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL DEFAULT '',
    `mobile` varchar(64) NOT NULL DEFAULT '',
    `avatar` varchar(1024) NOT NULL DEFAULT '',
    `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    `last_login_time` TIMESTAMP NOT NULL DEFAULT NOW(),
    `disabled` TINYINT NOT NULL DEFAULT 0,
    PRIMARY KEY(`id`),
    INDEX `idx_user_mobile` (`mobile`),
    UNIQUE KEY `idx_users_account` (`account`),
    INDEX email_idx(email)
);