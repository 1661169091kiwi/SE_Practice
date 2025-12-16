# 体育赛事管理系统（SE_Practice）

一个用于校园体育赛事管理与数据采集的全栈项目，包含学生端赛事浏览与订阅、采集员端数据录入、比赛详情与积分榜展示等功能。

## 技术栈
- 前端：Vue 3、Vite、Pinia、Vue Router
- 后端：Go、MySQL
- 开发工具：Node.js（>= 20）、Go（>= 1.22）、Vite Dev Server

## 目录结构
```
SE_Practice/
├─ backend/                # Go 后端（API 服务）
│  ├─ cmd/api              # 入口命令
│  └─ internal/...         # 业务与服务实现
├─ frontend/               # Vue 前端（学生端/采集员端）
│  ├─ src/views            # 主要页面
│  └─ src/stores           # 状态管理（Pinia）
├─ database_struction.md   # 数据库设计与建表语句
└─ README.md               # 项目说明（当前文件）
```

## 快速开始
> 以下步骤以本地开发为例，默认后端端口 `8081`，前端端口 `5173`。

### 1. 准备环境
- 安装 Node.js（推荐 20+）
- 安装 Go（推荐 1.22+）
- 安装并启动 MySQL，创建数据库：
  ```
  CREATE DATABASE sports_management;
  ```

### 2. 启动后端
在 `backend` 目录：
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
成功日志示例：
```
Database connection established successfully
Database schema initialized successfully
server starting at :8081
```

### 3. 启动前端
在 `frontend` 目录：
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
打开浏览器访问：`http://localhost:5173/`

## 常用脚本（前端）
- 开发启动：`npm run dev`
- 生产构建：`npm run build`
- 预览构建：`npm run preview`
- 代码检查：`npm run lint`

## 访问入口
- 学生端赛事列表：`/student/events`
- 采集员赛事选择：`/events`
- 比赛详情页：`/match/:id`
- 积分榜：`/student/standings/:id`

## 常见问题
- 端口占用：修改 `APP_PORT` 或通过 `vite --port` 指定前端端口
- 接口 404：确认前端 `VITE_API_BASE_URL` 指向后端 `http://localhost:8081/api`
- 数据库连接失败：检查 `DB_DSN` 格式、数据库是否已创建、账号密码是否正确

## 本地部署与运行指南
为协作者准备了更详细的步骤说明，参见：`LOCAL_DEPLOYMENT.md`
