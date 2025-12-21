ALTER TABLE subscriptions
  ADD INDEX idx_sub_student (student_id);

ALTER TABLE player_ratings
  ADD INDEX idx_pr_match_created (match_id, created_at);

ALTER TABLE match_lineups
  ADD UNIQUE KEY unique_lineup_member (match_id, team_id, student_id),
  ADD INDEX idx_lineups_match_team (match_id, team_id);
