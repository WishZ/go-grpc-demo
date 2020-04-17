table_name:m_todo
```sql
CREATE TABLE `m_todo` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `title` varchar(200) DEFAULT NULL,
    `description` varchar(1024) DEFAULT NULL,
    `reminder` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created_at',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated_at',
    `deleted_at` datetime NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT 'deleted_at',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB CHARSET=utf8mb4 COMMENT "待办事项";
```