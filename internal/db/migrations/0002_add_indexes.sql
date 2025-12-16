-- 增量迁移：添加索引与唯一约束
-- 可先执行：USE sports_management;

-- 1) subscriptions：为 student_id 添加索引
ALTER TABLE subscriptions
    ADD INDEX idx_sub_student (student_id);

-- 2) player_ratings：为 (match_id, created_at) 添加组合索引
ALTER TABLE player_ratings
    ADD INDEX idx_pr_match_created (match_id, created_at);

-- 3) match_lineups：添加唯一约束，以及 (match_id, team_id) 组合索引
ALTER TABLE match_lineups
    ADD UNIQUE KEY unique_lineup_member (match_id, team_id, student_id),
    ADD INDEX idx_lineups_match_team (match_id, team_id);

