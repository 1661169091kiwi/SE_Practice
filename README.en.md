# 软件工程中级实训代码仓库

#### Description
{**When you're done, you can delete the content in this README and update the file with details for others getting started with your repository**}

#### Software Architecture
Software architecture description

#### Installation

1.  xxxx
2.  xxxx
3.  xxxx

#### Instructions

1.  xxxx
2.  xxxx
3.  xxxx

#### Contribution

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request


#### Gitee Feature

1.  You can use Readme\_XXX.md to support different languages, such as Readme\_en.md, Readme\_zh.md
2.  Gitee blog [blog.gitee.com](https://blog.gitee.com)
3.  Explore open source project [https://gitee.com/explore](https://gitee.com/explore)
4.  The most valuable open source project [GVP](https://gitee.com/gvp)
5.  The manual of Gitee [https://gitee.com/help](https://gitee.com/help)
6.  The most popular members  [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)


### Start the project
# 本地启动与预览

## 环境要求
- Node.js 20+（或 22+）
- Go 1.21+
- MySQL 已运行，并已初始化 `sports_management` 数据库（脚本在 `backend/internal/db/migrations/0001_init.sql`）

## 1. 启动后端（端口 8080）
```powershell
cd c:\Users\16611\Desktop\中级实训\Code\SE_Practice\backend
$env:APP_PORT = "8080"
$env:DB_DSN = "root:你的密码@tcp(localhost:3306)/sports_management?charset=utf8mb4&parseTime=True&loc=Local"
go run cmd/api/main.go
```

验证健康：
```powershell
curl http://localhost:8080/api/health
# 预期 {"code":200,"message":"success","data":{"status":"ok"}}
```

## 2. 启动前端（两种方式）

### 方式A：开发模式（热更新）
```powershell
cd c:\Users\16611\Desktop\中级实训\Code\SE_Practice\frontend
npm install
npm run dev
# 访问 http://localhost:5173/
```

### 方式B：构建后预览（静态预览）
```powershell
cd c:\Users\16611\Desktop\中级实训\Code\SE_Practice\frontend
npm install
npm run build
npm run preview
# 访问 http://localhost:4173/
```

> 前端默认请求后端 `http://localhost:8080/api`。如需修改：
```powershell
# 示例：切换到其它后端地址
$env:VITE_API_BASE_URL = "http://localhost:8080/api"
npm run dev
```

## 3. 常用接口
- 健康检查: `GET http://localhost:8080/api/health`
- 登录接口: `POST http://localhost:8080/api/login`
- 比赛列表: `GET http://localhost:8080/api/matches`
- 阵容读取: `GET http://localhost:8080/api/matches/{matchId}/lineups`
- 阵容删除（采集员权限）: `DELETE http://localhost:8080/api/matches/{matchId}/lineups`

## 4. 常见问题排查
- `ERR_CONNECTION_REFUSED`：后端未启动或端口不通。先访问 `http://localhost:8080/api/health` 验证后端是否在线。
- 登录失败 401：检查前端是否携带正确的 `Bearer Token`；测试登录时可先注册再登录。
- CORS 问题：后端已启用跨域（允许 `GET, POST, PUT, PATCH, DELETE, OPTIONS`），正常不会阻塞。