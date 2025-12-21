-- 添加 rank 字段到 standings 表
ALTER TABLE standings ADD COLUMN `rank` INT DEFAULT 0 COMMENT '排名';
ALTER TABLE standings ADD INDEX idx_event_rank (event_id, `rank`);

-- 比赛表索引优化
-- 注意：如果索引名重复可能会报错
CREATE INDEX idx_matches_event_status ON matches (event_id, status);
CREATE INDEX idx_matches_team_a ON matches (team_a_id);
CREATE INDEX idx_matches_team_b ON matches (team_b_id);
CREATE INDEX idx_matches_time ON matches (match_time);
