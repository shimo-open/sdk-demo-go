-- Create database
CREATE DATABASE IF NOT EXISTS `sdk_demo_go` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;

-- Use database
use sdk_demo_go;

-- Create table structure
DROP TABLE IF EXISTS `app_clients`;
CREATE TABLE `app_clients`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `app_id`     varchar(255) NOT NULL DEFAULT '' COMMENT 'appId',
    `app_secret` varchar(255) NOT NULL DEFAULT '' COMMENT 'appSecret',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE,
    KEY          `idx_app_id` (`app_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='App clients table';

DROP TABLE IF EXISTS `departments`;
CREATE TABLE `departments`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `name`       varchar(100) NOT NULL DEFAULT '' COMMENT 'Department name',
    `parent_id`  int(11) NOT NULL DEFAULT 0 COMMENT 'Parent ID',
    `team_id`    int(11) NOT NULL DEFAULT 0 COMMENT 'Team ID',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE,
    KEY          `idx_team_id` (`team_id`) USING BTREE,
    KEY          `idx_parent_id` (`parent_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='Departments table';

DROP TABLE IF EXISTS `dept_members`;
CREATE TABLE `dept_members`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `dept_id`    int(11) NOT NULL DEFAULT 0 COMMENT 'Department ID',
    `user_id`    int(11) NOT NULL DEFAULT 0 COMMENT 'Member ID',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uniq_dept_id_user_id` (`dept_id`,`user_id`) USING BTREE,
    KEY          `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='Department members table';

DROP TABLE IF EXISTS `events`;
CREATE TABLE `events`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `type`       varchar(30)   NOT NULL DEFAULT '' COMMENT 'Event type',
    `file_id`    varchar(128)  NOT NULL DEFAULT '' COMMENT 'File ID',
    `user_id`    varchar(128)  NOT NULL DEFAULT '' COMMENT 'Related user ID',
    `raw_data`   varchar(10000) NOT NULL DEFAULT '' COMMENT 'Message content',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    `headers`    varchar(5000) NOT NULL DEFAULT '' COMMENT 'Event headers',
    PRIMARY KEY (`id`) USING BTREE,
    KEY          `idx_file_id` (`file_id`) USING BTREE,
    KEY          `idx_user_id` (`user_id`) USING BTREE,
    KEY          `idx_type` (`type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='Events table';

DROP TABLE IF EXISTS `file_permissions`;
CREATE TABLE `file_permissions`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `file_id`     bigint(20) NOT NULL DEFAULT 0 COMMENT 'File ID',
    `user_id`     bigint(20) NOT NULL DEFAULT 0 COMMENT 'User ID',
    `role`        enum('owner','collaborator') NOT NULL DEFAULT 'collaborator' COMMENT 'Role (owner/collaborator)',
    `permissions` json COMMENT 'Permissions',
    `created_at`  bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at`  bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at`  bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uniq_file_id_user_id` (`file_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='File permissions table';

DROP TABLE IF EXISTS `files`;
CREATE TABLE `files`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `guid`          varchar(64) NOT NULL DEFAULT '' COMMENT 'File GUID (unique identifier)',
    `name`          varchar(255) NOT NULL DEFAULT '' COMMENT 'File name',
    `type`          varchar(255) NOT NULL DEFAULT '' COMMENT 'File type',
    `file_path`     varchar(255) NOT NULL DEFAULT '' COMMENT 'File path',
    `creator_id`    bigint(20) NOT NULL DEFAULT 0 COMMENT 'Creator ID',
    `is_shimo_file` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Is Shimo file',
    `shimo_type`    varchar(255) NOT NULL DEFAULT '' COMMENT 'Shimo file type',
    `created_at`    bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at`    bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at`    bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uniq_guid` (`guid`) USING BTREE,
    KEY             `files_id_creator_id_index` (`guid`,`creator_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='Files table';

DROP TABLE IF EXISTS `team_role`;
CREATE TABLE `team_role`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `team_id`    bigint(20) NOT NULL DEFAULT 0 COMMENT 'Team ID',
    `user_id`    bigint(20) NOT NULL DEFAULT 0 COMMENT 'User ID',
    `role`       enum('creator','manager','member') NOT NULL DEFAULT 'member' COMMENT 'Role (creator/manager/member)',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uniq_team_id_user_id` (`team_id`,`user_id`) USING BTREE,
    KEY          `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='Team members table';

DROP TABLE IF EXISTS `teams`;
CREATE TABLE `teams`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `name`       varchar(100) NOT NULL DEFAULT '' COMMENT 'Team name',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='Teams table';

DROP TABLE IF EXISTS `test_api`;
CREATE TABLE `test_api`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `test_id`    varchar(64) NOT NULL DEFAULT '' COMMENT 'Test UUID',
    `test_type`  varchar(64) NOT NULL DEFAULT '' COMMENT 'Test type',
    `api_name`   varchar(255) NOT NULL DEFAULT '' COMMENT 'API name',
    `success`    tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Success (0-false;1-true)',
    `http_code`  int(11) NOT NULL DEFAULT 0 COMMENT 'Status code',
    `http_resp`  text COMMENT 'Response result',
    `err_msg`    text COMMENT 'Error message',
    `path_str`   text COMMENT 'API request path',
    `body_req`   text COMMENT 'Body parameters',
    `query`      text NOT NULL DEFAULT '' COMMENT 'Query parameters',
    `form_data`  varchar(255) NOT NULL DEFAULT '' COMMENT 'Form data parameters',
    `file_ext`   varchar(64)  NOT NULL DEFAULT '' COMMENT 'File extension/Export file type',
    `time_consuming` varchar(64) NOT NULL DEFAULT '' COMMENT 'Time consuming',
    `start_time` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Test start time',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='System test table';

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `name`       varchar(255) NOT NULL DEFAULT '' COMMENT 'User name',
    `email`      varchar(255) NOT NULL DEFAULT '' COMMENT 'Email address',
    `avatar`     varchar(255) NOT NULL DEFAULT '' COMMENT 'Avatar URL',
    `password`   varchar(255) NOT NULL DEFAULT '' COMMENT 'Password',
    `app_id`     varchar(255) NOT NULL DEFAULT '' COMMENT 'appId',
    `created_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Created timestamp',
    `updated_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Updated timestamp',
    `deleted_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'Deleted timestamp',
    PRIMARY KEY (`id`) USING BTREE,
    KEY          `idx_email` (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='Users table';

