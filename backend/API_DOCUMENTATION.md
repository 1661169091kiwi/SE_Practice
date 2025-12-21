# 体育赛事管理系统 - API接口清单

本文档详细列出了体育赛事管理系统后端已实现的所有API接口，包括接口路径、请求方法、功能描述、请求参数、响应格式等信息。

## 1. 健康检查接口

### 1.1 获取系统健康状态
- **路径**: `/api/health`
- **方法**: `GET`
- **功能**: 检查系统是否正常运行
- **请求参数**: 无
- **响应格式**:
  ```json
  {
    "code": 200,
    "data": {
      "status": "ok"
    },
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/user.go`

## 2. 认证相关接口

### 2.1 用户注册
- **路径**: `/api/register`
- **方法**: `POST`
- **功能**: 注册新用户
- **请求参数** (JSON):
  ```json
  {
    "student_id": "2023000001",
    "password": "password123",
    "name": "张三"
  }
  ```
- **响应格式**:
  - 成功:
    ```json
    {
      "code": 200,
      "data": {
        "student_id": "2023000001",
        "name": "张三",
        "avatar": "",
        "college": ""
      },
      "message": "success"
    }
    ```
  - 失败 (用户已存在):
    ```json
    {
      "code": 409,
      "message": "user already exists"
    }
    ```
  - 失败 (参数错误):
    ```json
    {
      "code": 400,
      "message": "student_id, password, name are required"
    }
    ```
- **实现文件**: `backend/internal/handler/auth.go`

### 2.2 用户登录
- **路径**: `/api/login`
- **方法**: `POST`
- **功能**: 用户登录认证
- **请求参数** (JSON):
  ```json
  {
    "student_id": "2023000001",
    "password": "password123"
  }
  ```
- **响应格式**:
  - 成功:
    ```json
    {
      "code": 200,
      "data": {
        "token": "stub.jwt.token",
        "user": {
          "student_id": "2023000001",
          "name": "张三",
          "role": "student",
          "avatar": "",
          "college": ""
        }
      },
      "message": "success"
    }
    ```
  - 失败 (认证失败):
    ```json
    {
      "code": 401,
      "message": "invalid student_id or password"
    }
    ```
- **实现文件**: `backend/internal/handler/auth.go`

## 3. 用户资料相关接口

### 3.1 获取用户个人资料
- **路径**: `/api/user/profile/{student_id}`
- **方法**: `GET`
- **功能**: 根据学号获取用户详细信息
- **请求参数**:
  - 路径参数: `student_id` (学号)
- **响应格式**:
  - 成功:
    ```json
    {
      "code": 200,
      "data": {
        "student_id": "2023000001",
        "name": "张三",
        "avatar": "/uploads/avatars/2023000001.jpg",
        "college": "计算机科学与技术学院"
      },
      "message": "success"
    }
    ```
  - 失败 (用户不存在):
    ```json
    {
      "code": 404,
      "message": "user not found"
    }
    ```
- **实现文件**: `backend/internal/handler/auth.go`

### 3.2 更新用户学院信息
- **路径**: `/api/user/college`
- **方法**: `POST`
- **功能**: 更新用户的学院信息
- **请求参数** (JSON):
  ```json
  {
    "student_id": "2023000001",
    "college": "计算机科学与技术学院"
  }
  ```
- **响应格式**:
  - 成功:
    ```json
    {
      "code": 200,
      "data": {
        "student_id": "2023000001",
        "college": "计算机科学与技术学院",
        "message": "学院信息更新成功"
      },
      "message": "success"
    }
    ```
  - 失败 (用户不存在):
    ```json
    {
      "code": 404,
      "message": "user not found"
    }
    ```
  - 失败 (参数错误):
    ```json
    {
      "code": 400,
      "message": "student_id and college are required"
    }
    ```
- **实现文件**: `backend/internal/handler/auth.go`

### 3.3 更新用户头像
- **路径**: `/api/user/avatar/{student_id}`
- **方法**: `POST`
- **功能**: 上传并更新用户头像
- **请求参数**:
  - 路径参数: `student_id` (学号)
  - Form-data: `avatar` (图片文件，支持jpg、jpeg、png、gif格式，大小不超过10MB)
- **响应格式**:
  - 成功:
    ```json
    {
      "code": 200,
      "data": {
        "avatar_url": "/uploads/avatars/2023000001.jpg",
        "message": "头像更新成功"
      },
      "message": "success"
    }
    ```
  - 失败 (用户不存在):
    ```json
    {
      "code": 404,
      "message": "user not found"
    }
    ```
  - 失败 (文件错误):
    ```json
    {
      "code": 400,
      "message": "file too large"
    }
    ```
