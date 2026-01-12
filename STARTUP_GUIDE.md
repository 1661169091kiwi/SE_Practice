# 完整项目启动指南 (Full Project Startup Guide)

本文档记录了如何在**数据库模式**下完整启动本项目（前端 + 后端 + 数据库）。

## 1. 数据库准备 (Database Preparation)

确保本地 MySQL 服务已启动，且端口为 `3306`。

- **数据库名称**: `sports_management`
- **账号**: `root`
- **密码**: `15329554862ph` (根据您的环境配置)

如果尚未初始化数据库表结构，后端启动时会自动尝试初始化。

## 2. 后端启动 (Backend Startup)

后端服务负责处理业务逻辑并与 MySQL 数据库交互。

1. 打开终端（PowerShell）。
2. 进入 `backend` 目录。
3. 设置环境变量并启动服务：

```powershell
# 进入后端目录
cd backend

# 设置数据库连接字符串 (请确保密码正确)
$env:DB_DSN='root:15329554862ph@tcp(localhost:3306)/sports_management?charset=utf8mb4&parseTime=True&loc=Local'

# 设置应用端口 (必须为 8080 以匹配前端配置)
$env:APP_PORT='8080'

# 启动后端
go run ./cmd/api
```

> **成功标志**: 终端输出 `Database connection established successfully` 和 `server starting at :8080`。

## 3. 前端启动 (Frontend Startup)

前端界面基于 Vue 3 + Vite。

1. 打开一个新的终端窗口。
2. 进入 `frontend` 目录。
3. 启动开发服务器：

```powershell
# 进入前端目录
cd frontend

# 启动开发服务器
npm run dev
```

> **成功标志**: 终端输出 `Local: http://localhost:5173/`。

## 4. 访问项目 (Access)

- **前端页面**: [http://localhost:5173](http://localhost:5173)
- **后端 API**: `http://localhost:8080/api`

## 5. 常用测试账号

- **管理员**:
    - 账号: `admin`
    - 密码: `admin123`
- **普通用户**:
    - 您可以在登录页面点击“注册”创建一个新账号，数据将保存到数据库中。
