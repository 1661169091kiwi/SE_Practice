# 体育赛事管理系统 - 数据库设计文档

## 项目概述

这是一个面向中山大学的体育赛事管理系统，支持用户注册、赛事管理、数据采集、运动员管理等功能。


## 主要业务逻辑

### 用户注册流程
1. 用户在注册界面填写信息
2. 系统检查学号是否已存在
3. 密码加密后存入`users`表

### 运动员申请流程
1. 用户在个人界面点击"申请成为运动员"
2. 填写运动员信息（号码、选择队伍）
3. 系统在`athletes`表中创建记录
4. 如果选择"创建队伍"，同时在`teams`表中创建队伍

### 赛事订阅流程
1. 用户在赛事界面点击订阅按钮
2. 系统在`subscriptions`表中添加记录
3. 订阅界面只显示已订阅的赛事

### 数据采集流程
1. 采集员在个人界面报名采集比赛
2. 比赛开始后，采集员可以录入阵容、事件等信息
3. 系统更新对应的比赛数据表

## API接口建议

### 用户相关
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `PUT /api/user/profile` - 更新用户信息
- `POST /api/athlete/apply` - 申请成为运动员

### 赛事相关
- `GET /api/events` - 获取赛事列表
- `POST /api/events/subscribe` - 订阅赛事
- `GET /api/matches` - 获取比赛列表
- `GET /api/matches/{id}` - 获取比赛详情

### 数据相关
- `POST /api/ratings` - 提交评分
- `GET /api/standings` - 获取积分榜
- `POST /api/lineups` - 录入阵容

## 部署说明

### 环境要求
- MySQL 5.7+ 或 PostgreSQL
- 字符集：UTF-8

### 初始化步骤
1. 创建数据库：`CREATE DATABASE sports_management;`
2. 按顺序执行建表语句
3. 插入基础数据（运动类型等）
4. 配置数据库连接

## 注意事项

1. **数据一致性**：运动员信息与用户信息需要保持同步
2. **权限控制**：不同角色（普通用户、运动员、采集员、管理员）的权限不同
3. **性能优化**：赛事数据和用户交互数据量较大，需要考虑分页和索引优化
4. **安全考虑**：密码需要加密存储，用户输入需要验证和转义

## 扩展建议

1. 可以考虑添加缓存层（Redis）来提升性能
2. 对于大型赛事，可以考虑分库分表
3. 添加操作日志表记录重要操作
4. 考虑数据备份和恢复策略

# 数据库执行语句
## 用户相关表

```sql
-- 用户表
CREATE TABLE users (
    student_id VARCHAR(20) PRIMARY KEY,  -- 学号
    password VARCHAR(255) NOT NULL,      -- 密码
    name VARCHAR(50) NOT NULL,           -- 姓名
    college VARCHAR(100),                -- 学院
    grade VARCHAR(20),                   -- 年级
    avatar_url VARCHAR(255),             -- 头像URL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 运动员信息表
CREATE TABLE athletes (
    athlete_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    sport_type VARCHAR(50) NOT NULL,     -- 运动类型
    team_id INT NOT NULL,                -- 所属队伍
    jersey_number VARCHAR(10),           -- 运动员号码
    is_captain BOOLEAN DEFAULT FALSE,    -- 是否为队长
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    UNIQUE KEY unique_athlete (student_id, sport_type, team_id)
);

-- 采集员信息表
CREATE TABLE collectors (
    collector_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    is_approved BOOLEAN DEFAULT TRUE,    -- 是否被批准为采集员
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);

-- 管理员表
CREATE TABLE admins (
    admin_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    role VARCHAR(50) DEFAULT 'admin',    -- 管理员角色
    permissions TEXT,                    -- 权限列表
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);
```

## 赛事相关表

