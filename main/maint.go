package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

// 用户结构体
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// 登录请求
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// JWT声明
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// 用var声明变量
var jwtSecret = []byte("your-secret-key-change-in-production")
var db *sql.DB

func main() {
	// 初始化数据库
	initDB()

	// 创建Echo实例
	e := echo.New()

	// 中间件
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 路由
	e.POST("/api/login", loginHandler)
	e.GET("/api/protected", protectedHandler, jwtMiddleware)

	// 启动服务器
	log.Println("Server starting on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

// 初始化数据库
func initDB() {
	var err error
	// MySQL 连接字符串格式: "username:password@protocol(address)/dbname?param=value"
	dsn := "root:qqhyjy0514@tcp(localhost:3306)/multimodal_notes?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// 创建用户表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// 插入测试用户（使用 INSERT IGNORE 避免重复）
	hash, _ := hashPassword("123456")
	_, err = db.Exec(
		"INSERT IGNORE INTO users (email, password, name) VALUES (?, ?, ?)",
		"test@example.com", hash, "Test User",
	)
	if err != nil {
		log.Fatal("Failed to insert test user:", err)
	}

	log.Println("MySQL database initialized successfully")
}

// 密码加密
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 验证密码
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 生成JWT Token
func generateToken(userID int, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWT中间件
func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

		return next(c)
	}
}

// 登录处理器
func loginHandler(c echo.Context) error {
	var req LoginRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// 验证输入
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
	}

	// 查询用户
	var user User
	err := db.QueryRow(
		"SELECT id, email, password, name FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid credentials",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	// 验证密码
	if !checkPasswordHash(req.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// 生成token
	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not generate token",
		})
	}

	// 返回响应（不包含密码）
	user.Password = ""
	response := LoginResponse{
		Token: token,
		User:  &user,
	}

	return c.JSON(http.StatusOK, response)
}

// 受保护的路由示例
func protectedHandler(c echo.Context) error {
	userID := c.Get("userID").(int)
	userEmail := c.Get("userEmail").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Hello %s (ID: %d), this is a protected route!", userEmail, userID),
		"user_id": userID,
	})
}
