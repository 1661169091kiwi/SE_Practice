# 体育赛事管理系统 - 后端（Go，极简版）

## 启动

```bash
# 在仓库根目录
go run ./backend/cmd/api
# 默认端口 8080，可用 APP_PORT 覆盖
```

健康检查：
```bash
curl http://localhost:8080/api/health
```

已提供占位接口（返回 stub 数据）：
- POST /api/register
- POST /api/login
- GET /api/events
- POST /api/events/subscribe
- GET /api/matches
- GET /api/matches/{id}

## 配置
- 环境变量：`APP_PORT`（默认 8080）
- 示例 YAML：`backend/configs/config.yaml`（可忽略，当前未强依赖）

## 如何建立本地数据库？ （backend/internal/db/migrations/0001_init.sql）

准备数据库（MySQL/MariaDB）：

方式一：命令行客户端
```bash
# 进入 mysql 客户端后，依次执行
CREATE DATABASE IF NOT EXISTS sports_management DEFAULT CHARACTER SET utf8mb4;
USE sports_management;
SOURCE c:/Users/16611/Desktop/中级实训/Code/SE_Practice/backend/internal/db/migrations/0001_init.sql;
```
> 注：Windows 下请将路径替换为你本机实际绝对路径；Linux/Mac 使用对应绝对路径。

方式二：直接重定向执行
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS sports_management DEFAULT CHARACTER SET utf8mb4;"
mysql -u root -p sports_management < backend/internal/db/migrations/0001_init.sql
```

方式三：Docker MySQL 容器
```bash
docker exec -i <mysql_container_name> mysql -uroot -p<password> -e "CREATE DATABASE IF NOT EXISTS sports_management DEFAULT CHARACTER SET utf8mb4;"
type backend/internal/db/migrations/0001_init.sql | docker exec -i <mysql_container_name> mysql -uroot -p<password> sports_management
```

说明：
- `0001_init.sql` 内已包含 `SET FOREIGN_KEY_CHECKS` 开关，便于按顺序创建表。
- 完整表结构请从 `docs/DB.md` 或 `database_struction.md` 复制粘贴到该 SQL 文件内（当前仅内置 `users` 示例表）。
- 执行成功后再在应用中接入真实 DB 连接（当前 `internal/db/db.go` 仅为占位）。

## 下一步
- 将 `database_struction.md` 的建表语句分批拷贝到 `0001_init.sql` 并执行。
- 在 `repo`/`service` 中接入真实查询与业务逻辑。
- 给路由/处理器补充参数校验与错误处理。