```sql
-- 运动类型表
CREATE TABLE sports (
    sport_id INT AUTO_INCREMENT PRIMARY KEY,
    sport_name VARCHAR(50) NOT NULL,     -- 运动名称
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

-- 赛事表
CREATE TABLE events (
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    event_name VARCHAR(100) NOT NULL,    -- 赛事名称
    sport_id INT NOT NULL,               -- 运动类型
    season VARCHAR(20),                  -- 赛季
    round VARCHAR(50),                   -- 轮次
    format_type ENUM('points', 'group_knockout') NOT NULL, -- 赛制类型
    start_date DATE,
    end_date DATE,
    status ENUM('upcoming', 'ongoing', 'finished') DEFAULT 'upcoming',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sport_id) REFERENCES sports(sport_id)
);

-- 比赛表
CREATE TABLE matches (
    match_id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,               -- 所属赛事
    match_name VARCHAR(100),             -- 比赛名称
    round VARCHAR(50),                   -- 轮次
    match_time DATETIME NOT NULL,        -- 比赛时间
    team_a_id INT NOT NULL,              -- 队伍A
    team_b_id INT NOT NULL,              -- 队伍B
    score_team_a INT DEFAULT 0,          -- 队伍A比分
    score_team_b INT DEFAULT 0,          -- 队伍B比分
    half_score_team_a INT DEFAULT 0,     -- 半场比分A
    half_score_team_b INT DEFAULT 0,     -- 半场比分B
    status ENUM('not_started', 'ongoing', 'finished', 'cancelled') DEFAULT 'not_started',
    collector1_id INT,                   -- 采集员1
    collector2_id INT,                   -- 采集员2
    FOREIGN KEY (event_id) REFERENCES events(event_id),
    FOREIGN KEY (team_a_id) REFERENCES teams(team_id),
    FOREIGN KEY (team_b_id) REFERENCES teams(team_id),
    FOREIGN KEY (collector1_id) REFERENCES collectors(collector_id),
    FOREIGN KEY (collector2_id) REFERENCES collectors(collector_id)
);
```

## 队伍相关表

```sql
-- 队伍表
CREATE TABLE teams (
    team_id INT AUTO_INCREMENT PRIMARY KEY,
    team_name VARCHAR(100) NOT NULL,     -- 队伍名称
    sport_id INT NOT NULL,               -- 运动类型
    college VARCHAR(100),                -- 所属学院（如果是院队）
    team_type ENUM('college', 'school') NOT NULL, -- 队伍类型：院队/校队
    avatar_url VARCHAR(255),             -- 队伍头像
    description TEXT,
    created_by VARCHAR(20) NOT NULL,     -- 创建人学号
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sport_id) REFERENCES sports(sport_id),
    FOREIGN KEY (created_by) REFERENCES users(student_id)
);

-- 队伍成员表（记录队伍与运动员的关系）
CREATE TABLE team_members (
    team_member_id INT AUTO_INCREMENT PRIMARY KEY,
    team_id INT NOT NULL,
    athlete_id INT NOT NULL,
    join_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (athlete_id) REFERENCES athletes(athlete_id),
    UNIQUE KEY unique_team_member (team_id, athlete_id)
);
```

## 用户交互相关表

```sql
-- 订阅表
CREATE TABLE subscriptions (
    subscription_id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    event_id INT NOT NULL,
    subscribed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (event_id) REFERENCES events(event_id),
    UNIQUE KEY unique_subscription (student_id, event_id)
);

-- 球员评分表
CREATE TABLE player_ratings (
    rating_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,     -- 被评球员学号
    rater_id VARCHAR(20) NOT NULL,       -- 评分人学号
    score DECIMAL(3,1) NOT NULL,         -- 评分（1-10分）
    comment TEXT,                        -- 评论
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (rater_id) REFERENCES users(student_id),
    UNIQUE KEY unique_rating (match_id, student_id, rater_id)
);

-- 评论点赞表
CREATE TABLE comment_likes (
    like_id INT AUTO_INCREMENT PRIMARY KEY,
    rating_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,
    liked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rating_id) REFERENCES player_ratings(rating_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    UNIQUE KEY unique_like (rating_id, student_id)
);

-- 猜胜负表
CREATE TABLE match_predictions (
    prediction_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,
    prediction ENUM('home_win', 'draw', 'away_win') NOT NULL,
    predicted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    UNIQUE KEY unique_prediction (match_id, student_id)
);
```

