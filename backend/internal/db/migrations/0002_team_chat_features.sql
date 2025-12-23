-- 队内聊天、投票、通知、请假功能相关表
-- 执行前请确保已执行 0001_init.sql

SET FOREIGN_KEY_CHECKS = 0;

-- 队内聊天消息表
CREATE TABLE IF NOT EXISTS team_messages (
    message_id INT AUTO_INCREMENT PRIMARY KEY,
    team_id INT NOT NULL,
    sender_id VARCHAR(20) NOT NULL,
    message_type ENUM('text', 'image', 'file', 'vote', 'notification', 'leave_request') DEFAULT 'text',
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (sender_id) REFERENCES users(student_id),
    INDEX idx_team_created (team_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 消息已读状态表
CREATE TABLE IF NOT EXISTS message_read_status (
    read_id INT AUTO_INCREMENT PRIMARY KEY,
    message_id INT NOT NULL,
    reader_id VARCHAR(20) NOT NULL,
    read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (message_id) REFERENCES team_messages(message_id) ON DELETE CASCADE,
    FOREIGN KEY (reader_id) REFERENCES users(student_id),
    UNIQUE KEY unique_read (message_id, reader_id),
    INDEX idx_reader (reader_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 投票表
CREATE TABLE IF NOT EXISTS team_votes (
    vote_id INT AUTO_INCREMENT PRIMARY KEY,
    team_id INT NOT NULL,
    creator_id VARCHAR(20) NOT NULL,
    message_id INT NOT NULL, -- 关联到team_messages表
    title VARCHAR(200) NOT NULL,
    description TEXT,
    options TEXT NOT NULL, -- JSON格式存储选项，如 ["选项1", "选项2", "选项3"]
    is_multiple BOOLEAN DEFAULT FALSE, -- 是否允许多选
    deadline DATETIME,
    status ENUM('active', 'closed') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (creator_id) REFERENCES users(student_id),
    FOREIGN KEY (message_id) REFERENCES team_messages(message_id) ON DELETE CASCADE,
    INDEX idx_team_status (team_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 投票记录表
CREATE TABLE IF NOT EXISTS vote_records (
    record_id INT AUTO_INCREMENT PRIMARY KEY,
    vote_id INT NOT NULL,
    voter_id VARCHAR(20) NOT NULL,
    selected_options TEXT NOT NULL, -- JSON格式存储选中的选项索引，如 [0, 2]
    voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vote_id) REFERENCES team_votes(vote_id) ON DELETE CASCADE,
    FOREIGN KEY (voter_id) REFERENCES users(student_id),
    UNIQUE KEY unique_vote (vote_id, voter_id),
    INDEX idx_voter (voter_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 通知表
CREATE TABLE IF NOT EXISTS team_notifications (
    notification_id INT AUTO_INCREMENT PRIMARY KEY,
    team_id INT NOT NULL,
    sender_id VARCHAR(20) NOT NULL,
    message_id INT NOT NULL, -- 关联到team_messages表
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    notification_type ENUM('info', 'warning', 'urgent') DEFAULT 'info',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (sender_id) REFERENCES users(student_id),
    FOREIGN KEY (message_id) REFERENCES team_messages(message_id) ON DELETE CASCADE,
    INDEX idx_team_created (team_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 请假申请表
CREATE TABLE IF NOT EXISTS leave_requests (
    leave_id INT AUTO_INCREMENT PRIMARY KEY,
    team_id INT NOT NULL,
    applicant_id VARCHAR(20) NOT NULL,
    message_id INT, -- 关联到team_messages表（可选）
    leave_type ENUM('sick', 'personal', 'other') DEFAULT 'personal',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NOT NULL,
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    reviewer_id VARCHAR(20), -- 审核人（队长）
    review_comment TEXT,
    reviewed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (applicant_id) REFERENCES users(student_id),
    FOREIGN KEY (reviewer_id) REFERENCES users(student_id),
    FOREIGN KEY (message_id) REFERENCES team_messages(message_id) ON DELETE SET NULL,
    INDEX idx_team_status (team_id, status),
    INDEX idx_applicant (applicant_id),
    INDEX idx_reviewer (reviewer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;

