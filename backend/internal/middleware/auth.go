package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

func base64urlDecode(s string) ([]byte, bool) {
	b := s
	if m := len(b) % 4; m != 0 {
		b += strings.Repeat("=", 4-m)
	}
	d, err := base64.URLEncoding.DecodeString(b)
	if err != nil {
		return nil, false
	}
	return d, true
}

func verifyJWT(token string, secret string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	header := parts[0]
	payload := parts[1]
	sig := parts[2]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(header + "." + payload))
	expect := mac.Sum(nil)
	got, ok := base64urlDecode(sig)
	if !ok {
		return false
	}
	return hmac.Equal(got, expect)
}

// Claims 简易 JWT 载荷结构（只解析我们需要的字段）
type Claims struct {
	Sub  string `json:"sub"`
	Iat  int64  `json:"iat"`
	Exp  int64  `json:"exp"`
	Role string `json:"role"`
}

// extractClaims 解析 JWT 的 payload 部分
func extractClaims(token string) (*Claims, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, false
	}
	payload := parts[1]
	pb, ok := base64urlDecode(payload)
	if !ok {
		return nil, false
	}
	var c Claims
	if err := json.Unmarshal(pb, &c); err != nil {
		return nil, false
	}
	return &c, true
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev-secret"
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if !verifyJWT(token, secret) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RoleAuth 要求指定角色且未过期（若不满足，返回 403）
// 说明：仅在需要保护的路由上使用，或用于“仅保护特定方法”的场景
func RoleAuth(required string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "dev-secret"
			}
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			if !verifyJWT(token, secret) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims, ok := extractClaims(token)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// 过期校验
			if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// 角色校验
			if claims.Role != required {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CheckRole 在 handler 内部条件性调用（例如仅在 POST 方法时需要角色）
// 返回 true 表示通过，失败时已写入响应并返回 false
func CheckRole(w http.ResponseWriter, r *http.Request, required string) bool {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if !verifyJWT(token, secret) {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	claims, ok := extractClaims(token)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if claims.Role != required {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}

func AuthClaims(r *http.Request) (*Claims, bool) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, false
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if !verifyJWT(token, secret) {
		return nil, false
	}
	claims, ok := extractClaims(token)
	if !ok {
		return nil, false
	}
	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		return nil, false
	}
	return claims, true
}
