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

