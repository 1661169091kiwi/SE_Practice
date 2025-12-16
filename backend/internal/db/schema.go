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
	}

	for _, query := range queries {
		if _, err := Exec(query); err != nil {
			log.Printf("Failed to execute query: %s, error: %v", query, err)
			return err
		}
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
		// 修改唯一索引
		// 先尝试删除旧索引
		Exec("ALTER TABLE subscriptions DROP INDEX unique_subscription")
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