- **实现文件**: `backend/internal/handler/auth.go`

### 3.4 修改用户密码
- **路径**: `/api/user/change-password`
- **方法**: `POST`
- **功能**: 修改用户密码
- **请求参数** (JSON):
  ```json
  {
    "student_id": "2023000001",
    "old_password": "password123",
    "new_password": "newpassword456"
  }
  ```
- **响应格式**:
  - 成功:
    ```json
    {
      "code": 200,
      "data": {
        "message": "password changed successfully"
      },
      "message": "success"
    }
    ```
  - 失败 (认证失败):
    ```json
    {
      "code": 401,
      "message": "invalid credentials"
    }
    ```
  - 失败 (新旧密码相同):
    ```json
    {
      "code": 400,
      "message": "new password cannot be the same as old password"
    }
    ```
- **实现文件**: `backend/internal/handler/auth.go`

## 4. 赛事相关接口（基础实现）

### 4.1 获取赛事列表
- **路径**: `/api/events`
- **方法**: `GET`
- **功能**: 获取所有赛事列表
- **请求参数**: 无
- **响应格式**:
  ```json
  {
    "code": 200,
    "data": [],
    "message": "success"
  }
  ```
- **说明**: 目前仅返回空数组，需要进一步实现
- **实现文件**: `backend/internal/handler/event.go`

### 4.2 订阅赛事
- **路径**: `/api/events/subscribe`
- **方法**: `POST`
- **功能**: 订阅指定赛事
- **请求参数**: 未明确指定（需进一步实现）
- **响应格式**:
  ```json
  {
    "code": 200,
    "data": {
      "subscribed": true
    },
    "message": "success"
  }
  ```
- **说明**: 目前仅返回订阅成功状态，需要进一步实现参数验证和业务逻辑
- **实现文件**: `backend/internal/handler/event.go`

## 5. 比赛相关接口（基础实现）

### 5.1 获取比赛列表
- **路径**: `/api/matches`
- **方法**: `GET`
- **功能**: 获取所有比赛列表
- **请求参数**: 无
- **响应格式**:
  ```json
  {
    "code": 200,
    "data": [],
    "message": "success"
  }
  ```
- **说明**: 目前仅返回空数组，需要进一步实现
- **实现文件**: `backend/internal/handler/match.go`

### 5.2 获取比赛详情
- **路径**: `/api/matches/{id}`
- **方法**: `GET`
- **功能**: 根据ID获取比赛详细信息
- **请求参数**:
  - 路径参数: `id` (比赛ID)
- **响应格式**:
  ```json
  {
    "code": 200,
    "data": {
      "id": "match123"
    },
    "message": "success"
  }
  ```
- **说明**: 目前仅返回比赛ID，需要进一步实现完整的比赛信息
- **实现文件**: `backend/internal/handler/match.go`

## 接口实现状态说明

| 接口类型 | 实现状态 | 说明 |
|---------|---------|------|
| 认证接口 | 已实现 | 包含注册、登录功能 |
| 用户资料接口 | 已实现 | 包含获取资料、更新学院、更新头像、修改密码功能 |
| 健康检查 | 已实现 | 系统状态检查功能 |
| 赛事接口 | 基础实现 | 仅实现接口框架，需要完善业务逻辑 |
| 比赛接口 | 基础实现 | 仅实现接口框架，需要完善业务逻辑 |

## 注意事项

1. 所有接口统一使用`util.OK`和`util.Error`返回标准格式响应
2. 认证相关接口目前使用占位JWT token，实际部署时需要实现真实的JWT生成和验证
3. 文件上传功能目前使用模拟URL，实际部署时需要实现真实的文件存储逻辑
4. 赛事和比赛相关接口需要进一步完善业务逻辑和数据模型
5. 成功响应统一约定：`"message": "success"`（便于前端一致处理）

## 6. 采集员数据采集接口（collector）

### 6.1 获取采集视图数据
- **路径**: `/api/collector/{sportType}/events/{matchId}`
- **方法**: `GET`
- **功能**: 供采集员查看比赛基础信息与当前比分/阵容数据
- **鉴权**: 需携带`Authorization: Bearer <token>`（登录用户即可）；网关已对`/api/collector/*`做鉴权保护
- **请求参数**:
  - 路径参数: `sportType`（如`football`、`basketball`），`matchId`（比赛ID，整数）