## 比赛详情相关表

```sql
-- 比赛阵容表
CREATE TABLE match_lineups (
    lineup_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    team_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,     -- 球员学号
    position VARCHAR(50),                -- 位置
    is_starting BOOLEAN DEFAULT TRUE,    -- 是否首发
    jersey_number VARCHAR(10),           -- 比赛号码
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);

-- 比赛事件表（进球、红黄牌等）
CREATE TABLE match_events (
    event_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    event_type ENUM('goal', 'yellow_card', 'red_card', 'substitution', 'injury') NOT NULL,
    event_time VARCHAR(10),              -- 事件时间（如"23'"）
    student_id VARCHAR(20),              -- 相关球员
    team_id INT,                         -- 相关队伍
    description TEXT,                    -- 事件描述
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id)
);

-- 伤停信息表
CREATE TABLE injury_suspensions (
    injury_id INT AUTO_INCREMENT PRIMARY KEY,
    match_id INT NOT NULL,
    student_id VARCHAR(20) NOT NULL,     -- 球员学号
    injury_type ENUM('injury', 'suspension') NOT NULL,
    description TEXT,
    expected_return DATE,                -- 预计回归时间
    FOREIGN KEY (match_id) REFERENCES matches(match_id),
    FOREIGN KEY (student_id) REFERENCES users(student_id)
);
```

## 数据统计相关表

```sql
-- 积分榜表
CREATE TABLE standings (
    standing_id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    team_id INT NOT NULL,
    matches_played INT DEFAULT 0,        -- 已赛场次
    wins INT DEFAULT 0,                  -- 胜场
    draws INT DEFAULT 0,                 -- 平场
    losses INT DEFAULT 0,                -- 负场
    goals_for INT DEFAULT 0,             -- 进球数
    goals_against INT DEFAULT 0,         -- 失球数
    goal_difference INT DEFAULT 0,       -- 净胜球
    points INT DEFAULT 0,                -- 积分
    FOREIGN KEY (event_id) REFERENCES events(event_id),
    FOREIGN KEY (team_id) REFERENCES teams(team_id),
    UNIQUE KEY unique_standing (event_id, team_id)
);

-- 淘汰赛阶段表
CREATE TABLE knockout_stages (
    stage_id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT NOT NULL,
    stage_name VARCHAR(50) NOT NULL,     -- 阶段名称（如1/8决赛、半决赛）
    stage_order INT NOT NULL,            -- 阶段顺序
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);

-- 淘汰赛对阵表
CREATE TABLE knockout_matches (
    knockout_match_id INT AUTO_INCREMENT PRIMARY KEY,
    stage_id INT NOT NULL,
    match_id INT NOT NULL,
    prev_match_a_id INT,                 -- 前一阶段比赛A
    prev_match_b_id INT,                 -- 前一阶段比赛B
    FOREIGN KEY (stage_id) REFERENCES knockout_stages(stage_id),
    FOREIGN KEY (match_id) REFERENCES matches(match_id)
);
```

## 主要关系说明

1. **用户系统**：以`student_id`为唯一标识，支持普通用户、运动员、采集员、管理员多种角色
2. **赛事体系**：支持多种赛制（积分制、小组赛+淘汰赛），灵活的轮次管理
3. **队伍管理**：支持院队和校队，队长可以管理队伍成员
4. **数据采集**：采集员可以录入比赛阵容、事件、伤停等信息
5. **用户交互**：支持订阅、评分、评论、点赞、猜胜负等功能
6. **数据展示**：完整的积分榜和淘汰赛晋级情况展示

这个数据库结构能够很好地支持您描述的所有功能需求，并考虑了扩展性和数据一致性。
---

*此文档会根据项目进展持续更新，请团队成员及时同步最新版本*