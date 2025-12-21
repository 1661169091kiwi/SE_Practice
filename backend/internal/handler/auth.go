package handler //该包负责HTTP请求处理
//前后端交互流程：
//HTTP请求->handler（HTTP请求处理层）->service（业务逻辑层）->repo（数据访问层）->数据库（前端发起请求流程）
//数据库->repo->service->handler->HTTP响应（后端响应前端流程）
//代码实现时，从底层逐步向上实现，即编码顺序：数据库->repo->service->handler->HTTP响应

//数据库层：负责与数据库交互，执行SQL语句。
//数据访问层（repo）：负责定义数据访问接口，封装数据库操作。
//业务逻辑层（service）：负责实现业务逻辑，调用数据访问层进行数据操作。
//HTTP请求处理层（handler）：负责处理HTTP请求，调用业务逻辑层进行业务处理。
//HTTP响应：负责将业务处理结果返回给前端。

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http" //Go标准库的HTTP功能
	"os"
	"strconv"
	"strings"
	"time"

	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util" //项目内部的工具包（统一响应格式）
)

var authService = service.NewAuthService() //创建一个AuthService全局变量，该变量用于处理认证相关的业务逻辑

// 处理用户注册
func Register(w http.ResponseWriter, r *http.Request) {
	//------------检查是否为Post请求---------------//
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed") //"http.StatusMethodNotAllowed"返回错误响应（405 Method Not Allowed）
		return
	}
	//Post请求在传输用户账号、密码时更安全

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}
	if req.StudentID == "" || req.Password == "" || req.Name == "" {
		util.Error(w, http.StatusBadRequest, "student_id, password, name are required")
		return
	}

	user, err := authService.Register(&req)
	if err != nil {
		if err == service.ErrUserExists {
			util.Error(w, http.StatusConflict, "user already exists")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, user)
}

// 获取用户个人资料
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// 检查是否为GET请求
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 从URL路径中提取学号: /api/user/profile/{student_id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/user/profile/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}
	studentID := parts[0]

	// 调用服务层获取用户资料
	profile, err := authService.GetUserProfile(studentID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if profile == nil {
		util.Error(w, http.StatusNotFound, "user not found")
		return
	}

	util.OK(w, profile)
}