- **响应格式**（`message=success`）:
  ```json
  {
    "code": 200,
    "data": {
      "id": 123,
      "name": "中大杯小组赛A1",
      "teamA": "软件学院一队",
      "teamB": "物理学院一队",
      "time": "2025-12-10 19:30",
      "venue": "主体育场",
      "sportType": "football",
      "status": "ongoing",
      "currentData": {
        "scores": { "teamA": 1, "teamB": 0 },
        "lineups": {
          "teamA": [ { "id": "20230001", "name": "张三", "number": "9", "position": "FW", "isStarting": true } ],
          "teamB": [ { "id": "20230002", "name": "李四", "number": "1", "position": "GK", "isStarting": true } ]
        }
      }
    },
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/collector.go`

### 6.2 提交采集数据（比分/阵容）
- **路径**: `/api/collector/{sportType}/events/{matchId}/data`
- **方法**: `POST`
- **功能**: 采集员实时上报比赛比分与双方阵容（支持首发/号码/位置）
- **鉴权**: 需携带`Authorization: Bearer <token>`且具备`collector`角色（服务端在处理函数内校验）
- **请求体**（支持两种格式，推荐简化版）：
  - 简化版 `CollectorSubmitDataRequest`:
    ```json
    {
      "scoreTeamA": 2,
      "scoreTeamB": 1,
      "teamA": [
        { "studentId": "20230001", "position": "FW", "jerseyNumber": "9", "isStarting": true }
      ],
      "teamB": [
        { "studentId": "20230002", "position": "GK", "jerseyNumber": "1", "isStarting": true }
      ]
    }
    ```
  - 扩展版（含外层上下文 `env`）：
    ```json
    {
      "eventId": 88,
      "sportType": "football",
      "timestamp": "2025-12-10T19:45:00+08:00",
      "collectorName": "王采集",
      "status": "ongoing",
      "data": {
        "scores": { "teamA": 2, "teamB": 1 },
        "lineups": {
          "teamA": [ { "id": "20230001", "position": "FW", "number": "9", "isStarting": true } ],
          "teamB": [ { "id": "20230002", "position": "GK", "number": "1", "isStarting": true } ]
        }
      }
    }
    ```
- **响应格式**（`message=success`）:
  ```json
  {
    "code": 200,
    "data": {
      "sportType": "football",
      "matchId": 123,
      "updatedScore": true,
      "lineupsUpserted": 2
    },
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/collector.go`

## 7. 比赛评分接口（ratings）

### 7.1 获取比赛评分列表
- **路径**: `/api/matches/{matchId}/ratings`
- **方法**: `GET`
- **功能**: 列出指定比赛下所有球员的评分记录
- **鉴权**: 当前实现无需鉴权（后续可根据需求增加）
- **请求参数**:
  - 路径参数: `matchId`
- **响应格式**（`message=success`）:
  ```json
  {
    "code": 200,
    "data": [
      {
        "rating_id": 1,
        "match_id": 123,
        "student_id": "20230001",
        "rater_id": "20239999",
        "score": 7.5,
        "comment": "前场拿球稳定，射门果断",
        "created_at": "2025-12-10T19:40:00+08:00"
      }
    ],
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/rating.go`

### 7.2 提交/更新球员评分
- **路径**: `/api/matches/{matchId}/ratings`
- **方法**: `POST`
- **功能**: 对指定比赛中的某位球员打分，若存在则更新分数与评语
- **鉴权**: 当前实现无需鉴权（建议登录后操作以追踪`rater_id`）
- **请求体**:
  ```json
  {
    "student_id": "20230001",
    "rater_id": "20239999",
    "score": 8.0,
    "comment": "体能出色，贡献关键助攻"
  }
  ```
- **响应格式**（`message=success`）:
  ```json
  {
    "code": 200,
    "data": { "match_id": 123 },
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/rating.go`

## 8. 阵容接口（lineups）

### 8.1 获取比赛阵容
- **路径**: `/api/matches/{matchId}/lineups`
- **方法**: `GET`
- **功能**: 获取指定比赛的球员阵容（含队伍ID、位置、号码、是否首发）
- **鉴权**: 无需鉴权
- **响应格式**（`message=success`）:
  ```json
  {
    "code": 200,
    "data": [
      {
        "lineup_id": 11,
        "match_id": 123,
        "team_id": 1001,
        "student_id": "20230001",
        "position": "FW",
        "is_starting": true,
        "jersey_number": "9"
      }
    ],
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/lineup.go`

