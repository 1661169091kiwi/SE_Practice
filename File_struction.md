# 项目文件结构说明（简洁版）

本结构依据 `Project_struction.md` 的页面/流程与 `database_struction.md` 的数据表与接口建议进行设计，目标：简单清晰、功能分区明确、便于扩展。

## 目录总览（Go 后端）

```text
SE_Practice/
├─ backend/
│  ├─ cmd/api/main.go              # 程序入口
│  ├─ internal/
│  │  ├─ router/router.go          # 路由注册（按模块挂载）
│  │  ├─ handler/                  # HTTP 处理器（少量文件即可）
│  │  │  ├─ auth.go                # 登录/注册
│  │  │  ├─ user.go                # 用户资料
│  │  │  ├─ event.go               # 赛事/订阅
│  │  │  └─ match.go               # 比赛/阵容/评分（先合并）
│  │  ├─ service/                  # 业务逻辑
│  │  │  ├─ auth_service.go
│  │  │  ├─ user_service.go
│  │  │  ├─ event_service.go
│  │  │  └─ match_service.go
│  │  ├─ model/                    # 数据模型定义
│  │  │  └─ user.go                # 用户相关数据结构
│  │  ├─ repo/                     # 数据访问（SQL/ORM）
│  │  │  ├─ user_repo.go
│  │  │  ├─ event_repo.go
│  │  │  └─ match_repo.go
│  │  ├─ db/
│  │  │  ├─ db.go                  # 连接与初始化
│  │  │  └─ migrations/            # 迁移脚本（来源于 DB.md）
│  │  │     └─ 0001_init.sql
│  │  ├─ middleware/
│  │  │  └─ auth.go                # 简单 JWT 校验（可选）
│  │  └─ util/
│  │     └─ response.go            # 统一响应格式
│  ├─ configs/config.yaml          # 配置示例（端口、DB、JWT）
│  └─ go.mod                       # Go 依赖
└─ .env.example                    # 环境变量示例
```

说明：
- `handler/event.go` 覆盖：赛事列表、订阅；`handler/match.go` 合并比赛、阵容、评分的基础接口。
- `repo/*` 必要的三个仓储：用户、赛事、比赛。

## 最小接口集合（建议首批）
- 认证：POST `/api/register`、POST `/api/login`
- 赛事：GET `/api/events`、POST `/api/events/subscribe`
- 比赛：GET `/api/matches`、GET `/api/matches/{id}`
- 评分/阵容：POST `/api/ratings`、POST `/api/lineups`

## 运行与配置（精简）
- 配置从环境变量或 `configs/config.yaml` 加载（DB 连接、端口、JWT 密钥）。
- `db/migrations` 中保存建表 SQL，按需手动执行或写简单迁移脚本。
- 先不引入复杂脚手架与 Makefile，直接 `go run ./backend/cmd/api` 启动。

