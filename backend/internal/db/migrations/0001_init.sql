-- 初始数据库结构。表定义来自 database_struction.md
-- 可先执行：CREATE DATABASE IF NOT EXISTS sports_management DEFAULT CHARACTER SET utf8mb4;
-- 然后：USE sports_management;

SET FOREIGN_KEY_CHECKS = 0;

-- 运动类型表
CREATE TABLE IF NOT EXISTS sports (
    sport_id INT AUTO_INCREMENT PRIMARY KEY,
    sport_name VARCHAR(50) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    student_id VARCHAR(20) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    college VARCHAR(100),
    grade VARCHAR(20),
    avatar_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 队伍表
CREATE TABLE IF NOT EXISTS teams (
    team_id INT AUTO_INCREMENT PRIMARY KEY,
    team_name VARCHAR(100) NOT NULL,
    sport_id INT NOT NULL,
    college VARCHAR(100),
    team_type ENUM('college', 'school') NOT NULL,
    avatar_url VARCHAR(255),
    description TEXT,
    created_by VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sport_id) REFERENCES sports(sport_id),
    FOREIGN KEY (created_by) REFERENCES users(student_id)
);

-- 运动员信息表
CREATE TABLE IF NOT EXISTS athletes (
    athlete_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    sport_type VARCHAR(50) NOT NULL,
    team_id INT NOT NULL,
    jersey_number VARCHAR(10),
    is_captain BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    UNIQUE KEY unique_athlete (student_id, sport_type, team_id)
);

-- 队伍成员表
CREATE TABLE IF NOT EXISTS team_members (
    team_member_id INT AUTO_INCREMENT PRIMARY KEY,
    team_id INT NOT NULL,
    athlete_id INT NOT NULL,
    join_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (athlete_id) REFERENCES athletes(athlete_id),
    UNIQUE KEY unique_team_member (team_id, athlete_id)
);

-- 采集员信息表
CREATE TABLE IF NOT EXISTS collectors (
    collector_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    is_approved BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);

-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    admin_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    role VARCHAR(50) DEFAULT 'admin',
    permissions TEXT,
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);

-- 赛事表
CREATE TABLE IF NOT EXISTS events (
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    event_name VARCHAR(100) NOT NULL,
    sport_id INT NOT NULL,
    season VARCHAR(20),
    round VARCHAR(50),
    format_type ENUM('points', 'group_knockout') NOT NULL,
    start_date DATE,
    end_date DATE,
    status ENUM('upcoming', 'ongoing', 'finished') DEFAULT 'upcoming',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sport_id) REFERENCES sports(sport_id)
);

-- 比赛表
CREATE TABLE IF NOT EXISTS matches (
    match_id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    match_name VARCHAR(100),
    round VARCHAR(50),
    match_time DATETIME NOT NULL,
    team_a_id INT NOT NULL,
    team_b_id INT NOT NULL,
    score_team_a INT DEFAULT 0,
    score_team_b INT DEFAULT 0,
    half_score_team_a INT DEFAULT 0,
    half_score_team_b INT DEFAULT 0,
    status ENUM('not_started', 'ongoing', 'finished', 'cancelled') DEFAULT 'not_started',
    collector1_id INT,
    collector2_id INT,
    FOREIGN KEY (event_id) REFERENCES events(event_id),
    FOREIGN KEY (team_a_id) REFERENCES teams(team_id),
    FOREIGN KEY (team_b_id) REFERENCES teams(team_id),
    FOREIGN KEY (collector1_id) REFERENCES collectors(collector_id),
    FOREIGN KEY (collector2_id) REFERENCES collectors(collector_id)
);

-- 订阅表
CREATE TABLE IF NOT EXISTS subscriptions (
    subscription_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    event_id INT NOT NULL,
    subscribed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (event_id) REFERENCES events(event_id),
    UNIQUE KEY unique_subscription (student_id, event_id)
);

-- 球员评分表
CREATE TABLE IF NOT EXISTS player_ratings (
    rating_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,
    rater_id VARCHAR(20) NOT NULL,
    score DECIMAL(3,1) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (rater_id) REFERENCES users(student_id),
    UNIQUE KEY unique_rating (match_id, student_id, rater_id)
);

-- 评论点赞表
CREATE TABLE IF NOT EXISTS comment_likes (
    like_id INT AUTO_INCREMENT PRIMARY KEY,
    rating_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,
    liked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rating_id) REFERENCES player_ratings(rating_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    UNIQUE KEY unique_like (rating_id, student_id)
);

-- 比赛阵容表
CREATE TABLE IF NOT EXISTS match_lineups (
    lineup_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    team_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,
    position VARCHAR(50),
    is_starting BOOLEAN DEFAULT TRUE,
    jersey_number VARCHAR(10),
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);

-- 比赛事件表
CREATE TABLE IF NOT EXISTS match_events (
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    event_type ENUM('goal', 'yellow_card', 'red_card', 'substitution', 'injury') NOT NULL,
    event_time VARCHAR(10),
    student_id VARCHAR(20),
    team_id INT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id)
);

-- 伤停信息表
CREATE TABLE IF NOT EXISTS injury_suspensions (
    injury_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,
    injury_type ENUM('injury', 'suspension') NOT NULL,
    description TEXT,
    expected_return DATE,
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);

-- 积分榜表
CREATE TABLE IF NOT EXISTS standings (
    standing_id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    team_id INT NOT NULL,
    matches_played INT DEFAULT 0,
    wins INT DEFAULT 0,
    draws INT DEFAULT 0,
    losses INT DEFAULT 0,
    goals_for INT DEFAULT 0,
    goals_against INT DEFAULT 0,
    goal_difference INT DEFAULT 0,
    points INT DEFAULT 0,
    FOREIGN KEY (event_id) REFERENCES events(event_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    UNIQUE KEY unique_standing (event_id, team_id)
);

-- 淘汰赛阶段表
CREATE TABLE IF NOT EXISTS knockout_stages (
    stage_id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    stage_name VARCHAR(50) NOT NULL,
    stage_order INT NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);

-- 淘汰赛对阵表
CREATE TABLE IF NOT EXISTS knockout_matches (
    knockout_match_id INT AUTO_INCREMENT PRIMARY KEY,
    stage_id INT NOT NULL,
    match_id INT NOT NULL,
    prev_match_a_id INT,
    prev_match_b_id INT,
    FOREIGN KEY (stage_id) REFERENCES knockout_stages(stage_id),
    FOREIGN KEY (match_id) REFERENCES matches(match_id)
);

SET FOREIGN_KEY_CHECKS = 1;