### 8.2 写入/更新比赛阵容
- **路径**: `/api/matches/{matchId}/lineups`
- **方法**: `POST`
- **功能**: 批量写入或更新比赛阵容条目
- **鉴权**: 需携带`Authorization: Bearer <token>`且具备`collector`角色
- **请求体**:
  ```json
  {
    "items": [
      {
        "lineup_id": 0,
        "match_id": 0,
        "team_id": 1001,
        "student_id": "20230001",
        "position": "FW",
        "is_starting": true,
        "jersey_number": "9"
      }
    ]
  }
  ```
- **响应格式**（`message=success`）:
  ```json
  {
    "code": 200,
    "data": { "match_id": 123, "updated": 1 },
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/lineup.go`

## 鉴权与响应约定

- `Authorization`：除认证/公开读取接口外，`/api/collector/*`与阵容写入接口必须携带`Bearer Token`。
- `角色要求`：采集相关写入（`/api/collector/.../data`、`POST /api/matches/{id}/lineups`）需要`collector`角色；评分接口当前无需角色限制。
- `成功消息`：统一使用`"message": "success"`作为成功提示，便于前端移动端统一处理。

## 9. 管理员接口 (Admin)

### 9.1 创建赛事
- **路径**: `/api/events/create`
- **方法**: `POST`
- **功能**: 创建新的赛事
- **鉴权**: 需 `admin` 角色 (Authorization: Bearer <token>)
- **请求体**:
  ```json
  {
    "event_name": "2025 Spring Football",
    "sport_id": 1,
    "season": "2025",
    "round": "Group Stage",
    "format_type": "points", // points 或 group_knockout
    "start_date": "2025-03-01T00:00:00Z",
    "end_date": "2025-06-01T00:00:00Z"
  }
  ```
- **响应**: 成功返回创建的赛事信息。

### 9.2 创建队伍
- **路径**: `/api/teams/create`
- **方法**: `POST`
- **功能**: 创建新的队伍
- **鉴权**: 需 `admin` 角色
- **请求体**:
  ```json
  {
    "team_name": "CS Team A",
    "sport_id": 1,
    "college": "Computer Science",
    "team_type": "college", // college 或 school
    "avatar_url": "",
    "description": "CS Department Team",
    "created_by": "admin_001"
  }
  ```
- **响应**: 成功返回创建的队伍信息。

### 9.3 创建比赛
- **路径**: `/api/matches/create`
- **方法**: `POST`
- **功能**: 创建新的比赛
- **鉴权**: 需 `admin` 角色
- **请求体**:
  ```json
  {
    "event_id": 1,
    "match_name": "CS vs Physics",
    "round": "Round 1",
    "match_time": "2025-03-02T10:00:00Z",
    "team_a_id": 101,
    "team_b_id": 102
  }
  ```
- **响应**: 成功返回创建的比赛信息。

### 9.4 队伍列表
- **路径**: `/api/teams/list`
- **方法**: `GET`
- **参数**: `sport_id` (可选)
- **功能**: 列出所有队伍

## 10. 公共数据接口

### 10.1 获取赛事数据 (积分榜/赛程)
- **路径**: `/api/match/data`
- **方法**: `GET`
- **请求参数**:
  - `sportType`: 运动类型ID (如 1)
  - `matchTime`: 时间 (格式 YYYY-MM)
  - `dataType`: 数据类型 ("积分榜", "赛程", "历史")
- **响应**: 返回对应列表数据。积分榜会自动根据比赛结果更新。

## 11. 积分榜接口 (Standings)

### 11.1 获取赛事积分榜
- **路径**: `/api/events/{event_id}/standings`
- **方法**: `GET`
- **功能**: 获取指定赛事的积分榜数据
- **请求参数**:
  - 路径参数: `event_id` (赛事ID)
- **响应格式**:
  ```json
  {
    "code": 200,
    "data": [
      {
        "event_id": 1,
        "team_id": 101,
        "team_name": "CS Team A",
        "team_logo": "http://example.com/logo.png",
        "played": 5,
        "won": 3,
        "drawn": 1,
        "lost": 1,
        "goals_for": 10,
        "goals_against": 5,
        "points": 10,
        "rank": 1
      }
    ],
    "message": "success"
  }
  ```
- **实现文件**: `backend/internal/handler/standings.go`
