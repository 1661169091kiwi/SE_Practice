package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func makeToken(sub, role string, iat, exp int64, secret string) string {
	header := `{"alg":"HS256","typ":"JWT"}`
	payload := `{"sub":"` + sub + `","iat":` + intToStr(iat) + `,"exp":` + intToStr(exp) + `,"role":"` + role + `"}`
	hb := base64.RawURLEncoding.EncodeToString([]byte(header))
	pb := base64.RawURLEncoding.EncodeToString([]byte(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(hb + "." + pb))
	sig := mac.Sum(nil)
	sb := base64.RawURLEncoding.EncodeToString(sig)
	return hb + "." + pb + "." + sb
}

func intToStr(i int64) string { return strconv.FormatInt(i, 10) }

func TestRoleAuth_AllowsCollector(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	now := time.Now().Unix()
	token := makeToken("20230001", "collector", now, now+3600, "test-secret")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/matches/update-score/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	h := RoleAuth("collector")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestRoleAuth_RejectsWrongRole(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	now := time.Now().Unix()
	token := makeToken("20230001", "viewer", now, now+3600, "test-secret")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/matches/update-score/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	h := RoleAuth("collector")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

func TestRoleAuth_RejectsExpired(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	now := time.Now().Unix()
	token := makeToken("20230001", "collector", now-7200, now-3600, "test-secret")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/matches/update-score/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	h := RoleAuth("collector")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
