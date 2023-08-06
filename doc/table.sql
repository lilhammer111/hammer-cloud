-- create user table
CREATE TABLE `tbl_user`(
  `id` int NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '',
  `user_pwd` varchar(256) NOT NULL DEFAULT '',
  `email` varchar(64) DEFAULT '',
  `phone` varchar(128) DEFAULT '',
  `email_validated` tinyint(1) DEFAULT '0',
  `phone_validated` tinyint(1) DEFAULT '0',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `profile` text,
  `status` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- create user token table
CREATE TABLE `tbl_user_token` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '',
  `user_token` char(40) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- create unique file table
CREATE TABLE `tbl_file` (
    `id` int NOT NULL AUTO_INCREMENT,
    `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT 'file hash',
    `file_name` varchar(256) NOT NULL DEFAULT '',
    `file_size` bigint DEFAULT '0',
    `file_addr` varchar(1024) NOT NULL DEFAULT '',
    `create_at` datetime DEFAULT CURRENT_TIMESTAMP,
     `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `status` int NOT NULL DEFAULT '0',
    `ext1` int DEFAULT '0',
    `ext2` text,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_hash` (`file_sha1`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- create user file table
create table tbl_user_file (
    id int(11) not null key auto_increment,
    user_name varchar(64) not null,
    file_sha1 varchar(64) not null default '',
    file_size bigint(20) default '0',
    file_name varchar(256) NOT NULL DEFAULT '',
    upload_at datetime default current_timestamp,
    last_update datetime default current_timestamp on update current_timestamp,
    status int(11) not null default '0',
    unique key idx_user_file (user_name, file_sha1),
    key idx_status (status),
    key idx_user_id (user_name)
) engine innodb default charset utf8mb4;