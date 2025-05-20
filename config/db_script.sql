CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL COMMENT 'Store only bcrypt hashed passwords',
    is_privileged ENUM('yes', 'no') NOT NULL DEFAULT 'no',
    privilege_tier ENUM('basic', 'premium', 'admin') NULL COMMENT 'Only populated if is_privileged=yes',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_privilege (is_privileged, privilege_tier),
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;