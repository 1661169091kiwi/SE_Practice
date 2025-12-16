# 团队协作：本地部署与运行指南

本文面向 GitHub/Gitee 协作者，说明如何在本地快速部署并运行本项目（前后端联调）。

## 环境准备
- Node.js：推荐 20+（前端开发）
- Go：推荐 1.22+（后端开发）
- MySQL：5.7+（推荐 8.0）
- Git：用于拉取与协作

## 获取代码
```bash
git clone <your_repo_url>
cd SE_Practice
```

## 数据库初始化
1. 启动本地 MySQL
2. 创建数据库：
   ```sql
   CREATE DATABASE sports_management;
   ```
3. 后端会进行自动建表与初始化（见后端启动日志）

## 启动后端（API）
进入 `backend` 目录：
```powershell
# Windows（PowerShell）
$env:DB_DSN='root:<your_password>@tcp(localhost:3306)/sports_management?charset=utf8mb4&parseTime=True&loc=Local'
$env:APP_PORT='8081'
go run ./cmd/api
```
```bash
# macOS / Linux（bash/zsh）
export DB_DSN='root:<your_password>@tcp(localhost:3306)/sports_management?charset=utf8mb4&parseTime=True&loc=Local'
export APP_PORT='8081'
go run ./cmd/api
```
看到如下日志表示成功：
```
Database connection established successfully
Database schema initialized successfully
server starting at :8081
```

## 启动前端（学生端/采集员端）
进入 `frontend` 目录：
```bash
npm install
```
配置后端地址并启动开发服务器：
```powershell
# Windows（PowerShell）
$env:VITE_API_BASE_URL='http://localhost:8081/api'
npm run dev
```
```bash
# macOS / Linux（bash/zsh）
export VITE_API_BASE_URL='http://localhost:8081/api'
npm run dev
```
访问：`http://localhost:5173/`

## 常用入口
- 学生端赛事列表：`/student/events`
- 采集员赛事选择：`/events`
- 比赛详情页：`/match/:id`
- 积分榜：`/student/standings/:id`

## 开发与协作建议
- 前端脚本：
  - `npm run dev` 开发
  - `npm run build` 构建
  - `npm run preview` 预览
  - `npm run lint` 代码检查
- 提交前建议运行 `npm run lint`，确保基础代码质量

## 常见问题排查
- 端口占用
  - 后端：修改 `APP_PORT`
  - 前端：`vite --port <port>` 或修改 `Vite` 配置
- 接口访问失败
  - 确认后端运行中，访问 `http://localhost:8081/api/...`
  - 前端 `VITE_API_BASE_URL` 是否设置为 `http://localhost:8081/api`
- 数据库连接失败
  - 检查 `DB_DSN` 中的账号密码与数据库是否存在
  - 确保 MySQL 已启动，端口默认 `3306`

## 如何停止服务
- 后端：在运行终端按 `Ctrl + C`
- 前端：在运行终端按 `Ctrl + C`

## 目录结构参考
```
SE_Practice/
├─ backend/                # Go 后端（API 服务）
│  ├─ cmd/api              # 入口命令
│  └─ internal/...         # 业务与服务实现
├─ frontend/               # Vue 前端（学生端/采集员端）
│  ├─ src/views            # 页面
│  └─ src/stores           # 状态管理
└─ LOCAL_DEPLOYMENT.md     # 本地部署与运行指南（当前文件）
```
