Write-Host "正在启动体育赛事管理系统..." -ForegroundColor Green

# 获取脚本所在目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# 1. 启动后端 (Start Backend)
Write-Host "正在启动后端服务 (Database Mode)..." -ForegroundColor Cyan

# 构造后端启动命令
# 注意：这里我们直接将路径注入到命令字符串中，并使用双引号确保变量被解析
# DSN 中的特殊字符已经包含在单引号中
$BackendCmd = "Set-Location '$ScriptDir\backend'; `$env:DB_DSN='root:15329554862ph@tcp(127.0.0.1:3306)/sports_management?charset=utf8mb4&parseTime=True&loc=Local'; go run ./cmd/api"

Start-Process powershell -ArgumentList "-NoExit", "-Command", $BackendCmd

# 2. 启动前端 (Start Frontend)
Write-Host "正在启动前端服务..." -ForegroundColor Cyan
$FrontendCmd = "Set-Location '$ScriptDir\frontend'; npm run dev"
Start-Process powershell -ArgumentList "-NoExit", "-Command", $FrontendCmd

# 等待几秒
Start-Sleep -Seconds 5
# 手动打开浏览器到登录页
Start-Process "http://localhost:5173/login?force=true"

Write-Host "所有服务已启动！前端页面将自动打开。" -ForegroundColor Green
