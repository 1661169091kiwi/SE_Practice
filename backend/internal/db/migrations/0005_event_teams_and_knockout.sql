-- 赛事参赛队伍与淘汰赛配置增强

-- 1) 扩展赛事赛制类型：新增 knockout
ALTER TABLE events
  MODIFY COLUMN format_type ENUM('points', 'knockout', 'group_knockout') NOT NULL;

-- 2) 新增赛事参赛队伍表
CREATE TABLE IF NOT EXISTS event_teams (
  event_id INT NOT NULL,
  team_id INT NOT NULL,
  group_name VARCHAR(50) DEFAULT NULL,
  slot INT DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (event_id, team_id),
  UNIQUE KEY unique_event_slot (event_id, slot),
  FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE,
  FOREIGN KEY (team_id) REFERENCES teams(team_id) ON DELETE CASCADE
);

-- 3) 允许淘汰赛对阵未绑定真实比赛（便于先搭建对阵图再生成比赛）
ALTER TABLE knockout_matches
  MODIFY COLUMN match_id INT NULL;

-- 4) 对阵显示顺序
ALTER TABLE knockout_matches
  ADD COLUMN match_order INT DEFAULT 0;