// 更新用户学院信息
func UpdateCollege(w http.ResponseWriter, r *http.Request) {
	// 检查是否为POST请求
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析请求体
	var req model.UpdateCollegeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	// 验证请求参数
	if req.StudentID == "" || req.College == "" {
		util.Error(w, http.StatusBadRequest, "student_id and college are required")
		return
	}

	// 调用服务层更新学院信息
	response, err := authService.UpdateCollege(&req)
	if err != nil {
		if err == service.ErrUserNotFound {
			util.Error(w, http.StatusNotFound, "user not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回成功响应
	util.OK(w, response)
}

// 更新用户头像
func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	// 检查是否为POST请求
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 从URL路径中提取学号: /api/user/avatar/{student_id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/user/avatar/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}
	studentID := parts[0]

	// 设置最大文件大小限制（10MB）
	r.ParseMultipartForm(10 << 20)

	// 获取上传的文件
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		util.Error(w, http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	// 验证文件大小（不超过10MB）
	if handler.Size > 10<<20 {
		util.Error(w, http.StatusBadRequest, "file too large")
		return
	}

	// 验证文件类型
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	ext := handler.Filename[strings.LastIndex(handler.Filename, "."):]
	if !allowedTypes[strings.ToLower(ext)] {
		util.Error(w, http.StatusBadRequest, "invalid file type")
		return
	}

	// 在实际应用中，这里应该将文件保存到服务器或云存储
	// 为了简化，我们使用一个模拟的URL
	// 实际部署时需要替换为真实的文件保存逻辑
	avatarURL := "/uploads/avatars/" + studentID + ext

	// 调用服务层更新头像
	err = authService.UpdateAvatar(studentID, avatarURL)
	if err != nil {
		if err == service.ErrUserNotFound {
			util.Error(w, http.StatusNotFound, "user not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回成功响应
	response := model.UpdateAvatarResponse{
		AvatarURL: avatarURL,
		Message:   "头像更新成功",
	}
	util.OK(w, response)
}

// 修改用户密码
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	// 检查是否为POST请求
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析请求体
	var req model.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 调用服务层修改密码
	err := authService.ChangePassword(&req)
	if err != nil {
		// 根据错误类型返回不同的状态码
		switch err {
		case service.ErrInvalidCredential:
			util.Error(w, http.StatusUnauthorized, err.Error())
		case service.ErrPasswordSame:
			util.Error(w, http.StatusBadRequest, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, "failed to change password")
		}
		return
	}

	// 成功响应
	util.OK(w, map[string]string{"message": "password changed successfully"})
}

// ApplyCollector 申请成为采集员
func ApplyCollector(w http.ResponseWriter, r *http.Request) {
	// 检查是否为POST请求
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析请求体
	var req struct {
		StudentID string `json:"student_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if req.StudentID == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}

	// 调用服务层
	if err := authService.ApplyCollector(req.StudentID); err != nil {
		util.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	util.OK(w, map[string]string{"message": "申请已提交，等待审核"})
}

// ListPendingCollectors 列出待审核的采集员申请
func ListPendingCollectors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	apps, err := authService.ListPendingCollectors()
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, apps)
}

// ListCollectors 列出所有已批准的采集员
func ListCollectors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	collectors, err := authService.ListCollectors()
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, collectors)
}

// DeleteCollector 删除采集员
func DeleteCollector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}

	if err := authService.DeleteCollector(studentID); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, "collector deleted")
}

// ApproveCollector 批准采集员申请
func ApproveCollector(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	var req struct {
		StudentID string `json:"student_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.StudentID == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}

	if err := authService.ApproveCollector(req.StudentID); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, map[string]string{"message": "已批准申请"})
}

// ApplyAthlete 申请成为运动员
func ApplyAthlete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.ApplyAthleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if req.StudentID == "" || req.TeamID == 0 || req.SportType == "" {
		util.Error(w, http.StatusBadRequest, "student_id, team_id and sport_type are required")
		return
	}

	if err := authService.ApplyAthlete(&req); err != nil {
		// 根据错误类型返回不同的状态码
		switch err {
		case service.ErrUserNotFound:
			util.Error(w, http.StatusNotFound, err.Error())
		case service.ErrTeamNotFound:
			util.Error(w, http.StatusBadRequest, err.Error())
		case service.ErrAthleteExists:
			util.Error(w, http.StatusConflict, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "申请成功"})
}

// GetPendingAthletes 获取待审核的运动员申请
func GetPendingAthletes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	apps, err := authService.GetPendingAthletes()
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, apps)
}

// ApproveAthlete 批准运动员
func ApproveAthlete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	var req struct {
		TeamMemberID int64 `json:"team_member_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := authService.ApproveAthlete(req.TeamMemberID); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, map[string]string{"message": "approved"})
}

// 处理用户登录
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	user, err := authService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredential {
			util.Error(w, http.StatusUnauthorized, "invalid student_id or password")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	token := signToken(user.StudentID, user.Role)
	resp := model.LoginResponse{Token: token, User: *user}
	util.OK(w, resp)
}

// signToken 增加 role、exp 字段（默认 role=collector，过期时间 24 小时）
func signToken(studentID, role string) string {
	header := `{"alg":"HS256","typ":"JWT"}`
	now := time.Now().Unix()
	exp := now + 24*60*60
	payload := `{"sub":"` + studentID + `","iat":` + fmtInt(now) + `,"exp":` + fmtInt(exp) + `,"role":"` + role + `"}`
	hb := base64.RawURLEncoding.EncodeToString([]byte(header))
	pb := base64.RawURLEncoding.EncodeToString([]byte(payload))
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(hb + "." + pb))
	sig := mac.Sum(nil)
	sb := base64.RawURLEncoding.EncodeToString(sig)
	return hb + "." + pb + "." + sb
}

func fmtInt(i int64) string {
	return strconv.FormatInt(i, 10)
}
