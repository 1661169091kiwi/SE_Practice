-- Seed data for testing collector workflow

-- 1. Insert Sports
INSERT INTO sports (sport_name, description) VALUES 
('Football', 'Standard 11-a-side Football'),
('Basketball', 'Standard Basketball')
ON DUPLICATE KEY UPDATE description=VALUES(description);

-- Variable to store Football ID
SET @football_id = (SELECT sport_id FROM sports WHERE sport_name = 'Football' LIMIT 1);

-- 2. Insert Users
-- 2023000001: Collector
-- 2023000002: Student/Team Creator
-- Password is '123456' (Plain text for dev/test environment as requested)
INSERT INTO users (student_id, password, name, college, grade) VALUES
('2023000001', '123456', 'Collector User', 'Computer Science', '2023'),
('2023000002', '123456', 'Student User', 'Software Engineering', '2023')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 3. Insert Collectors
-- Ensure 2023000001 is a collector
INSERT INTO collectors (student_id, is_approved) 
SELECT '2023000001', TRUE 
WHERE NOT EXISTS (SELECT 1 FROM collectors WHERE student_id = '2023000001');

-- 4. Insert Teams
-- Teams created by 2023000002
INSERT INTO teams (team_name, sport_id, college, team_type, created_by) 
SELECT 'Team A', @football_id, 'Computer Science', 'college', '2023000002' 
WHERE NOT EXISTS (SELECT 1 FROM teams WHERE team_name = 'Team A');

INSERT INTO teams (team_name, sport_id, college, team_type, created_by) 
SELECT 'Team B', @football_id, 'Software Engineering', 'college', '2023000002' 
WHERE NOT EXISTS (SELECT 1 FROM teams WHERE team_name = 'Team B');

-- Variables for Team IDs
SET @team_a_id = (SELECT team_id FROM teams WHERE team_name = 'Team A' LIMIT 1);
SET @team_b_id = (SELECT team_id FROM teams WHERE team_name = 'Team B' LIMIT 1);

-- 5. Insert Events
INSERT INTO events (event_name, sport_id, format_type, start_date, status) 
SELECT 'Campus Cup 2024', @football_id, 'points', CURDATE(), 'ongoing'
WHERE NOT EXISTS (SELECT 1 FROM events WHERE event_name = 'Campus Cup 2024');

-- Variable for Event ID
SET @event_id = (SELECT event_id FROM events WHERE event_name = 'Campus Cup 2024' LIMIT 1);

-- 6. Insert Matches
-- Match between Team A and Team B
INSERT INTO matches (event_id, match_name, match_time, team_a_id, team_b_id, status)
SELECT @event_id, 'Group Stage Match 1', DATE_ADD(NOW(), INTERVAL 1 DAY), @team_a_id, @team_b_id, 'not_started'
WHERE NOT EXISTS (SELECT 1 FROM matches WHERE event_id = @event_id AND team_a_id = @team_a_id AND team_b_id = @team_b_id);
