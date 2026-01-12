package db

import (
	"log"
)

// InitSchema 初始化数据库表结构
func InitSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS match_comments (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			match_id BIGINT NOT NULL,
			student_id VARCHAR(64) NOT NULL,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_match_id (match_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS match_stats (
			match_id BIGINT PRIMARY KEY,
			home_possession INT DEFAULT 50,
			away_possession INT DEFAULT 50,
			home_shots INT DEFAULT 0,
			away_shots INT DEFAULT 0,
			home_shots_on_target INT DEFAULT 0,
			away_shots_on_target INT DEFAULT 0,
			home_fouls INT DEFAULT 0,
			away_fouls INT DEFAULT 0,
			home_corners INT DEFAULT 0,
			away_corners INT DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS match_lineups (
			lineup_id BIGINT AUTO_INCREMENT PRIMARY KEY,
			match_id BIGINT NOT NULL,
			team_id BIGINT NOT NULL,
			student_id VARCHAR(64) NOT NULL,
			position VARCHAR(32),
			is_starting BOOLEAN DEFAULT FALSE,
			jersey_number VARCHAR(8),
			INDEX idx_match_team (match_id, team_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS match_lineups_external (
			ext_id BIGINT AUTO_INCREMENT PRIMARY KEY,
			match_id BIGINT NOT NULL,
			team_id BIGINT NOT NULL,
			player_name VARCHAR(50) NOT NULL,
			position VARCHAR(50),
			is_starting BOOLEAN DEFAULT TRUE,
			jersey_number VARCHAR(10),
			INDEX idx_ext_match_team (match_id, team_id),
			FOREIGN KEY (match_id) REFERENCES matches(match_id),
			FOREIGN KEY (team_id) REFERENCES teams(team_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS event_teams (
			event_id INT NOT NULL,
			team_id INT NOT NULL,
			group_name VARCHAR(50) DEFAULT NULL,
			slot INT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (event_id, team_id),
			UNIQUE KEY unique_event_slot (event_id, slot),
			FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE,
			FOREIGN KEY (team_id) REFERENCES teams(team_id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS knockout_stages (
			stage_id INT AUTO_INCREMENT PRIMARY KEY,
			event_id INT NOT NULL,
			stage_name VARCHAR(50) NOT NULL,
			stage_order INT NOT NULL,
			FOREIGN KEY (event_id) REFERENCES events(event_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS knockout_matches (
			knockout_match_id INT AUTO_INCREMENT PRIMARY KEY,
			stage_id INT NOT NULL,
			match_id INT NULL,
			match_order INT DEFAULT 0,
			prev_match_a_id INT,
			prev_match_b_id INT,
			FOREIGN KEY (stage_id) REFERENCES knockout_stages(stage_id) ON DELETE CASCADE,
			FOREIGN KEY (match_id) REFERENCES matches(match_id) ON DELETE SET NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS carousel_images (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			image_url VARCHAR(255) NOT NULL,
			title VARCHAR(100),
			link_url VARCHAR(255),
			sort_order INT DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}

	for _, query := range queries {
		if _, err := Exec(query); err != nil {
			log.Printf("Failed to execute query: %s, error: %v", query, err)
			return err
		}
	}

	// 动态添加 match_order 到 knockout_matches 表
	checkKnockoutMatchOrderQuery := `SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'knockout_matches' AND COLUMN_NAME = 'match_order'`
	var kmCount int
	if err := QueryRow(checkKnockoutMatchOrderQuery).Scan(&kmCount); err != nil {
		log.Printf("Failed to check column match_order in knockout_matches: %v", err)
	} else if kmCount == 0 {
		log.Println("Adding match_order column to knockout_matches table...")
		if _, err := Exec("ALTER TABLE knockout_matches ADD COLUMN match_order INT DEFAULT 0"); err != nil {
			return err
		}
	}

	// 动态添加 scheduled_time 到 knockout_matches 表
	checkKnockoutScheduledTimeQuery := `SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'knockout_matches' AND COLUMN_NAME = 'scheduled_time'`
	var kmtCount int
	if err := QueryRow(checkKnockoutScheduledTimeQuery).Scan(&kmtCount); err != nil {
		log.Printf("Failed to check column scheduled_time in knockout_matches: %v", err)
	} else if kmtCount == 0 {
		log.Println("Adding scheduled_time column to knockout_matches table...")
		if _, err := Exec("ALTER TABLE knockout_matches ADD COLUMN scheduled_time DATETIME NULL"); err != nil {
			return err
		}
	}

	// 确保 knockout_matches.match_id 允许为 NULL
	modifyKnockoutMatchIdQuery := "ALTER TABLE knockout_matches MODIFY COLUMN match_id INT NULL"
	if _, err := Exec(modifyKnockoutMatchIdQuery); err != nil {
		log.Printf("Note: attempt to modify match_id nullable returned: %v", err)
	}

	// 尝试更新 events 表的 format_type 枚举以支持 knockout
	alterEnumQuery := "ALTER TABLE events MODIFY COLUMN format_type ENUM('points', 'knockout', 'group_knockout') NOT NULL"
	if _, err := Exec(alterEnumQuery); err != nil {
		// 可能是已经存在或者其他非致命错误，记录日志但不中断启动（除非真的很严重，但通常 ENUM 修改失败意味着无需修改或不兼容）
		log.Printf("Note: attempt to update format_type enum returned: %v", err)
	}

	// 动态添加 match_id 到 subscriptions 表
	// 检查列是否存在
	checkColQuery := `SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscriptions' AND COLUMN_NAME = 'match_id'`
	var count int
	if err := QueryRow(checkColQuery).Scan(&count); err != nil {
		log.Printf("Failed to check column match_id: %v", err)
	} else if count == 0 {
		// 添加列
		log.Println("Adding match_id column to subscriptions table...")
		if _, err := Exec("ALTER TABLE subscriptions ADD COLUMN match_id INT DEFAULT NULL"); err != nil {
			return err
		}
		// 添加外键
		if _, err := Exec("ALTER TABLE subscriptions ADD CONSTRAINT fk_subscriptions_match FOREIGN KEY (match_id) REFERENCES matches(match_id)"); err != nil {
			log.Printf("Warning: failed to add foreign key: %v", err)
		}
	}

	// 修改唯一索引 (无论 match_id 是否是新加的，都要确保索引包含 match_id)
	// 先尝试删除旧索引（如果存在）
	// 注意：如果旧索引不存在，DROP INDEX 会报错，所以这里只作为尝试，忽略错误
	Exec("ALTER TABLE subscriptions DROP INDEX unique_subscription")

	// 检查新索引是否已存在，如果不存在则添加
	checkIndexQuery := `SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'subscriptions' AND INDEX_NAME = 'unique_subscription_entry'`
	var indexCount int
	if err := QueryRow(checkIndexQuery).Scan(&indexCount); err != nil {
		log.Printf("Failed to check index unique_subscription_entry: %v", err)
	} else if indexCount == 0 {
		// 添加新索引
		if _, err := Exec("ALTER TABLE subscriptions ADD UNIQUE KEY unique_subscription_entry (student_id, event_id, match_id)"); err != nil {
			log.Printf("Warning: failed to add unique key: %v", err)
		}
	}

	// 动态添加 is_approved 到 teams 表
	checkTeamColQuery := `SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'teams' AND COLUMN_NAME = 'is_approved'`
	var teamCount int
	if err := QueryRow(checkTeamColQuery).Scan(&teamCount); err != nil {
		log.Printf("Failed to check column is_approved in teams: %v", err)
	} else if teamCount == 0 {
		log.Println("Adding is_approved column to teams table...")
		if _, err := Exec("ALTER TABLE teams ADD COLUMN is_approved BOOLEAN DEFAULT FALSE"); err != nil {
			return err
		}
	}

	// 动态添加 is_approved 到 team_members 表
	checkTeamMemberColQuery := `SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'team_members' AND COLUMN_NAME = 'is_approved'`
	var teamMemberCount int
	if err := QueryRow(checkTeamMemberColQuery).Scan(&teamMemberCount); err != nil {
		log.Printf("Failed to check column is_approved in team_members: %v", err)
	} else if teamMemberCount == 0 {
		log.Println("Adding is_approved column to team_members table...")
		if _, err := Exec("ALTER TABLE team_members ADD COLUMN is_approved BOOLEAN DEFAULT FALSE"); err != nil {
			return err
		}
	}

	log.Println("Database schema initialized successfully")
	return nil
}
