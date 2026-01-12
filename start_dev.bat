@echo off
chcp 65001
cd /d "%~dp0"

echo Starting Sports Management System...

:: Set DB_DSN (escaped ampersands just in case)
set "DB_DSN=root:15329554862ph@tcp(127.0.0.1:3306)/sports_management?charset=utf8mb4^&parseTime=True^&loc=Local"

echo Starting Backend...
start "Backend API" cmd /k "cd backend && go run ./cmd/api || echo Backend failed && pause"

echo Starting Frontend...
start "Frontend" cmd /k "cd frontend && npm run dev || echo Frontend failed && pause"

echo Waiting for services to start...
ping 127.0.0.1 -n 6 >nul

echo Opening Browser...
start "" "http://localhost:5173/login?force=true"

echo Done.
pause
